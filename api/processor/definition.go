package processor

import (
	"bufio"
	"github.com/elek/flekszible/api/yaml"
	"github.com/sirupsen/logrus"
	"strings"
)

//split file to header(optional) and tranformation definition
func splitDefinitionFile(data []byte) ([]byte, []byte) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	head := ""
	body := ""
	headerPart := true
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" && headerPart {
			headerPart = false
		} else if headerPart {
			head += line + "\n"
		} else {
			body += line + "\n"
		}
	}
	return []byte(head), []byte(body)
}

func ParseDefinition(date []byte) (ProcessorMetadata, string, error) {
	return ProcessorMetadata{}, "", nil
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
			proc, err := CreateTransformationWithConfig(typeName.(string), &processorConfig)
			if err != nil {
				logrus.Error("Transformation can't be instantiated: " + typeName.(string) + " " + err.Error())
			} else if proc == nil {
				logrus.Info("Optional transformation depends on an unknown transformation type: " + typeName.(string) + " Additional import may be required to use optional features.")
			} else {
				if scope, found := processorConfig.Get("scope"); found {
					proc.SetScope(scope.(string))
				}

				processors = append(processors, proc)
			}
		} else {
			panic("Type is missing from the config")
		}
	}
	return processors, nil
}
