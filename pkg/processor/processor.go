package processor

import (
	"github.com/elek/flekszible/pkg/data"
)

type Processor interface {
	data.Visitor

	Before(ctx *RenderContext, resources []data.Resource)
	After(ctx *RenderContext, resources []data.Resource)

	BeforeResource(*data.Resource)
	AfterResource(*data.Resource)
}

type DefaultProcessor struct {
	data.DefaultVisitor
	Type            string
	File            string
	CurrentResource *data.Resource
}

func (processor *DefaultProcessor) Before(ctx *RenderContext, resources []data.Resource) {}
func (processor *DefaultProcessor) After(ctx *RenderContext, resources []data.Resource)  {}

func (p *DefaultProcessor) BeforeResource(resource *data.Resource) {
	p.CurrentResource = resource
}

func (p *DefaultProcessor) AfterResource(*data.Resource) {
	p.CurrentResource = nil
}
