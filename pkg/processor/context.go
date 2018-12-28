package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type RenderContext struct {
	OutputDir           string
	InputDir            []string
	Conf                data.Configuration
	Mode                string
	ImageOverride       string
	Namespace           string
	Resources           []data.Resource
	ProcessorRepository ProcessorRepository
}

func CreateRenderContext(mode string, inputDir string, outputDir string) *RenderContext {
	return &RenderContext{
		OutputDir:           outputDir,
		Mode:                mode,
		InputDir:            []string{inputDir,},
		ProcessorRepository: CreateProcessorRepository(),
	}
}

//read config file and follow the import path
func (ctx *RenderContext) ParseDir(dir string) (map[string][]byte, error) {
	transformationDefinitions := make(map[string][]byte)
	configFile := path.Join(dir, "flekszible.yaml")
	conf, err := data.ReadConfiguration(configFile)
	if err != nil {
		return transformationDefinitions, err
	}
	for _, importDefinition := range conf.Import {
		var importedDir string
		if !filepath.IsAbs(importDefinition.Path) {
			importedDir = path.Join(dir, importDefinition.Path)
		} else {
			importedDir = importDefinition.Path
		}
		ctx.InputDir = append([]string{importedDir}, ctx.InputDir...)
		transformationForImports, err := ctx.ParseDir(importedDir)
		for k, v := range transformationForImports {
			transformationDefinitions[k] = v
		}
		if err != nil {
			return transformationDefinitions, err
		}
		if importDefinition.Transformations != nil {
			rawTransformation, err := yaml.Marshal(importDefinition.Transformations)
			if err != nil {
				return transformationDefinitions, err
			}
			transformationDefinitions[path.Join(dir, importDefinition.Path)] = rawTransformation
		}

	}
	return transformationDefinitions, nil
}

//Parse the imports/configs from the specific input dirs
func (ctx *RenderContext) ReadConfigs() (map[string][]byte, error) {
	transformationDefinitions := make(map[string][]byte)
	ctx.Conf = data.Configuration{
		Import: make([]data.ImportConfiguration, 0),
	}

	for _, dir := range ctx.InputDir {
		tf, err := ctx.ParseDir(dir)
		for k, v := range tf {
			transformationDefinitions[k] = v
		}
		if err != nil {
			return transformationDefinitions, err
		}
	}
	return transformationDefinitions, nil
}

//Collect all the resource files from the input dirs
func (ctx *RenderContext) ReadResources() {
	resources := make([]data.Resource, 0)
	for _, inputDir := range ctx.InputDir {
		resources = append(resources, data.ReadResourcesFromDir(inputDir)...)
	}
	ctx.Resources = resources

}

//load transformation definitions from ./definitions dir (all dir)
func (ctx *RenderContext) LoadDefinitions() {
	for _, dir := range ctx.InputDir {
		defDir := path.Join(dir, "definitions")
		if _, err := os.Stat(defDir); !os.IsNotExist(err) {
			files, err := ioutil.ReadDir(defDir)
			if err != nil {
				logrus.Warningf("Can't load definition directory %s: %s", defDir, err.Error())
			}
			for _, file := range files {
				definitionFile := path.Join(defDir, file.Name())
				err := parseDefintion(definitionFile)
				if err != nil {
					logrus.Errorf("Can't parse the definition file %s: %s", definitionFile, err.Error())
				}
			}
		}
	}
}
