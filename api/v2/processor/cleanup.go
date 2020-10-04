package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Cleanup struct {
	DefaultProcessor
	ResourceOutputDir string
	All               bool
	cleanedDirs       map[string]bool
}

func (cleanup *Cleanup) Before(ctx *RenderContext, node *ResourceNode) error {
	cleanup.ResourceOutputDir = ctx.OutputDir
	return nil
}

func createOutputFileName(name string, kind string) string {
	return strings.ToLower(name) + "-" + strings.ToLower(kind) + ".yaml"
}
func (cleanup *Cleanup) createOutputPath(outputDir, name, kind string, destination string, destinationFile string) string {
	if destinationFile != "" {
		return path.Join(outputDir, destination, destinationFile)
	} else {
		return path.Join(outputDir, destination, createOutputFileName(name, kind))
	}
}

func (processor *Cleanup) After(ctx *RenderContext, node *ResourceNode) error {
	if processor.All {
		for dir := range processor.cleanedDirs {
			logrus.Info("Deleting all YAML files from the directory ", dir)
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, f := range files {
				if filepath.Ext(f.Name()) == ".yaml" && !f.IsDir() {
					err = os.Remove(path.Join(dir, f.Name()))
					if err != nil {
						logrus.Warn("Delete was unsuccessfull ", err)
					}
				}
			}
		}
	}
	return nil
}

func (cleanup *Cleanup) AfterResource(resource *data.Resource) error {
	outputDir := cleanup.ResourceOutputDir
	outputFile := cleanup.createOutputPath(outputDir, resource.Name(), resource.Kind(), resource.Destination, resource.DestinationFileName)
	if !cleanup.All {
		logrus.Info("Deleting generated file", outputFile)
		err := os.Remove(outputFile)
		if err != nil {
			logrus.Warn("Delete was unsuccessfull ", err)
		}
	} else {
		cleanup.cleanedDirs[filepath.Dir(outputFile)] = true
	}
	return nil
}

func CreateCleanup(outputDir string, All bool) Processor {
	cl := &Cleanup{ResourceOutputDir: outputDir, All: All}
	cl.cleanedDirs = make(map[string]bool)
	return cl
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "cleanup",
			Description: "Internal transformation to delete generated yaml files",
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Cleanup{}, config)
		},
	})
}
