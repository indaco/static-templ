package generator

import (
	"html/template"
	"os"

	"github.com/nokacper24/templ-static-generator/internal/finder"
)

const outputScript = `// Code generated by TEMPL STATIC GENERATOR; DO NOT EDIT.
package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"log"
	{{- range .Imports }}
	"{{ . }}"
	{{- end }}
)

func main() {
	fileFuncs := []fileAndFunc{
	{{- range .FilesToGenerate }}
		{
			{{ .PackageName }}.{{ .FuncToCall }}().Render,
			"{{.FilePath}}",
		},
	{{- end }}
	}

	ctx := context.Background()

	for _, ff := range fileFuncs {
		if err := os.MkdirAll(filepath.Dir(ff.path), os.ModePerm); err != nil {
			log.Fatal("error creating dirs:", err)
		}
	
		file, err := os.Create(ff.path)
		if err != nil {
			log.Fatal("error creating file:", err)
		}
		defer file.Close()
		ff.function(ctx,file)
	}

}

type fileAndFunc struct {
	function func(ctx context.Context, w io.Writer) error
	path string
}
`

type InputData struct {
	Imports         []string
	FilesToGenerate []StringedData
}

type StringedData struct {
	FuncToCall  string
	FilePath    string
	PackageName string
}

func Generate(outputPath string, imports []string, funcs []finder.FileToGenerate) error {
	var stringed []StringedData
	for _, f := range funcs {
		stringed = append(stringed, StringedData{
			f.FunctionName,
			f.ToGenerate("web/pages", "dist"),
			f.PackageName,
		})
	}

	data := InputData{
		Imports:         imports,
		FilesToGenerate: stringed,
	}

	tmpl, err := template.New("output").Parse(outputScript)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}