package mesh_test

import (
	"math"
	"testing"

	"github.com/ssqwet1324/viewing/internal/mesh"
	"github.com/ssqwet1324/viewing/internal/stl"
	"github.com/ssqwet1324/viewing/internal/testutil"
)

// Склейка не должна менять геометрию: каждый треугольник восстанавливается как был.
func checkSameTriangles(t *testing.T, want []stl.Triangle, m mesh.Indexed) {
	t.Helper()
	if len(m.Indices) != 3*len(want) {
		t.Fatalf("индексов %d, want %d", len(m.Indices), 3*len(want))
	}
	for i, tri := range want {
		for k := range 3 {
			if got := m.Positions[m.Indices[3*i+k]]; got != tri[k] {
				t.Fatalf("tri %d vertex %d: %v != %v", i, k, got, tri[k])
			}
		}
	}
}

func TestWeld(t *testing.T) {
	tests := []struct {
		name  string
		tris  []stl.Triangle
		verts int
	}{
		{"cube", testutil.Cube(10), 8},
		{"cylinder", testutil.Cylinder(2, 10, 32), 2*32 + 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mesh.Weld(tt.tris)
			if len(m.Positions) != tt.verts {
				t.Fatalf("вершин %d, want %d", len(m.Positions), tt.verts)
			}
			checkSameTriangles(t, tt.tris, m)
		})
	}
}

func TestWeldNegativeZero(t *testing.T) {
	nz := float32(math.Copysign(0, -1))
	tris := []stl.Triangle{
		{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}},
		{{nz, nz, 0}, {0, 1, 0}, {-1, 0, 0}},
	}
	if m := mesh.Weld(tris); len(m.Positions) != 4 {
		t.Fatalf("вершин %d, want 4: -0 и +0 должны склеиться", len(m.Positions))
	}
}

func TestWeldDropsDegenerate(t *testing.T) {
	tris := []stl.Triangle{
		{{0, 0, 0}, {1, 0, 0}, {0, 1, 0}},
		{{0, 0, 0}, {0, 0, 0}, {0, 1, 0}}, // вырожденный
	}
	if m := mesh.Weld(tris); len(m.Indices) != 3 {
		t.Fatalf("индексов %d, want 3", len(m.Indices))
	}
}

func TestQuantizeError(t *testing.T) {
	tris := testutil.Translate(testutil.Cylinder(2.1, 11.5, 64), stl.Vec3{-56.96, -69.12, -18.24})
	m := mesh.Weld(tris)
	q := mesh.Quantize(m)
	if q.NumVertices() != len(m.Positions) || len(q.Indices) != len(m.Indices) {
		t.Fatal("размеры не совпали")
	}
	for i, p := range m.Positions {
		got := q.Vertex(i)
		for a := range 3 {
			// scale/2 – теоретический предел; запас на float32-арифметику.
			lim := float64(q.Scale[a])/2 + 1e-5
			if d := math.Abs(float64(got[a] - p[a])); d > lim {
				t.Fatalf("вершина %d ось %d: ошибка %g > %g", i, a, d, lim)
			}
		}
	}
}

func TestQuantizeFlatAxis(t *testing.T) {
	// Все z одинаковые: scale по z не должен быть 0 или NaN.
	m := mesh.Weld([]stl.Triangle{{{0, 0, 5}, {1, 0, 5}, {0, 1, 5}}})
	q := mesh.Quantize(m)
	for i := range 3 {
		if v := q.Vertex(i); v[2] != 5 {
			t.Fatalf("z = %v, want 5", v[2])
		}
	}
	if q.Scale[2] == 0 || math.IsNaN(float64(q.Scale[2])) {
		t.Fatalf("scale z = %v", q.Scale[2])
	}
}
