// Package htmlout собирает итоговый самодостаточный HTML: шаблон + gzip/base64-блобы + вьюер.
package htmlout

import (
	_ "embed"
	"errors"
	"io"
)

// viewerJS – бандл вьюера, собирается `task viewer` (esbuild) и коммитится в репо.
//
//go:embed assets/viewer.js
var viewerJS []byte

// pageTmpl – шаблон страницы с inline-загрузчиком.
//
//go:embed page.html.tmpl
var pageTmpl string

// ID блобов. Должны совпадать с viewer/src/ids.ts.
const (
	IDScene  = "dv-scene"
	IDVolume = "dv-volume"
	IDViewer = "dv-viewer"
)

type Blob struct {
	ID   string
	Data []byte // сырые байты контейнера, сжимает Render
}

type Page struct {
	Title string
	Blobs []Blob // вьюер Render добавляет сам, последним
}

var errTODO = errors.New("htmlout: не реализовано")

// Render пишет страницу в w.
//
// Шаблон – text/template, а не html/template: html/template считает содержимое <script>
// JS-контекстом и экранирует base64 как JS-строку. Поэтому Title экранируем вручную
// (html.EscapeString). Base64 состоит из [A-Za-z0-9+/=], его экранировать не нужно.
// Данные для шаблона: {Title string; Blobs []struct{ID, Base64 string}}.
func Render(w io.Writer, p Page) error {
	// TODO
	return errTODO
}

// encodeBlob = base64.StdEncoding(gzip(data, gzip.BestCompression)).
// Не забудьте Close у gzip.Writer до чтения буфера, иначе хвост потока потеряется.
func encodeBlob(data []byte) (string, error) {
	// TODO
	return "", errTODO
}
