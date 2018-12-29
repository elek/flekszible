package processor

import (
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
)

var ProcessorTypeRegistry ProcessorTypes

type ProcessorFactory = func() Processor

type ProcessorTypes struct {
	TypeMap map[string]ProcessorFactory
}

func (pt *ProcessorTypes) Add(processor Processor) {
	if pt.TypeMap == nil {
		pt.TypeMap = make(map[string]ProcessorFactory)
	}
	name := reflect.TypeOf(processor).Elem().Name()
	factory := func() Processor {
		processorType := reflect.TypeOf(processor).Elem()
		value := reflect.New(processorType)
		processor := value.Interface()
		return processor.(Processor)
	}
	pt.TypeMap[name] = factory
}

func (pt *ProcessorTypes) AddComposit(name string, factory ProcessorFactory) {
	logrus.Infof("Registring composit definition %s", name)
	pt.TypeMap[name] = factory
}

type ProcessorRepository struct {
	Processors []Processor
}

func (repository *ProcessorRepository) Append(processor Processor) {
	repository.Processors = append(repository.Processors, processor)
}

func (repository *ProcessorRepository) AppendAll(processors []Processor) {
	repository.Processors = append(repository.Processors, processors...)
}

//read processor definition file from a file
func ReadProcessorDefinitionFile(filename string) ([]Processor, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ReadProcessorDefinition(data)
}

//read processor definitions from raw yaml file
func ReadProcessorDefinition(data []byte) ([]Processor, error) {
	processors := make([]Processor, 0)
	processorsConfigs := make([]yaml.MapSlice, 0)
	err := yaml.Unmarshal(data, &processorsConfigs)
	if err != nil {
		return nil, err
	}
	for _, processorConfig := range processorsConfigs {
		typeName, ok := processorConfig.Get("type")
		if ok {
			proc := CreateProcessor(typeName.(string)).(Processor)
			processors = append(processors, proc)
			processorConfigYaml, err := yaml.Marshal(processorConfig)
			if err != nil {
				panic(err)
			}
			println(string(processorConfigYaml))
			err = yaml.Unmarshal(processorConfigYaml, proc)
			if err != nil {
				panic(err)
			}
		} else {
			panic("Type is missing from the config")
		}
	}
	return processors, nil
}

func CreateProcessor(processorTypeName string) interface{} {
	if factory, ok := ProcessorTypeRegistry.TypeMap[processorTypeName]; ok {
		processor := factory()
		return processor
	} else {
		logrus.Error("Unknown processor type: " + processorTypeName)
		logrus.Info("Available processor types:")
		for k := range ProcessorTypeRegistry.TypeMap {
			logrus.Info(k)
		}

		panic("Unknown processor")
	}

}

func CreateProcessorRepository() *ProcessorRepository {
	return &ProcessorRepository{
		Processors: make([]Processor, 0),
	}
}

func (repository *ProcessorRepository) ParseProcessors(inputDir string) {
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
