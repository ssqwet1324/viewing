package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ssqwet1324/viewing/internal/testutil"
)

// Сквозной тест: STL на диске → run → HTML.
func TestRunEndToEnd(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, data []byte) {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("bone.stl", testutil.EncodeBinary(testutil.Cube(20), "solid bone"))
	write("impl.stl", testutil.EncodeASCII(testutil.Cylinder(2, 10, 24), "impl"))
	write("plan.json", []byte(`{"title":"e2e","meshes":[
		{"file":"bone.stl","kind":"bone"},
		{"file":"impl.stl","kind":"implant","name":"16"}]}`))

	out := filepath.Join(dir, "out.html")
	if err := run([]string{"-c", filepath.Join(dir, "plan.json"), "-o", out}, io.Discard); err != nil {
		t.Fatal(err)
	}
	html, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{"<title>e2e</title>", `id="dv-scene"`, `id="dv-viewer"`} {
		if !strings.Contains(string(html), s) {
			t.Errorf("нет %q", s)
		}
	}
	if _, err := os.Stat(out + ".tmp"); !os.IsNotExist(err) {
		t.Error("остался .tmp")
	}
}
