package processor

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//the instatiated transformations for a specific run
type TransformationRepository struct {
	Processors []Processor
}

func (repository *TransformationRepository) Append(processor Processor) {
	repository.Processors = append(repository.Processors, processor)
}

func (repository *TransformationRepository) InsertToBeginning(processor Processor) {
	repository.Processors = append([]Processor{processor}, repository.Processors...)
}

func (repository *TransformationRepository) AppendAll(processors []Processor) {
	repository.Processors = append(processors, repository.Processors...)
}

//read processor definitions from a file (./definitions/xxx)
func (registry *ProcessorTypes) ReadProcessorDefinitionFile(filename string) ([]Processor, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return registry.ReadProcessorDefinition(data)
}

func (registry *ProcessorTypes) CreateTransformation(processorTypeName string) (Processor, error) {
	return registry.CreateTransformationWithConfig(processorTypeName, &yaml.MapSlice{})
}

func (registry *ProcessorTypes) CreateTransformationWithConfig(processorTypeName string, config *yaml.MapSlice) (Processor, error) {
	if definition, ok := registry.TypeMap[strings.ToLower(processorTypeName)]; ok {
		processor, err := definition.Factory(config)
		if err != nil {
			return nil, err
		}
		return processor, nil
	} else {
		if optional, found := config.Get("optional"); found && optional.(bool) {
			return nil, nil
		}
		logrus.Error("Unknown processor type: " + processorTypeName)
		logrus.Info("Available processor types:")
		for k := range registry.TypeMap {
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
func (registry *ProcessorTypes) ParseTransformations(inputDir string) ([]Processor, error) {
	result := make([]Processor, 0)
	mixinDir := path.Join(inputDir, "transformations")
	if _, err := os.Stat(mixinDir); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(mixinDir)
		if err != nil {
			return result, err
		}
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
				fullPath := path.Join(mixinDir, file.Name())
				processors, err := registry.ReadProcessorDefinitionFile(fullPath)
				if err != nil {
					return result, errors.Wrap(err, "Processor configuration can't be loaded from "+fullPath)
				}
				result = append(result, processors...)
			}
		}
	}
	return result, nil
}
