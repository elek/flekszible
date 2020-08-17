package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
	"github.com/pkg/errors"
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

func (processor *Change) BeforeResource(resource *data.Resource) error {
	var re = regexp.MustCompile(processor.Pattern)
	if !processor.Trigger.active(resource) {
		return nil
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
				return errors.Wrap(err, "Invalid replacement value: number can be replaced only with number and not: "+processor.Replacement)
			}
		case string:
			key.Value = re.ReplaceAllString(key.Value.(string), processor.Replacement)

		default:
			return errors.New("Invalid replacement only string or int can be replaced: " + result.Path.ToString())
		}

	}
	return nil
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Change",
			Description: "Replace existing value literal in the yaml struct",
			Parameters: []ProcessorParameter{
				PathParameter,
				TriggerParameter,
				{
					Name:        "pattern",
					Description: "Regular expression to test the existing value. Value will be changed only if matches.",
					Default:     ".*",
				},
				{
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
