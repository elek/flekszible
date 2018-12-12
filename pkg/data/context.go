package data

import (
	"path"
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

func (ctx *RenderContext) ParseDir(dir string) {
	configFile := path.Join(dir, "flekszible.yaml")
	conf, err := ReadConfiguration(configFile)
	if err != nil {
		panic(err)
	}
	for _, importDir := range conf.Import {
		absDir := path.Join(dir, importDir)
		ctx.InputDir = append([]string{absDir}, ctx.InputDir...)
		ctx.ParseDir(absDir)
	}
}

func (ctx *RenderContext) ReadConfigs() {
	ctx.Conf = Configuration{
		Import: make([]string, 0),
	}

	for _, dir := range ctx.InputDir {
		ctx.ParseDir(dir)
	}
}


