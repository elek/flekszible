package pkg

import (
	"github.com/elek/flekszible/pkg/processor"
)

func Run(context *processor.RenderContext, minikube bool) {
	err := context.Init()
	if err != nil {
		panic(err)
	}
	AddInternalTransformations(context, minikube)
	context.Render()
}

func AddInternalTransformations(context *processor.RenderContext, minikube bool) {
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
	if (minikube) {
		context.RootResource.ProcessorRepository.Append(&processor.DaemonToStatefulSet{})
		context.RootResource.ProcessorRepository.Append(&processor.PublishStatefulSet{})
	}
	if context.Mode == "k8s" {
		context.RootResource.ProcessorRepository.Append(&processor.K8sWriter{})
	}
}
