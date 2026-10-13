package scene_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/ssqwet1324/viewing/internal/mesh"
	"github.com/ssqwet1324/viewing/internal/scene"
	"github.com/ssqwet1324/viewing/internal/stl"
	"github.com/ssqwet1324/viewing/internal/testutil"
)

var update = flag.Bool("update", false, "перезаписать golden-файлы")

// decodeGeom – независимый от Builder декодер тела: так тест ловит расхождение со спекой.
func decodeGeom(t *testing.T, g scene.Geom, body []byte) (pos []uint16, idx []uint32) {
	t.Helper()
	if g.PosOff%scene.SectionAlign != 0 || g.IdxOff%scene.SectionAlign != 0 {
		t.Fatalf("секции не выровнены: %d, %d", g.PosOff, g.IdxOff)
	}
	if g.PosLen != g.Vertices*3*2 {
		t.Fatalf("pos_len %d, want %d", g.PosLen, g.Vertices*6)
	}
	p := body[g.PosOff : g.PosOff+g.PosLen]
	for i := 0; i < len(p); i += 2 {
		pos = append(pos, binary.LittleEndian.Uint16(p[i:]))
	}
	ix := body[g.IdxOff : g.IdxOff+g.IdxLen]
	switch g.Index {
	case scene.IndexU16:
		for i := 0; i < len(ix); i += 2 {
			idx = append(idx, uint32(binary.LittleEndian.Uint16(ix[i:])))
		}
	case scene.IndexU32:
		for i := 0; i < len(ix); i += 4 {
			idx = append(idx, binary.LittleEndian.Uint32(ix[i:]))
		}
	default:
		t.Fatalf("index = %q", g.Index)
	}
	if len(idx) != g.Triangles*3 {
		t.Fatalf("индексов %d, want %d", len(idx), g.Triangles*3)
	}
	return pos, idx
}

func TestBuilderRoundTrip(t *testing.T) {
	small := mesh.Quantize(mesh.Weld(testutil.Cylinder(2, 10, 7))) // нечётные длины → проверка добивки
	// > 65536 вершин → u32. Геометрия не обязана быть осмысленной.
	big := mesh.Quantized{Scale: [3]float32{1, 1, 1}, Pos: make([]uint16, 3*70000), Indices: []uint32{0, 69999, 12345}}
	for i := range big.Pos {
		big.Pos[i] = uint16(i * 7)
	}

	b := scene.NewBuilder("t")
	qs := []mesh.Quantized{small, big, small}
	for _, q := range qs {
		b.AddGeom(q)
	}
	hdr, body := b.Finish()
	if len(hdr.Geoms) != len(qs) {
		t.Fatalf("geoms %d", len(hdr.Geoms))
	}
	wantIdx := []scene.IndexType{scene.IndexU16, scene.IndexU32, scene.IndexU16}
	for i, q := range qs {
		g := hdr.Geoms[i]
		if g.Index != wantIdx[i] {
			t.Errorf("geom %d: index %q, want %q", i, g.Index, wantIdx[i])
		}
		pos, idx := decodeGeom(t, g, body)
		if !slices.Equal(pos, q.Pos) || !slices.Equal(idx, q.Indices) {
			t.Errorf("geom %d: данные не совпали", i)
		}
		if g.Origin != q.Origin || g.Scale != q.Scale {
			t.Errorf("geom %d: origin/scale", i)
		}
	}
}

func TestAddMeshBadGeom(t *testing.T) {
	b := scene.NewBuilder("t")
	if err := b.AddMesh(scene.Mesh{ID: "x", Geom: 0}); err == nil {
		t.Fatal("ожидалась ErrBadGeom")
	}
}

// TestGolden фиксирует JSON-заголовок небольшой сцены. Любое изменение формата видно в diff.
// Обновить: go test ./internal/scene -run Golden -update
func TestGolden(t *testing.T) {
	b := scene.NewBuilder("Синтетика")
	bone := b.AddGeom(mesh.Quantize(mesh.Weld(testutil.Translate(testutil.Cube(20), stl.Vec3{-10, -10, 0}))))
	impl := b.AddGeom(mesh.Quantize(mesh.Weld(testutil.Cylinder(2, 10, 16))))
	tr := scene.Matrix4{1, 0, 0, 5, 0, 1, 0, -3, 0, 0, 1, 8, 0, 0, 0, 1}
	meshes := []scene.Mesh{
		{ID: "bone", Name: "Кость", Kind: scene.KindBone, Color: "#e8e0d0", Opacity: 0.7, Visible: true, Shading: scene.ShadingSmooth, Geom: bone},
		{ID: "i16", Name: "16", Kind: scene.KindImplant, Color: "#9aa6b4", Opacity: 1, Visible: true, Shading: scene.ShadingCrease, CreaseDeg: 30, Transform: &tr, Geom: impl},
	}
	for _, m := range meshes {
		if err := b.AddMesh(m); err != nil {
			t.Fatal(err)
		}
	}
	hdr, _ := b.Finish()
	got, err := json.MarshalIndent(hdr, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got = append(got, '\n')

	path := filepath.Join("testdata", "scene.golden.json")
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (запустите с -update)", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("заголовок отличается от %s:\n%s", path, got)
	}
}
