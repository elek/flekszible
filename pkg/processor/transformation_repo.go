package processor

import (
	"errors"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//the instatiated transformations for a specific run
type TransformationRepository struct {
	Processors []Processor
}

func (repository *TransformationRepository) Append(processor Processor) {
	repository.Processors = append(repository.Processors, processor)
}

func (repository *TransformationRepository) AppendAll(processors []Processor) {
	repository.Processors = append(repository.Processors, processors...)
}

//read processor definitions from a file (./definitions/xxx)
func ReadProcessorDefinitionFile(filename string) ([]Processor, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ReadProcessorDefinition(data)
}

func CreateTransformation(processorTypeName string, config *yaml.MapSlice) (Processor, error) {
	if definition, ok := ProcessorTypeRegistry.TypeMap[strings.ToLower(processorTypeName)]; ok {
		processor, err := definition.Factory(config)
		if err != nil {
			return nil, err
		}
		return processor, nil
	} else {
		logrus.Error("Unknown processor type: " + processorTypeName)
		logrus.Info("Available processor types:")
		for k := range ProcessorTypeRegistry.TypeMap {
			logrus.Info(k)
		}

		return nil, errors.New("Unknown processor: " + processorTypeName)
	}

}

func CreateProcessorRepository() *TransformationRepository {
	return &TransformationRepository{
		Processors: make([]Processor, 0),
	}
}

//read transformations from ./transformations/... files
func (repository *TransformationRepository) ParseTransformations(inputDir string) {
	mixinDir := path.Join(inputDir, "transformations")
	if _, err := os.Stat(mixinDir); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(mixinDir)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
				fullPath := path.Join(mixinDir, file.Name())
				logrus.Info("Loading processor configuration from " + fullPath)
				processors, err := ReadProcessorDefinitionFile(fullPath)
				if err != nil {
					logrus.Error("Processor configuration can't be loaded from " + fullPath + " " + err.Error())
				}
				repository.AppendAll(processors)
			}
		}
	}
}
