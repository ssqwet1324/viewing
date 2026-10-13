// Package container пишет и читает бинарные блобы: magic | version | hdrLen | JSON | body.
// Подробности формата лежат в docs/FORMAT.md.
package container

import (
	"encoding/json"
	"errors"
	"io"
)

// Magic – 4 байта в начале блоба.
type Magic [4]byte

var (
	MagicScene  = Magic{'D', 'V', 'S', 'C'}
	MagicVolume = Magic{'D', 'V', 'V', 'L'}
)

// prefixLen – magic + version + hdrLen.
const prefixLen = 12

// BodyAlign – тело начинается со смещения, кратного BodyAlign. Заголовок добивается пробелами.
const BodyAlign = 8

var (
	ErrBadMagic  = errors.New("container: неизвестный magic")
	ErrTruncated = errors.New("container: данные обрезаны")
	errTODO      = errors.New("container: не реализовано")
)

// Write сериализует header в JSON, добивает пробелами до выравнивания и пишет блоб в w.
func Write(w io.Writer, magic Magic, version uint32, header any, body []byte) error {
	// TODO: json.Marshal → pad := (BodyAlign - (prefixLen+len(js))%BodyAlign) % BodyAlign
	// Порядок записи: magic, version, hdrLen (с учётом pad), js, pad пробелов, body.
	return errTODO
}

// File – разобранный блоб. Header и Body ссылаются на исходный срез, без копий.
type File struct {
	Magic   Magic
	Version uint32
	Header  json.RawMessage
	Body    []byte
}

// Parse разбирает блоб. Проверяет длины (ErrTruncated) и что Header – валидный JSON.
// Magic не проверяет, это делает вызывающий код через Expect.
func Parse(data []byte) (File, error) {
	// TODO
	return File{}, errTODO
}

// Expect проверяет magic и версию: ErrBadMagic или ошибка о неподдерживаемой версии.
func (f File) Expect(magic Magic, version uint32) error {
	// TODO
	return errTODO
}
