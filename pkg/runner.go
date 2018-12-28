package pkg

import (
	"github.com/elek/flekszible/pkg/processor"
)

func Run(context *processor.RenderContext) {
	transformations, err := context.ReadConfigs()
	if err != nil {
		panic(err)
	}
	context.LoadDefinitions()
	processors := AddInternalTransformations(context, transformations)
	processor.Generate(processors, context)
}

func AddInternalTransformations(context *processor.RenderContext,
	preImportTransformations map[string][]byte) processor.ProcessorRepository {
	if len(context.ImageOverride) > 0 {
		context.ProcessorRepository.Append(&processor.Image{
			Image: context.ImageOverride,
		})
	}
	if len(context.Namespace) > 0 {
		context.ProcessorRepository.Append(&processor.Namespace{
			Namespace: context.Namespace,
		})
	}

	//initialize all the transformations from the input dir structure
	for _, directory := range context.InputDir {
		context.ProcessorRepository.ParseProcessors(directory)
	}
	if context.Mode == "helm" {
		context.ProcessorRepository.Append(&processor.HelmDecorator{})
		context.ProcessorRepository.Append(&processor.HelmWriter{
			OutputDir: context.OutputDir,
		})
	}
	if context.Mode == "k8s" {
		context.ProcessorRepository.Append(&processor.K8sWriter{})
	}
	return context.ProcessorRepository
}
