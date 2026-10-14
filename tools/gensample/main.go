// gensample пишет синтетическую «челюсть» в testdata/sample: дугу кости, импланты,
// абатменты, втулки и wax-up. Координаты как в FORMAT.md: Z вверх, −Y вперёд.
//
//	go run ./tools/gensample
package main

import (
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/ssqwet1324/viewing/internal/stl"
	"github.com/ssqwet1324/viewing/internal/testutil"
)

const outDir = "testdata/sample"

func main() {
	write("bone.stl", testutil.EncodeBinary(arch(25, 8, 14, 96), "solid synthetic bone"))

	var wax []stl.Triangle
	for i := range 10 {
		a := math.Pi * (0.08 + 0.84*float64(i)/9)
		c := archPoint(25, a)
		wax = append(wax, testutil.Translate(testutil.Cube(5), stl.Vec3{c[0] - 2.5, c[1] - 2.5, 14.5})...)
	}
	write("waxup.stl", testutil.EncodeBinary(wax, "waxup"))

	// Импланты стоят в кости с наклоном, вершина (платформа) на уровне гребня.
	sites := []struct {
		tooth string
		angle float64 // положение на дуге, рад
		tilt  float64 // наклон оси, рад
	}{{"15", 0.2 * math.Pi, 0.2}, {"12", 0.42 * math.Pi, 0.1}, {"22", 0.58 * math.Pi, -0.1}, {"25", 0.8 * math.Pi, -0.2}}
	for _, s := range sites {
		top := archPoint(25, s.angle)
		top[2] = 14
		place := func(tris []stl.Triangle) []stl.Triangle {
			return testutil.Translate(rotX(tris, s.tilt), top)
		}
		write("implant_"+s.tooth+".stl", testutil.EncodeBinary(place(threaded(1.9, 11, 48)), "implant"))
		write("abutment_"+s.tooth+".stl", testutil.EncodeASCII(place(testutil.Translate(testutil.Cylinder(1.6, 4, 24), stl.Vec3{0, 0, 0})), "abutment"))
		write("sleeve_"+s.tooth+".stl", testutil.EncodeBinary(place(testutil.Translate(tube(2.6, 3.4, 5, 32), stl.Vec3{0, 0, 6})), "sleeve"))
	}
}

func write(name string, data []byte) {
	if err := os.WriteFile(filepath.Join(outDir, name), data, 0o644); err != nil {
		log.Fatal(err)
	}
}

func archPoint(r, a float64) stl.Vec3 {
	return stl.Vec3{float32(r * math.Cos(a)), float32(-r*math.Sin(a)) + 5, 0}
}

// arch – прямоугольное сечение w×h, протянутое по полуокружности радиуса r.
func arch(r, w, h float64, n int) []stl.Triangle {
	sec := func(a float64) [4]stl.Vec3 {
		c, s := math.Cos(a), math.Sin(a)
		p := func(rr, z float64) stl.Vec3 { return stl.Vec3{float32(rr * c), float32(-rr*s) + 5, float32(z)} }
		return [4]stl.Vec3{p(r-w/2, 0), p(r+w/2, 0), p(r+w/2, h), p(r-w/2, h)}
	}
	var out []stl.Triangle
	quad := func(a, b, c, d stl.Vec3) { out = append(out, stl.Triangle{a, b, c}, stl.Triangle{a, c, d}) }
	prev := sec(0)
	quad(prev[0], prev[3], prev[2], prev[1])
	for i := 1; i <= n; i++ {
		cur := sec(math.Pi * float64(i) / float64(n))
		for k := range 4 {
			j := (k + 1) % 4
			quad(prev[k], prev[j], cur[j], cur[k])
		}
		prev = cur
	}
	quad(prev[0], prev[1], prev[2], prev[3])
	return out
}

// threaded – «имплант»: конус с резьбой, ось −Z от 0 до -h. Нужен, чтобы проверить crease-нормали.
func threaded(r, h float64, n int) []stl.Triangle {
	const rings = 60
	ringAt := func(k int) []stl.Vec3 {
		z := -h * float64(k) / rings
		rr := r * (1 - 0.25*float64(k)/rings)
		if k%2 == 1 && k < rings-2 {
			rr *= 0.85
		}
		pts := make([]stl.Vec3, n)
		for i := range n {
			a := 2 * math.Pi * float64(i) / float64(n)
			pts[i] = stl.Vec3{float32(rr * math.Cos(a)), float32(rr * math.Sin(a)), float32(z)}
		}
		return pts
	}
	var out []stl.Triangle
	prev := ringAt(0)
	top := stl.Vec3{0, 0, 0}
	for i := range n {
		out = append(out, stl.Triangle{top, prev[i], prev[(i+1)%n]})
	}
	for k := 1; k <= rings; k++ {
		cur := ringAt(k)
		for i := range n {
			j := (i + 1) % n
			out = append(out, stl.Triangle{prev[i], cur[i], cur[j]}, stl.Triangle{prev[i], cur[j], prev[j]})
		}
		prev = cur
	}
	bot := stl.Vec3{0, 0, float32(-h)}
	for i := range n {
		out = append(out, stl.Triangle{bot, prev[(i+1)%n], prev[i]})
	}
	return out
}

// tube – полый цилиндр (втулка) вдоль Z.
func tube(ri, ro, h float64, n int) []stl.Triangle {
	p := func(r, a, z float64) stl.Vec3 {
		return stl.Vec3{float32(r * math.Cos(a)), float32(r * math.Sin(a)), float32(z)}
	}
	var out []stl.Triangle
	quad := func(a, b, c, d stl.Vec3) { out = append(out, stl.Triangle{a, b, c}, stl.Triangle{a, c, d}) }
	for i := range n {
		a0, a1 := 2*math.Pi*float64(i)/float64(n), 2*math.Pi*float64(i+1)/float64(n)
		quad(p(ro, a0, 0), p(ro, a1, 0), p(ro, a1, h), p(ro, a0, h))
		quad(p(ri, a1, 0), p(ri, a0, 0), p(ri, a0, h), p(ri, a1, h))
		quad(p(ri, a0, h), p(ro, a0, h), p(ro, a1, h), p(ri, a1, h))
		quad(p(ri, a1, 0), p(ro, a1, 0), p(ro, a0, 0), p(ri, a0, 0))
	}
	return out
}

func rotX(tris []stl.Triangle, a float64) []stl.Triangle {
	c, s := float32(math.Cos(a)), float32(math.Sin(a))
	out := make([]stl.Triangle, len(tris))
	for i, t := range tris {
		for k, v := range t {
			out[i][k] = stl.Vec3{v[0], c*v[1] - s*v[2], s*v[1] + c*v[2]}
		}
	}
	return out
}
