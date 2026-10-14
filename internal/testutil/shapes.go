// Package testutil – синтетические фигуры и STL-энкодеры для тестов и testdata/sample.
// Никаких реальных медицинских данных.
package testutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ssqwet1324/viewing/internal/stl"
)

// Cube – куб [0,size]³ из 12 треугольников, обход против часовой снаружи.
func Cube(size float32) []stl.Triangle {
	s := size
	v := [8]stl.Vec3{
		{0, 0, 0}, {s, 0, 0}, {s, s, 0}, {0, s, 0},
		{0, 0, s}, {s, 0, s}, {s, s, s}, {0, s, s},
	}
	quads := [6][4]int{
		{0, 3, 2, 1}, // низ (-Z)
		{4, 5, 6, 7}, // верх (+Z)
		{0, 1, 5, 4}, // -Y
		{2, 3, 7, 6}, // +Y
		{1, 2, 6, 5}, // +X
		{3, 0, 4, 7}, // -X
	}
	out := make([]stl.Triangle, 0, 12)
	for _, q := range quads {
		out = append(out,
			stl.Triangle{v[q[0]], v[q[1]], v[q[2]]},
			stl.Triangle{v[q[0]], v[q[2]], v[q[3]]})
	}
	return out
}

// Cylinder – закрытый цилиндр вдоль Z от z=0 до z=h, n сегментов.
// Вершин после склейки: 2n + 2, треугольников: 4n.
func Cylinder(r, h float32, n int) []stl.Triangle {
	ring := func(z float32) []stl.Vec3 {
		pts := make([]stl.Vec3, n)
		for i := range n {
			a := 2 * math.Pi * float64(i) / float64(n)
			pts[i] = stl.Vec3{r * float32(math.Cos(a)), r * float32(math.Sin(a)), z}
		}
		return pts
	}
	lo, hi := ring(0), ring(h)
	cLo, cHi := stl.Vec3{0, 0, 0}, stl.Vec3{0, 0, h}
	out := make([]stl.Triangle, 0, 4*n)
	for i := range n {
		j := (i + 1) % n
		out = append(out,
			stl.Triangle{lo[i], lo[j], hi[j]},
			stl.Triangle{lo[i], hi[j], hi[i]},
			stl.Triangle{cHi, hi[i], hi[j]},
			stl.Triangle{cLo, lo[j], lo[i]})
	}
	return out
}

// Translate сдвигает все вершины.
func Translate(tris []stl.Triangle, d stl.Vec3) []stl.Triangle {
	out := make([]stl.Triangle, len(tris))
	for i, t := range tris {
		for k := range 3 {
			for a := range 3 {
				out[i][k][a] = t[k][a] + d[a]
			}
		}
	}
	return out
}

// EncodeBinary пишет binary STL. Header обрезается или дополняется до 80 байт;
// "solid ..." в header – валидный и коварный случай.
func EncodeBinary(tris []stl.Triangle, header string) []byte {
	buf := make([]byte, 84, 84+50*len(tris))
	copy(buf[:80], header)
	binary.LittleEndian.PutUint32(buf[80:], uint32(len(tris)))
	for _, t := range tris {
		buf = append(buf, make([]byte, 12)...) // нормаль = 0, читатель её игнорирует
		for _, v := range t {
			for _, c := range v {
				buf = binary.LittleEndian.AppendUint32(buf, math.Float32bits(c))
			}
		}
		buf = append(buf, 0, 0)
	}
	return buf
}

// EncodeASCII пишет ASCII STL. Формат %g с точностью float32 даёт точный round-trip.
func EncodeASCII(tris []stl.Triangle, name string) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "solid %s\n", name)
	for _, t := range tris {
		b.WriteString("  facet normal 0 0 0\n    outer loop\n")
		for _, v := range t {
			fmt.Fprintf(&b, "      vertex %s %s %s\n", f32(v[0]), f32(v[1]), f32(v[2]))
		}
		b.WriteString("    endloop\n  endfacet\n")
	}
	fmt.Fprintf(&b, "endsolid %s\n", name)
	return b.Bytes()
}

func f32(x float32) string { return fmt.Sprint(x) }
