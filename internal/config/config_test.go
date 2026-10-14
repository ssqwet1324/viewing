package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ssqwet1324/viewing/internal/config"
	"github.com/ssqwet1324/viewing/internal/scene"
)

func writeConfig(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "plan.json")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadDefaults(t *testing.T) {
	p := writeConfig(t, `{
		"title": "Тест",
		"meshes": [
			{"file": "Bone.stl", "kind": "bone"},
			{"file": "/abs/impl.stl", "kind": "implant", "opacity": 0, "visible": false, "name": "16"}
		]
	}`)
	c, err := config.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	bone, impl := c.Meshes[0].ToScene(), c.Meshes[1].ToScene()

	if bone.ID != "bone" || bone.Name != "bone" || bone.Color != config.DefaultColors[scene.KindBone] ||
		bone.Opacity != 0.7 || !bone.Visible || bone.Shading != scene.ShadingSmooth {
		t.Errorf("bone: %+v", bone)
	}
	// Явные opacity: 0 и visible: false не должны перетираться дефолтами.
	if impl.Opacity != 0 || impl.Visible || impl.Shading != scene.ShadingCrease || impl.CreaseDeg != config.DefaultCreaseDeg {
		t.Errorf("implant: %+v", impl)
	}
	if got := c.Path("Bone.stl"); got != filepath.Join(filepath.Dir(p), "Bone.stl") {
		t.Errorf("Path = %s", got)
	}
}

func TestValidate(t *testing.T) {
	p := writeConfig(t, `{
		"title": "",
		"meshes": [
			{"file": "a.stl", "kind": "tooth", "color": "red", "opacity": 2},
			{"file": "b.stl", "kind": "bone", "id": "a"},
			{"file": "a.stl", "kind": "bone"}
		]
	}`)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
	// Все проблемы – в одной ошибке.
	for _, s := range []string{"title", "tooth", "red", "opacity", "a"} {
		if !strings.Contains(err.Error(), s) {
			t.Errorf("в ошибке нет %q: %v", s, err)
		}
	}
}

func TestUnknownField(t *testing.T) {
	p := writeConfig(t, `{"title": "x", "meshes": [{"file": "a.stl", "kind": "bone", "colour": "#ffffff"}]}`)
	if _, err := config.Load(p); err == nil {
		t.Fatal("опечатка в ключе должна быть ошибкой")
	}
}
