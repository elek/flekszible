package pkg

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/processor"
)

func Run(context *data.RenderContext) {
	context.ReadConfigs()

	processor.LoadDefinitions(context)
	processors := initProcessors(context)
	processors = filter(&processors, context.Mode)
	processor.Generate(processors, context)
}

func filter(repository *processor.ProcessorRepository, mode string) processor.ProcessorRepository {
	filtered := processor.ProcessorRepository{}
	for _, p := range repository.Processors {
			filtered.Append(p)

	}

	return filtered
}

func initProcessors(context *data.RenderContext) processor.ProcessorRepository {
	repository := processor.CreateProcessorRepository()
	repository.Append(&processor.Initializer{})

	if len(context.ImageOverride) > 0 {
		repository.Append(&processor.ImageSet{
			Image: context.ImageOverride,
		})
	}
	if len(context.Namespace) > 0 {
		repository.Append(&processor.Namespace{
			Namespace: context.Namespace,
		})
	}
	for _, directory := range context.InputDir {
		repository.ParseProcessors(directory)
	}
	if context.Mode == "helm" {
		repository.Append(&processor.HelmDecorator{})
		repository.Append(&processor.HelmWriter{
			OutputDir: context.OutputDir,
		})
	}
	if context.Mode == "k8s" {
		repository.Append(&processor.K8sWriter{})
	}
	return repository
}
