// dentview собирает автономный HTML-вьюер из STL (и позже DICOM) по JSON-конфигу.
//
//	dentview -c plan.json -o plan.html
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ssqwet1324/viewing/internal/config"
	"github.com/ssqwet1324/viewing/internal/container"
	"github.com/ssqwet1324/viewing/internal/htmlout"
	"github.com/ssqwet1324/viewing/internal/mesh"
	"github.com/ssqwet1324/viewing/internal/scene"
	"github.com/ssqwet1324/viewing/internal/stl"
)

func main() {
	if err := run(os.Args[1:], os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

// run отделён от main, чтобы его можно было вызвать из теста.
func run(args []string, log io.Writer) error {
	fs := flag.NewFlagSet("dentview", flag.ContinueOnError)
	fs.SetOutput(log)
	cfgPath := fs.String("c", "", "путь к JSON-конфигу (обязательно)")
	outPath := fs.String("o", "out.html", "куда писать HTML")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *cfgPath == "" {
		fs.Usage()
		return fmt.Errorf("не задан -c")
	}

	start := time.Now()
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		return err
	}

	b := scene.NewBuilder(cfg.Title)
	for _, mc := range cfg.Meshes {
		tris, err := stl.ReadFile(cfg.Path(mc.File))
		if err != nil {
			return err
		}
		q := mesh.Quantize(mesh.Weld(tris))
		sm := mc.ToScene()
		sm.Geom = b.AddGeom(q)
		if err := b.AddMesh(sm); err != nil {
			return err
		}
		fmt.Fprintf(log, "%-24s %8d тр. → %7d верш.\n", sm.ID, len(tris), q.NumVertices())
	}

	hdr, body := b.Finish()
	var sceneBuf bytes.Buffer
	if err := container.Write(&sceneBuf, container.MagicScene, scene.Version, hdr, body); err != nil {
		return err
	}

	// Пишем во временный файл и переименовываем: упавшая генерация не оставит битый HTML.
	tmp := *outPath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	page := htmlout.Page{
		Title: cfg.Title,
		Blobs: []htmlout.Blob{{ID: htmlout.IDScene, Data: sceneBuf.Bytes()}},
	}
	if err := htmlout.Render(f, page); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, *outPath); err != nil {
		return err
	}

	st, _ := os.Stat(*outPath)
	fmt.Fprintf(log, "готово: %s, %.1f МБ, %v\n", *outPath, float64(st.Size())/(1<<20), time.Since(start).Round(time.Millisecond))
	return nil
}
