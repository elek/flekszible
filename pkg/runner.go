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

func addSourceToTable(manager *data.SourceCacheManager, table *termtables.Table, source data.Source) {
	typ, value := source.ToString()
	source.GetPath(manager, "")
	table.AddRow(typ, value)
}
func ListSources(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}

	cacheManager := data.SourceCacheManager{RootPath: context.RootResource.Dir}

	table := termtables.CreateTable()
	table.AddHeaders("source", "location", "path")

	addSourceToTable(&cacheManager, table, &data.EnvSource{})
	addSourceToTable(&cacheManager, table, &data.LocalSource{RelativeTo: context.RootResource.Dir})

	listSources(&cacheManager, context.RootResource, table)
	fmt.Println("Detected sources:")
	fmt.Println(table.Render())
}

func listSources(manager *data.SourceCacheManager, node *processor.ResourceNode, table *termtables.Table) {
	addSourceToTable(manager, table, node.Origin)
	for _, child := range node.Children {
		listSources(manager, child, table)
	}
}

func ListComponent(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}

	table := termtables.CreateTable()
	table.AddHeaders("component", "source")

	listComponent(context.RootResource, table)
	fmt.Println("Detected components (dirs):")
	fmt.Println(table.Render())
}

func listComponent(node *processor.ResourceNode, table *termtables.Table) {
	tpe, _ := node.Origin.ToString()
	table.AddRow(node.Dir, tpe)
	for _, child := range node.Children {
		listComponent(child, table)
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
