package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/nokacper24/templ-static-generator/internal/finder"
	"github.com/nokacper24/templ-static-generator/internal/generator"
)

const (
	outputScriptDirPath  string = "temp"
	outputScriptFileName string = "templ_static_generate_script.go"
)

func main() {
	outputDir := "dist"
	inputDir := "web/pages"

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal("err creating dirs:", err)
	}

	modName, err := finder.FindModulePath()
	if err != nil {
		log.Fatal("Error finding module name:", err)
	}

	allFiles, err := finder.FindAllFiles(inputDir)
	if err != nil {
		log.Fatal("Error finding go dirs:", err)
	}

	groupedFiles := allFiles.ToGroupedFiles()

	funcs, err := finder.FindFuncsToCall(groupedFiles.TemplGoFiles)
	if err != nil {
		log.Fatal("Error finding funcs:", err)
	}

	importsMap := map[string]bool{}
	for _, f := range funcs {
		importsMap[f.DirPath()] = true
	}
	var importsSlice []string
	for imp := range importsMap {
		importsSlice = append(importsSlice, fmt.Sprintf("%s/%s", modName, imp))
	}

	err = os.RemoveAll(outputDir)
	if err != nil {
		log.Fatal("err removing files", err)
	}

	err = os.MkdirAll(outputScriptDirPath, os.ModePerm)
	if err != nil {
		log.Fatal("err creating temp dir:", err)
	}

	err = generator.Generate(getOutputScriptPath(), importsSlice, funcs)
	if err != nil {
		log.Fatal("err generating script", err)
	}

	cmd := exec.Command("go", "run", getOutputScriptPath())
	err = cmd.Start()
	if err != nil {
		log.Fatal("err starting script", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal("err running script", err)
	}

	err = os.Remove(getOutputScriptPath())
	if err != nil {
		log.Fatal("err removing enerated script file", err)
	}
}

func getOutputScriptPath() string {
	return fmt.Sprintf("%s/%s", outputScriptDirPath, outputScriptFileName)
}