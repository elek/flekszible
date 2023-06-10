package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"strings"
)

type Substitute struct {
	DefaultProcessor
	Trigger Trigger
	From    string
	To      string
}

func (st *Substitute) ToString() string {
	return CreateToString("replace").
		Add("from", st.From).
		Add("to", st.To).
		Build()
}

func (st *Substitute) BeforeResource(resource *data.Resource) error {
	if !st.Trigger.active(resource) {
		return nil
	}

	target := data.GetStrings{}
	resource.Content.Accept(&target)
	for _, match := range target.Result {
		match.Value.Value = strings.ReplaceAll(match.Value.Value.(string), st.From, st.To)
	}
	return nil
}

func ActivateSubstitute(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Substitute",
			Description: "Substitute one string with an other EVERYWHERE in Yaml files.",
			Doc:         TriggerDoc,
			Parameters: []ProcessorParameter{
				TriggerParameter,
				{
					Name:        "from",
					Description: "A string which supposed to be replaced",
				},
				{
					Name:        "to",
					Description: "the replacement string",
				},
			},
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Substitute{}, config)
		},
	})
}
