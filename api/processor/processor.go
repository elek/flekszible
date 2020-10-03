package processor

import (
	"encoding/json"
	"fmt"
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"strconv"
)

type Processor interface {
	data.Visitor
	RegisterResources(ctx *RenderContext, node *ResourceNode) error
	Before(ctx *RenderContext, node *ResourceNode) error
	After(ctx *RenderContext, node *ResourceNode) error
	GetType() string
	BeforeResource(resource *data.Resource, location *ResourceNode) error
	AfterResource(resource *data.Resource, location *ResourceNode) error
	GetScope() string
	SetScope(scope string)
	ToString() string
}

type DefaultProcessor struct {
	data.DefaultVisitor
	Type            string
	Scope           string
	File            string
	CurrentResource *data.Resource
}

func (processor *DefaultProcessor) ToString() string {
	return "___"
}
func (processor *DefaultProcessor) RegisterResources(ctx *RenderContext, node *ResourceNode) error {
	return nil
}

func (processor *DefaultProcessor) GetType() string {
	return processor.Type
}
func (processor *DefaultProcessor) Before(ctx *RenderContext, node *ResourceNode) error {
	return nil
}
func (processor *DefaultProcessor) After(ctx *RenderContext, node *ResourceNode) error {
	return nil
}
func (processor *DefaultProcessor) GetScope() string {
	return processor.Scope
}
func (processor *DefaultProcessor) SetScope(scope string) {
	processor.Scope = scope
}
func (p *DefaultProcessor) BeforeResource(resource *data.Resource, location *ResourceNode) error {
	p.CurrentResource = resource
	return nil
}

func (p *DefaultProcessor) AfterResource(resource *data.Resource, location *ResourceNode) error {
	p.CurrentResource = nil
	return nil
}

func configureProcessorFromYamlFragment(processor Processor, config *yaml.MapSlice) (Processor, error) {
	//remove built-in keys
	clean := config.Remove("type")
	clean = clean.Remove("scope")
	config = &clean
	processorConfigYaml, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	err = yaml.UnmarshalStrict(processorConfigYaml, processor)
	if err != nil {
		return nil, err
	}
	return processor, nil
}

type ToStringBuilder struct {
	content    string
	parameters bool
}

func CreateToString(name string) *ToStringBuilder {
	return &ToStringBuilder{content: name}
}

func (builder *ToStringBuilder) AddBool(key string, value bool) *ToStringBuilder {
	if value {
		return builder.Add(key, "true")
	} else {
		return builder.Add(key, "false")
	}
}
func (builder *ToStringBuilder) Add(key string, value string) *ToStringBuilder {
	if !builder.parameters {
		builder.content += ":"
		builder.parameters = true
	} else {
		builder.content += ","
	}
	builder.content = builder.content + key + "=" + value
	return builder
}

func anyToString(value interface{}) string {
	switch typedValue := value.(type) {
	case []interface{}:
		valueString := ""
		for _, elem := range typedValue {
			if len(valueString) > 0 {
				valueString += ","
			}
			valueString += anyToString(elem)
		}
		return valueString
	case string:
		return typedValue
	case int:
		return strconv.Itoa(typedValue)
	case yaml.MapSlice:
		rawYaml, err := yaml.Marshal(typedValue)
		if err != nil {
			return "unparesable"
		}
		rawData := make(map[string]interface{})
		err = yaml.Unmarshal(rawYaml, &rawData)
		if err != nil {
			return "unparesable"
		}
		jsonData, err := json.Marshal(rawData)
		if err != nil {
			return "unparesable"
		}
		return string(jsonData)
	default:
		return fmt.Sprintf("?%T?", typedValue)
	}
}

func (builder *ToStringBuilder) AddValue(key string, value interface{}) *ToStringBuilder {
	return builder.Add(key, anyToString(value))
}

func (builder *ToStringBuilder) Build() string {
	return builder.content
}
