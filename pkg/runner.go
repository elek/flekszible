package pkg

import (
	"github.com/elek/flekszible/pkg/processor"
)

func Run(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}
	AddInternalTransformations(context)
	context.Render()
}

func AddInternalTransformations(context *processor.RenderContext) {
	if len(context.ImageOverride) > 0 {
		context.RootResource.ProcessorRepository.Append(&processor.Image{
			Image: context.ImageOverride,
		})
	}
	if len(context.Namespace) > 0 {
		context.RootResource.ProcessorRepository.Append(&processor.Namespace{
			Namespace: context.Namespace,
		})
	}

	if context.Mode == "helm" {
		context.RootResource.ProcessorRepository.Append(&processor.HelmDecorator{})
		context.RootResource.ProcessorRepository.Append(&processor.HelmWriter{
			OutputDir: context.OutputDir,
		})
	}
	if context.Mode == "k8s" {
		context.RootResource.ProcessorRepository.Append(&processor.K8sWriter{})
	}
}
