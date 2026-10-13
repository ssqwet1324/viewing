package container_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ssqwet1324/viewing/internal/container"
)

type hdr struct {
	Title string    `json:"title"`
	Dims  [3]int    `json:"dims"`
	Vals  []float32 `json:"vals"`
}

func TestRoundTrip(t *testing.T) {
	// Разная длина заголовка, чтобы проверить все варианты добивки.
	for n := range 10 {
		want := hdr{Title: "Тест" + strings.Repeat("ё", n), Dims: [3]int{343, 246, 183}, Vals: []float32{0.32}}
		body := bytes.Repeat([]byte{1, 2, 3}, 100+n)

		var buf bytes.Buffer
		if err := container.Write(&buf, container.MagicScene, 7, want, body); err != nil {
			t.Fatal(err)
		}
		f, err := container.Parse(buf.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		if err := f.Expect(container.MagicScene, 7); err != nil {
			t.Fatal(err)
		}
		var got hdr
		if err := json.Unmarshal(f.Header, &got); err != nil {
			t.Fatal(err)
		}
		if got.Title != want.Title || got.Dims != want.Dims {
			t.Fatalf("header: %+v", got)
		}
		if !bytes.Equal(f.Body, body) {
			t.Fatal("body не совпало")
		}
		if off := buf.Len() - len(body); off%container.BodyAlign != 0 {
			t.Fatalf("тело начинается с %d, не кратно %d", off, container.BodyAlign)
		}
	}
}

func TestExpect(t *testing.T) {
	var buf bytes.Buffer
	if err := container.Write(&buf, container.MagicVolume, 1, struct{}{}, nil); err != nil {
		t.Fatal(err)
	}
	f, err := container.Parse(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Expect(container.MagicScene, 1); !errors.Is(err, container.ErrBadMagic) {
		t.Fatalf("err = %v, want ErrBadMagic", err)
	}
	if err := f.Expect(container.MagicVolume, 2); err == nil {
		t.Fatal("ожидалась ошибка версии")
	}
}

func TestParseTruncated(t *testing.T) {
	var buf bytes.Buffer
	if err := container.Write(&buf, container.MagicScene, 1, hdr{Title: "x"}, nil); err != nil {
		t.Fatal(err)
	}
	data := buf.Bytes()
	for _, n := range []int{0, 3, 11, 12, len(data) - 1} {
		if _, err := container.Parse(data[:n]); !errors.Is(err, container.ErrTruncated) {
			t.Errorf("len %d: err = %v, want ErrTruncated", n, err)
		}
	}
}

func FuzzParse(f *testing.F) {
	var buf bytes.Buffer
	_ = container.Write(&buf, container.MagicScene, 1, hdr{Title: "seed"}, []byte{1, 2})
	f.Add(buf.Bytes())
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = container.Parse(data) // не паниковать, в т.ч. на hdrLen = 0xFFFFFFFF
	})
}
