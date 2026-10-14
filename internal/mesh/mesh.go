// Package mesh превращает суп треугольников в индексированную сетку и квантует её в uint16.
package mesh

import (
	"github.com/ssqwet1324/viewing/internal/stl"
)

// Indexed – сетка с общими вершинами.
type Indexed struct {
	Positions []stl.Vec3
	Indices   []uint32 // len % 3 == 0, обход как в исходном STL
}

// Weld склеивает вершины с побитово равными координатами.
//
// Правила:
//   - Сравнение точное (map[stl.Vec3]uint32), без эпсилона. CAD-экспорт пишет общие вершины
//     одинаковыми float, а сравнение с допуском склеивает лишнее на мелких деталях (резьба импланта).
//   - -0 и +0 считаются одной вершиной (нормализуйте перед поиском в map).
//   - Порядок треугольников и обход внутри треугольника сохраняются.
//   - Вырожденные треугольники (после склейки два индекса совпали) выбрасываются.
//   - Новые вершины нумеруются в порядке первого появления. Это даёт детерминизм для golden-тестов.
func Weld(tris []stl.Triangle) Indexed {
	// TODO
	panic("mesh.Weld: не реализовано")
}

// Bounds возвращает min/max по каждой оси. Для пустой сетки – нули.
func (m Indexed) Bounds() (lo, hi stl.Vec3) {
	// TODO
	panic("mesh.Bounds: не реализовано")
}

// Quantized – позиции в uint16 относительно bbox.
type Quantized struct {
	Origin  [3]float32 // = Bounds().lo
	Scale   [3]float32 // (hi-lo)/65535 по оси; если hi == lo, то 1 (любое ненулевое)
	Pos     []uint16   // 3 на вершину, xyz подряд
	Indices []uint32
}

// MaxQ – максимальное значение кванта.
const MaxQ = 65535

// Quantize: q = round((p - origin) / scale), с clamp в [0, MaxQ].
// Считайте во float64, иначе на краях bbox бывает 65536.
func Quantize(m Indexed) Quantized {
	// TODO
	panic("mesh.Quantize: не реализовано")
}

// Vertex восстанавливает i-ю вершину: origin + q*scale.
func (q Quantized) Vertex(i int) stl.Vec3 {
	// TODO
	panic("mesh.Vertex: не реализовано")
}

// NumVertices – число вершин.
func (q Quantized) NumVertices() int { return len(q.Pos) / 3 }

// NumTriangles – число треугольников.
func (q Quantized) NumTriangles() int { return len(q.Indices) / 3 }
