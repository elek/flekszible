package pkg

import (
	"fmt"
	"github.com/apcera/termtables"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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
	table.AddHeaders("name", "description")
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

func listUniqSources(context *processor.RenderContext) []data.Source {

	sources := make([]data.Source, 0)
	cacheManager := data.SourceCacheManager{RootPath: context.RootResource.Dir}

	sources = append(sources, &data.EnvSource{})
	sources = append(sources, &data.CurrentDir{CurrentDir: context.RootResource.Dir})

	nodes := context.ListResourceNodes()

	sourceSet := make(map[string]bool)
	id, _ := context.RootResource.Origin.GetPath(&cacheManager, "")
	sourceSet[id] = true

	for _, node := range nodes {
		for
		_, source := range node.Source {
			id, _ := source.GetPath(&cacheManager, "")
			if _, hasKey := sourceSet[id]; !hasKey {
				sources = append(sources, source)
				sourceSet[id] = true
			}
		}

	}

	return sources
}
func ListSources(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}

	cacheManager := data.SourceCacheManager{RootPath: context.RootResource.Dir}

	table := termtables.CreateTable()
	table.AddHeaders("source", "location", "path")

	for _, source := range listUniqSources(context) {
		typ, value := source.ToString()
		path, _ := source.GetPath(&cacheManager, "")
		table.AddRow(typ, value, path)
	}
	fmt.Println("Detected sources:")
	fmt.Println(table.Render())
}

func SearchComponent(context *processor.RenderContext) {
	err := context.Init()
	if err != nil {
		panic(err)
	}

	table := termtables.CreateTable()
	table.AddHeaders("path", "description")
	cacheManager := data.SourceCacheManager{RootPath: context.RootResource.Dir}
	for _, source := range listUniqSources(context) {
		findApps(source, &cacheManager, table)

	}
	fmt.Println(table.Render())
}

func findApps(source data.Source, manager *data.SourceCacheManager, table *termtables.Table) {

	dir, err := source.GetPath(manager, "")
	if dir == "" {
		return
	}
	if err != nil {
		logrus.Error("Can't find real path of the source")
	}
	err = filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".cache" {
			return filepath.SkipDir
		}
		if path.Base(filePath) == "flekszible.yaml" {
			relpath, err := filepath.Rel(dir, filepath.Dir(filePath))
			if relpath == "." {
				return nil
			}
			if err != nil {
				logrus.Error("Can't find relative path of" + filePath + " " + err.Error())
			}
			fleksz := make(map[string]interface{})
			bytes, err := ioutil.ReadFile(filePath)
			if err != nil {
				logrus.Error("Can't read flekszible.yaml from " + filePath + " " + err.Error())
			}
			name := ""
			err = yaml.Unmarshal(bytes, &fleksz)
			if err != nil {
				logrus.Error("Can't parse flekszible.yaml from " + filePath + " " + err.Error())
			}
			if declaredName, found := fleksz["description"]; found {
				name = declaredName.(string)
				table.AddRow(relpath, name)
			}
		}

		return nil
	})

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
