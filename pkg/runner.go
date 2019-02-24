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
	nodes := context.ListResourceNodes()
	for _, node := range nodes {
		for _, resource := range node.Resources {
			table.AddRow(resource.Name(), resource.Kind())
		}
	}
	fmt.Println("Detected resources:")
	fmt.Println(table.Render())
}


func ListProcessor(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}
	table := termtables.CreateTable()

	for name, definition := range processor.ProcessorTypeRegistry.TypeMap {
		table.AddRow(name, definition.Metadata.Description)
	}
	fmt.Println(table.Render())

}

func ShowProcessor(context *processor.RenderContext, command string) {
	err := context.Init()
	if err != nil {
		panic(err)
	}

	if procDefinition, found := processor.ProcessorTypeRegistry.TypeMap[command]; found {
		fmt.Println("")
		fmt.Println("### " + command)
		fmt.Println()
		fmt.Println(procDefinition.Metadata.Description)
		fmt.Println()
		fmt.Println("#### Parameters")
		fmt.Println("")
		table := termtables.CreateTable()
		table.AddHeaders("name", "default", "description")
		for _, parameter := range procDefinition.Metadata.Parameter {
			table.AddRow(parameter.Name, parameter.Default, parameter.Description)
		}
		fmt.Println(table.Render())
		fmt.Println()
		fmt.Println(procDefinition.Metadata.Doc)

	} else {
		fmt.Println("No such processor definition: " + command)
	}

}

func addSourceToTable(manager *data.SourceCacheManager, table *termtables.Table, source data.Source) {
	typ, value := source.ToString()
	path, _ := source.GetPath(manager, "")
	table.AddRow(typ, value, path)
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

	nodes := context.ListResourceNodes()

	sourceSet := make(map[string]bool)
	id, _ := context.RootResource.Origin.GetPath(&cacheManager, "")
	sourceSet[id] = true

	for _, node := range nodes {
		id, _ := node.Origin.GetPath(&cacheManager, "")
		if _, hasKey := sourceSet[id]; !hasKey {
			addSourceToTable(&cacheManager, table, node.Origin)
			sourceSet[id] = true
		}

	}

	fmt.Println("Detected sources:")
	fmt.Println(table.Render())
}

func SearchComponent(context *processor.RenderContext) {
}

func ListApp(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}

	table := termtables.CreateTable()
	table.AddHeaders("dir")

	nodes := context.ListResourceNodes()
	for _, node := range nodes {
		table.AddRow(node.Dir)
	}
	fmt.Println("Detected components (dirs):")
	fmt.Println(table.Render())
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
