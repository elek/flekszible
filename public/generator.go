package flekszible

import (
	"github.com/elek/flekszible/public/data"
	"github.com/elek/flekszible/public/processor"
)

func Initialize(inputDir string) (*processor.RenderContext, error) {
	context := processor.CreateRenderContext("k8s", inputDir, "")
	err := context.Init()
	if err != nil {
		return nil, err
	}
	return context, nil
}

func Generate(inputDir string, processors []processor.Processor) ([]*data.Resource, error) {
	context := processor.CreateRenderContext("k8s", inputDir, "")
	err := context.Init()
	if err != nil {
		return nil, err
	}
	for _, processor := range processors {
		context.RootResource.ProcessorRepository.Append(processor)
	}
	context.Render()
	return context.Resources(), nil
}
