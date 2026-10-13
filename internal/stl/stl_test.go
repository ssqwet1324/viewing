package stl_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ssqwet1324/viewing/internal/stl"
	"github.com/ssqwet1324/viewing/internal/testutil"
)

func TestParse(t *testing.T) {
	cyl := testutil.Translate(testutil.Cylinder(2.1, 11.5, 24), stl.Vec3{-13.96, -33.75, 22.75})
	tests := []struct {
		name string
		data []byte
	}{
		{"binary", testutil.EncodeBinary(cyl, "dentview test")},
		// Заголовок начинается с "solid" – детект по префиксу здесь ломается.
		{"binary with solid header", testutil.EncodeBinary(cyl, "solid exported by CAD")},
		{"ascii", testutil.EncodeASCII(cyl, "cyl")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stl.Parse(tt.data)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, cyl) {
				t.Fatalf("треугольники не совпали: got %d, want %d", len(got), len(cyl))
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	bin := testutil.EncodeBinary(testutil.Cube(1), "solid x")
	tests := []struct {
		name string
		data []byte
	}{
		{"empty", nil},
		{"truncated binary", bin[:len(bin)-7]},
		{"ascii vertex count not multiple of 3", []byte("solid a\nvertex 0 0 0\nvertex 1 0 0\nendsolid a\n")},
		{"ascii bad number", []byte("solid a\nvertex 0 0 x\nvertex 1 0 0\nvertex 0 1 0\nendsolid a\n")},
		{"binary zero triangles", testutil.EncodeBinary(nil, "")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := stl.Parse(tt.data); err == nil {
				t.Fatal("ожидалась ошибка")
			}
		})
	}
	t.Run("zero triangles is ErrEmpty", func(t *testing.T) {
		if _, err := stl.Parse(testutil.EncodeBinary(nil, "")); !errors.Is(err, stl.ErrEmpty) {
			t.Fatalf("err = %v, want ErrEmpty", err)
		}
	})
}

func FuzzParse(f *testing.F) {
	f.Add(testutil.EncodeBinary(testutil.Cube(1), "solid"))
	f.Add(testutil.EncodeASCII(testutil.Cube(1), "c"))
	f.Fuzz(func(t *testing.T, data []byte) {
		tris, err := stl.Parse(data) // главное – не паниковать
		if err == nil && len(tris) == 0 {
			t.Fatal("nil error при пустом результате")
		}
	})
}
