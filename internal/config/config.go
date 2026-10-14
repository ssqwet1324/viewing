// Package config читает JSON-конфиг генератора (пример: testdata/sample/config.json,
// схема: docs/config.schema.json).
package config

import (
	"errors"

	"github.com/ssqwet1324/viewing/internal/scene"
)

type Config struct {
	Schema string `json:"$schema,omitempty"` // для автодополнения в редакторе, игнорируется
	Title  string `json:"title"`
	Meshes []Mesh `json:"meshes"`

	dir string // каталог конфига, от него резолвятся относительные пути
}

// Mesh – описание одного STL. Необязательные поля сделаны указателями: так «не задано»
// отличается от нуля (opacity: 0 и visible: false – валидные значения).
type Mesh struct {
	File      string         `json:"file"`
	ID        string         `json:"id,omitempty"`   // по умолчанию имя файла без расширения, в нижнем регистре
	Name      string         `json:"name,omitempty"` // по умолчанию = ID
	Kind      scene.Kind     `json:"kind"`
	Color     string         `json:"color,omitempty"` // #rrggbb, по умолчанию DefaultColors[Kind]
	Opacity   *float32       `json:"opacity,omitempty"`
	Visible   *bool          `json:"visible,omitempty"`
	Shading   scene.Shading  `json:"shading,omitempty"` // по умолчанию DefaultShading(Kind)
	CreaseDeg *float32       `json:"crease_deg,omitempty"`
	Transform *scene.Matrix4 `json:"transform,omitempty"`
}

var DefaultColors = map[scene.Kind]string{
	scene.KindAnatomy:  "#d9c7b0",
	scene.KindBone:     "#e8e0d0",
	scene.KindImplant:  "#9aa6b4",
	scene.KindAbutment: "#c9a96e",
	scene.KindSleeve:   "#5f9fd6",
	scene.KindWaxup:    "#f3efe2",
}

// DefaultOpacity – кость полупрозрачная, чтобы импланты были видны сразу.
var DefaultOpacity = map[scene.Kind]float32{
	scene.KindBone: 0.7,
}

const DefaultCreaseDeg = 30

// DefaultShading: для металла и пластика (implant/abutment/sleeve) – crease, для остального – smooth.
func DefaultShading(k scene.Kind) scene.Shading {
	switch k {
	case scene.KindImplant, scene.KindAbutment, scene.KindSleeve:
		return scene.ShadingCrease
	default:
		return scene.ShadingSmooth
	}
}

var errTODO = errors.New("config: не реализовано")

// Load: os.ReadFile → json.Decoder с DisallowUnknownFields (опечатка в ключе должна быть
// ошибкой, а не молча игнорироваться) → dir = filepath.Dir(path) → applyDefaults → Validate.
func Load(path string) (*Config, error) {
	// TODO
	return nil, errTODO
}

// applyDefaults заполняет незаданные поля значениями по умолчанию (см. комментарии к Mesh).
func (c *Config) applyDefaults() {
	// TODO
}

// Validate собирает все ошибки сразу через errors.Join, а не останавливается на первой.
// Проверяет: title не пуст; meshes не пуст; file задан; kind из scene.Kinds;
// color – #rrggbb; opacity в [0,1]; shading известен; id уникальны и подходят под [a-z0-9_-]+.
func (c *Config) Validate() error {
	// TODO
	return errTODO
}

// Path резолвит путь из конфига: абсолютный возвращается как есть, относительный – от каталога конфига.
func (c *Config) Path(p string) string {
	// TODO
	return p
}

// ToScene превращает описание в scene.Mesh (без geom – его проставит вызывающий код).
// Вызывать после applyDefaults.
func (m Mesh) ToScene() scene.Mesh {
	// TODO
	panic("config.ToScene: не реализовано")
}
