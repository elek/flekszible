package processor

import (
	"github.com/elek/flekszible/pkg/data"
)

type Processor interface {
	data.Visitor

	Before(ctx *RenderContext)
	After(ctx *RenderContext)

	BeforeResource(*data.Resource)
	AfterResource(*data.Resource)
	//restrict processor execution only for one file
	OnlyForFiles(string)
}

type DefaultProcessor struct {
	data.DefaultVisitor
	Type            string
	File            string
	CurrentResource *data.Resource
}

func (processor *DefaultProcessor) Before(ctx *RenderContext) {}
func (processor *DefaultProcessor) After(ctx *RenderContext)  {}


func (p *DefaultProcessor) BeforeResource(resource *data.Resource) {
	p.CurrentResource = resource
}

func (p *DefaultProcessor) AfterResource(*data.Resource) {
	p.CurrentResource = nil
}

func (p *DefaultProcessor) OnlyForFiles(file string) {
	p.File = file
}