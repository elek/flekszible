package processor

import (
	"bufio"
	"strings"

	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// split file to header(optional) and tranformation definition
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

// ReadProcessorDefinition reads processor definitions from raw yaml file.
func (registry *ProcessorTypes) ReadProcessorDefinition(data []byte) ([]Processor, error) {
	processors := make([]Processor, 0)
	processorsConfigs := make([]yaml.MapSlice, 0)
	err := yaml.Unmarshal(data, &processorsConfigs)
	if err != nil {
		return nil, err
	}
	for _, processorConfig := range processorsConfigs {
		typeName, ok := processorConfig.Get("type")
		if ok {
			proc, err := registry.CreateTransformationWithConfig(typeName.(string), &processorConfig)
			if err != nil {
				return processors, errors.Wrapf(err, "Transformation '%s' can't be instantiated", typeName.(string))
			} else if proc == nil {
				logrus.Debug("Optional transformation depends on an unknown transformation type: " + typeName.(string) + " Additional import may be required to use optional features.")
			} else {
				if scope, found := processorConfig.Get("scope"); found {
					proc.SetScope(scope.(string))
				}

				processors = append(processors, proc)
			}
		} else {
			return processors, errors.New("Type tag is missing from a processor definition")
		}
	}
	return processors, nil
}
