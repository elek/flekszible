package data

import (
	"path"
	"path/filepath"
)

type RenderContext struct {
	OutputDir     string
	InputDir      []string
	Conf          Configuration
	Mode          string
	ImageOverride string
	Namespace     string
	Resources     []Resource
}

//read config file and follow the import path
func (ctx *RenderContext) ParseDir(dir string) {
	configFile := path.Join(dir, "flekszible.yaml")
	conf, err := ReadConfiguration(configFile)
	if err != nil {
		panic(err)
	}
	for _, importDir := range conf.Import {
		var importedDir string
		if !filepath.IsAbs(importDir.Path) {
			importedDir = path.Join(dir, importDir.Path)
		} else {
			importedDir = importDir.Path
		}
		ctx.InputDir = append([]string{importedDir}, ctx.InputDir...)
		ctx.ParseDir(importedDir)
	}
}

func (ctx *RenderContext) ReadConfigs() {
	ctx.Conf = Configuration{
		Import: make([]ImportConfiguration, 0),
	}

	for _, dir := range ctx.InputDir {
		ctx.ParseDir(dir)
	}
}


