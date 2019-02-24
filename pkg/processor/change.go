package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
)

type Change struct {
	DefaultProcessor
	Trigger     Trigger
	Path        data.Path
	Pattern     string
	Replacement string
}

func (processor *Change) BeforeResource(resource *data.Resource) {
	var re = regexp.MustCompile(processor.Pattern)
	if !processor.Trigger.active(resource) {
		return
	}
	getter := data.SmartGetAll{Path: processor.Path}
	resource.Content.Accept(&getter)
	for _, result := range getter.Result {
		key := result.Value.(*data.KeyNode)
		switch key.Value.(type) {
		case int:
			intval, err := strconv.Atoi(processor.Replacement)
			key.Value = intval
			if err != nil {
				logrus.Error("Invalid replacement value: number can be replaced only with number and not: " + processor.Replacement)
			}
		case string:
			key.Value = re.ReplaceAllString(key.Value.(string), processor.Replacement)

		default:
			logrus.Error("Invalid replacement only string or int can be replaced: " + result.Path.ToString())
		}

	}
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Change",
			Description: "Replace existing value literal in the yaml struct.",
			Parameter: []ProcessorParameter{
				ProcessorParameter{
					Name:        "pattern",
					Description: "Regular expression to test the existing value. Value will be changed only if matches.",
					Default:     ".*",
				},
				ProcessorParameter{
					Name:        "replacement",
					Description: "The value to replace the field in case of match",
					Required:    true,
				},
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Change{Pattern: ".*"}, config)
		},
	})
}
