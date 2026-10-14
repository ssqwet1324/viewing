package htmlout_test

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/ssqwet1324/viewing/internal/htmlout"
)

var blobRe = regexp.MustCompile(`<script type="application/gzip" id="([a-z-]+)">([^<]*)</script>`)

// extractBlobs делает то же, что загрузчик в браузере: base64 → gunzip.
func extractBlobs(t *testing.T, html string) map[string][]byte {
	t.Helper()
	out := map[string][]byte{}
	for _, m := range blobRe.FindAllStringSubmatch(html, -1) {
		raw, err := base64.StdEncoding.DecodeString(m[2])
		if err != nil {
			t.Fatalf("%s: base64: %v", m[1], err)
		}
		zr, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			t.Fatalf("%s: gzip: %v", m[1], err)
		}
		data, err := io.ReadAll(zr)
		if err != nil {
			t.Fatalf("%s: gunzip: %v", m[1], err)
		}
		out[m[1]] = data
	}
	return out
}

func TestRender(t *testing.T) {
	scene := bytes.Repeat([]byte("DVSC payload "), 1000)
	var buf bytes.Buffer
	err := htmlout.Render(&buf, htmlout.Page{
		Title: `План <b>&</b> "тест"`,
		Blobs: []htmlout.Blob{{ID: htmlout.IDScene, Data: scene}},
	})
	if err != nil {
		t.Fatal(err)
	}
	html := buf.String()

	if !strings.Contains(html, "<title>План &lt;b&gt;&amp;&lt;/b&gt; &#34;тест&#34;</title>") {
		t.Error("title не экранирован")
	}
	blobs := extractBlobs(t, html)
	if !bytes.Equal(blobs[htmlout.IDScene], scene) {
		t.Error("блоб сцены не совпал")
	}
	if len(blobs[htmlout.IDViewer]) == 0 {
		t.Error("нет блоба вьюера")
	}
	// Вьюер должен идти после данных: к моменту его запуска все блобы уже в DOM.
	if strings.Index(html, `id="`+htmlout.IDViewer+`"`) < strings.Index(html, `id="`+htmlout.IDScene+`"`) {
		t.Error("вьюер раньше сцены")
	}
}
