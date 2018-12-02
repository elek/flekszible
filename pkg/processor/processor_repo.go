package processor

import (
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"reflect"
)

var ProcessorTypeRegistry ProcessorTypes

type ProcessorTypes struct {
	TypeMap map[string]reflect.Type
}

func (pt *ProcessorTypes) Add(processor Processor) {
	if pt.TypeMap == nil {
		pt.TypeMap = make(map[string]reflect.Type)
	}
	name := reflect.TypeOf(processor).Elem().Name()
	pt.TypeMap[name] = reflect.TypeOf(processor).Elem()
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

func ReadProcessorDefinitionFile(filename string) ([]Processor, error) {
	processors := make([]Processor, 0)
	processorsConfigs := make([]yaml.MapSlice, 0)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &processorsConfigs)
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

func ReadProcessorConfig(processor *Processor, configYamlContent string) {

}

func CreateProcessor(processorTypeName string) interface{} {
	if processorType, ok := ProcessorTypeRegistry.TypeMap[processorTypeName]; ok {
		value := reflect.New(processorType)
		processor := value.Interface()
		return processor
	} else {
		logrus.Error("Unknown processor type: " + processorTypeName)
		logrus.Info("Available processor types:")
		for k, _ := range ProcessorTypeRegistry.TypeMap {
			logrus.Info(k)
		}

		panic("Unknown processor")
	}

}

func CreateProcessorRepository() ProcessorRepository {
	return ProcessorRepository{
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
			if !file.IsDir() {
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
