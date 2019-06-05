package flekszible

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/processor"
)

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
