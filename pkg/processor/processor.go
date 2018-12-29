package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/sirupsen/logrus"
	"regexp"
)

type Processor interface {
	data.Visitor

	Before(ctx *RenderContext, resources []data.Resource)
	After(ctx *RenderContext, resources []data.Resource)

	BeforeResource(*data.Resource)
	AfterResource(*data.Resource)
	//restrict processor execution only for one file
	OnlyForFiles(string)
	//returns true if the processor should be used for the specific resource
	Valid(resource data.Resource) bool
}

type DefaultProcessor struct {
	data.DefaultVisitor
	Type            string
	File            string
	CurrentResource *data.Resource
}

func (processor *DefaultProcessor) Valid(resource data.Resource) bool {
	if processor.File == "" {
		return true
	}
	r, err := regexp.Compile("^" + processor.File + ".*$")
	if err != nil {
		logrus.Warn("Invalid regular expression", err)
	}
	return r.Match([]byte(resource.Filename))
}

func (processor *DefaultProcessor) Before(ctx *RenderContext, resources []data.Resource) {}
func (processor *DefaultProcessor) After(ctx *RenderContext, resources []data.Resource)  {}


func (p *DefaultProcessor) BeforeResource(resource *data.Resource) {
	p.CurrentResource = resource
}

func (p *DefaultProcessor) AfterResource(*data.Resource) {
	p.CurrentResource = nil
}

func (p *DefaultProcessor) OnlyForFiles(file string) {
	p.File = file
}