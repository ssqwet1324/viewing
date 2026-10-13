// Package stl читает STL (binary и ASCII) в «суп» треугольников без общих вершин.
package stl

import (
	"errors"
)

// Vec3 – точка в мм.
type Vec3 [3]float32

// Triangle – три вершины в порядке обхода из файла. Нормаль из STL не храним:
// она часто кривая, а вьюер всё равно пересчитывает.
type Triangle [3]Vec3

var (
	ErrEmpty     = errors.New("stl: файл не содержит треугольников")
	ErrTruncated = errors.New("stl: файл обрезан")
	errTODO      = errors.New("stl: не реализовано")
)

const (
	binHeaderLen = 80
	binRecordLen = 50 // normal 3×f32 + 3 вершины × 3×f32 + attr uint16
)

// ReadFile читает файл целиком и разбирает его через Parse.
// Файлы STL не больше ~100 МБ, поэтому стриминг не нужен.
func ReadFile(path string) ([]Triangle, error) {
	// TODO: os.ReadFile + Parse; ошибку оборачивать с именем файла: fmt.Errorf("%s: %w", path, err)
	return nil, errTODO
}

// Parse определяет формат и разбирает данные.
//
// Подвох: многие binary-STL пишут "solid ..." в 80-байтный заголовок, поэтому по префиксу
// формат определять нельзя. Надёжный признак binary: len(data) == 84 + 50*count,
// где count – uint32 LE по смещению 80. Если не совпало – пробуем ASCII.
func Parse(data []byte) ([]Triangle, error) {
	// TODO
	return nil, errTODO
}

// isBinary – проверка размера, описанная в Parse.
func isBinary(data []byte) bool {
	// TODO
	return false
}

// parseBinary разбирает binary-STL. Запись: normal(3×f32), v0, v1, v2 (по 3×f32), attr(u16).
// Подсказка: binary.LittleEndian.Uint32 + math.Float32frombits, без binary.Read (он медленный).
func parseBinary(data []byte) ([]Triangle, error) {
	// TODO
	return nil, errTODO
}

// parseASCII разбирает ASCII-STL. Структура facet/outer loop нам не важна: достаточно
// идти по токенам (bytes.Fields) и после каждого "vertex" читать три числа
// strconv.ParseFloat(s, 32). Число вершин должно делиться на 3.
func parseASCII(data []byte) ([]Triangle, error) {
	// TODO
	return nil, errTODO
}
