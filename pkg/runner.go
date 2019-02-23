package pkg

import (
	"fmt"
	"github.com/apcera/termtables"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/processor"
)

func ListResources(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}
	table := termtables.CreateTable()
	table.AddHeaders("name", "kind")
	listResources(context.RootResource, table)
	fmt.Println("Detected resources:")
	fmt.Println(table.Render())
}

func listResources(node *processor.ResourceNode, table *termtables.Table) {
	for _, resource := range node.Resources {
		table.AddRow(resource.Name(), resource.Kind())
	}
	for _, child := range node.Children {
		listResources(child, table)
	}
}

func ListSources(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}

	table := termtables.CreateTable()
	table.AddHeaders("source")

	env := data.EnvSource{}
	table.AddRow(env.ToString())

	local := data.LocalSource{RelativeTo: context.RootResource.Dir}
	table.AddRow(local.ToString())

	listSources(context.RootResource, table)
	fmt.Println("Detected sources:")
	fmt.Println(table.Render())
}

func listSources(node *processor.ResourceNode, table *termtables.Table) {
	for _, source := range node.Source {
		table.AddRow(source.ToString())
	}
	for _, child := range node.Children {
		listSources(child, table)
	}
}

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
