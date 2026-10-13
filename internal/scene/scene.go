// Package scene описывает JSON-заголовок сцены (DVSC) и собирает её бинарное тело.
package scene

import (
	"bytes"
	"errors"

	"github.com/ssqwet1324/viewing/internal/mesh"
)

// Version – версия блоба DVSC.
const Version = 1

type Kind string

const (
	KindAnatomy  Kind = "anatomy"
	KindBone     Kind = "bone"
	KindImplant  Kind = "implant"
	KindAbutment Kind = "abutment"
	KindSleeve   Kind = "sleeve"
	KindWaxup    Kind = "waxup"
)

// Kinds – все допустимые значения в порядке отображения во вьюере.
var Kinds = []Kind{KindAnatomy, KindBone, KindImplant, KindAbutment, KindSleeve, KindWaxup}

type Shading string

const (
	ShadingSmooth Shading = "smooth"
	ShadingFlat   Shading = "flat"
	ShadingCrease Shading = "crease"
)

type IndexType string

const (
	IndexU16 IndexType = "u16"
	IndexU32 IndexType = "u32"
)

// Matrix4 – 4×4 row-major: трансляция лежит в [3], [7], [11].
type Matrix4 [16]float32

type Scene struct {
	Title  string `json:"title"`
	Units  string `json:"units"`
	Meshes []Mesh `json:"meshes"`
	Geoms  []Geom `json:"geoms"`
}

type Mesh struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Kind      Kind     `json:"kind"`
	Color     string   `json:"color"`
	Opacity   float32  `json:"opacity"`
	Visible   bool     `json:"visible"`
	Shading   Shading  `json:"shading"`
	CreaseDeg float32  `json:"crease_deg,omitempty"`
	Transform *Matrix4 `json:"transform,omitempty"`
	Geom      int      `json:"geom"`
}

type Geom struct {
	Origin    [3]float32 `json:"origin"`
	Scale     [3]float32 `json:"scale"`
	Vertices  int        `json:"vertices"`
	Triangles int        `json:"triangles"`
	Index     IndexType  `json:"index"`
	PosOff    int        `json:"pos_off"`
	PosLen    int        `json:"pos_len"`
	IdxOff    int        `json:"idx_off"`
	IdxLen    int        `json:"idx_len"`
}

// SectionAlign – выравнивание секций внутри тела.
const SectionAlign = 4

var ErrBadGeom = errors.New("scene: ссылка на несуществующий geom")

// Builder накапливает геометрии и меши. Не потокобезопасен.
type Builder struct {
	scene Scene
	body  bytes.Buffer
}

func NewBuilder(title string) *Builder {
	return &Builder{scene: Scene{Title: title, Units: "mm"}}
}

// AddGeom дописывает позиции и индексы в тело и возвращает индекс geom.
//
// Порядок: выровнять body до SectionAlign → записать Pos (uint16 LE) → выровнять →
// записать Indices как u16 (если NumVertices <= 65536) или u32.
// Подсказка: binary.LittleEndian.AppendUint16 в заранее выделенный []byte, затем body.Write.
func (b *Builder) AddGeom(q mesh.Quantized) int {
	// TODO
	panic("scene.AddGeom: не реализовано")
}

// AddMesh добавляет меш. Если m.Geom вне диапазона, возвращает ErrBadGeom.
func (b *Builder) AddMesh(m Mesh) error {
	// TODO
	return errors.New("scene.AddMesh: не реализовано")
}

// Finish возвращает заголовок и тело. Повторно Builder не используется.
func (b *Builder) Finish() (Scene, []byte) {
	return b.scene, b.body.Bytes()
}
