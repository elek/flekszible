package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
)

type Processor interface {
	data.Visitor
	Before(ctx *RenderContext, resources []*data.Resource) error
	After(ctx *RenderContext, resources []*data.Resource) error

	BeforeResource(*data.Resource) error
	AfterResource(*data.Resource) error
	GetScope() string
	SetScope(scope string)
}

type DefaultProcessor struct {
	data.DefaultVisitor
	Type            string
	Scope           string
	File            string
	CurrentResource *data.Resource
}

func (processor *DefaultProcessor) Before(ctx *RenderContext, resources []*data.Resource) error {
	return nil
}
func (processor *DefaultProcessor) After(ctx *RenderContext, resources []*data.Resource) error {
	return nil
}
func (processor *DefaultProcessor) GetScope() string {
	return processor.Scope
}
func (processor *DefaultProcessor) SetScope(scope string) {
	processor.Scope = scope
}
func (p *DefaultProcessor) BeforeResource(resource *data.Resource) error {
	p.CurrentResource = resource
	return nil
}

func (p *DefaultProcessor) AfterResource(*data.Resource) error {
	p.CurrentResource = nil
	return nil
}

func configureProcessorFromYamlFragment(processor Processor, config *yaml.MapSlice) (Processor, error) {
	processorConfigYaml, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(processorConfigYaml, processor)
	if err != nil {
		return nil, err
	}
	return processor, nil
}