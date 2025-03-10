package main

import (
	"context"
	"fmt"
	"io"
	golog "log"
	"net/http"
	"os"

	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2fonts"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

func compileAndRenderSketch(code string) ([]byte, error) {
	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return []byte{}, err
	}
	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		ThemeID: &d2themescatalog.GrapeSoda.ID,
		Sketch:  go2.Pointer(true),
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: func(engine string) (d2graph.LayoutGraph, error) {
			return d2dagrelayout.DefaultLayout, nil
		},
		Ruler:      ruler,
		FontFamily: go2.Pointer(d2fonts.HandDrawn),
	}
	ctx := log.WithDefault(context.Background())
	diagram, _, err := d2lib.Compile(ctx, code, compileOpts, renderOpts)
	if err != nil {
		return []byte{}, err
	}
	out, err := d2svg.Render(diagram, renderOpts)
	return out, err
}

func compileAndRender(code string) ([]byte, error) {
	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return []byte{}, err
	}
	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		ThemeID: &d2themescatalog.GrapeSoda.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: func(engine string) (d2graph.LayoutGraph, error) {
			return d2dagrelayout.DefaultLayout, nil
		},
		Ruler: ruler,
	}
	ctx := log.WithDefault(context.Background())
	diagram, _, err := d2lib.Compile(ctx, code, compileOpts, renderOpts)
	if err != nil {
		return []byte{}, err
	}
	out, err := d2svg.Render(diagram, renderOpts)
	return out, err
}

func generate(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	var code string
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
	} else {
		code = string(bytes)
	}
	out, err := compileAndRenderSketch(code)
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
	} else {
		w.Write(out)
	}
}

func view(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("index.html")
	if err == nil {
		w.Write(body)
	} else {
		w.Write([]byte("Error occurred: " + err.Error()))
	}
}

func main() {
	http.HandleFunc("/generate", generate)
	http.HandleFunc("/view", view)
	golog.Fatal(http.ListenAndServe(":8080", nil))
}
