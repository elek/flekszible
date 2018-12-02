package processor

import (
	"github.com/elek/flekszible/pkg/data"
)

type Processor interface {
	data.Visitor

	Before(ctx *data.RenderContext)
	After(ctx *data.RenderContext)

	BeforeResource(*data.Resource)
	AfterResource(*data.Resource)
}

type DefaultProcessor struct {
	data.DefaultVisitor
	Type            string
	CurrentResource *data.Resource
}

func (processor *DefaultProcessor) Before(ctx *data.RenderContext) {}
func (processor *DefaultProcessor) After(ctx *data.RenderContext)  {}


func (p *DefaultProcessor) BeforeResource(resource *data.Resource) {
	p.CurrentResource = resource
}

func (p *DefaultProcessor) AfterResource(*data.Resource) {
	p.CurrentResource = nil
}

