package processor

import (
	"encoding/json"
	"fmt"
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	jsonpatch "github.com/evanphx/json-patch"
)

type Patch struct {
	DefaultProcessor
	Op    string
	Value interface{}
	Path  string
}

func (processor *Patch) BeforeResource(res *data.Resource) error {

	p := struct {
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value"`
	}{
		Op:    processor.Op,
		Path:  processor.Path,
		Value: processor.Value,
	}

	rawPath, err := json.Marshal([]interface{}{p})
	if err != nil {
		return err
	}
	jsonPatch, err := jsonpatch.DecodePatch(rawPath)
	if err != nil {
		return err
	}

	content := res.Content.ToMap()

	rawContent, err := json.Marshal(content)
	if err != nil {
		return err
	}
	fmt.Println(string(rawContent))
	modified, err := jsonPatch.Apply(rawContent)
	if err != nil {
		return err
	}
	transformed, err := data.ReadManifestString(modified)
	if err != nil {
		return err
	}
	res.Content = transformed

	return nil
}

func ActivatePatch(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "patch",
			Description: "Apply Json Patch (rfc6902)",
			Doc:         addDoc,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Patch{}, config)
		},
	})
}
