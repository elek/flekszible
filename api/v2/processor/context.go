package processor

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RenderContext struct {
	OutputDir     string
	Conf          data.Configuration
	Mode          string
	ImageOverride string
	Namespace     string
	Registry      *ProcessorTypes
	RootResource  *ResourceNode
}

type ResourceNode struct {
	Name   string      //the original import name used in FLekszible files
	Dir    string      //absolut path of the current dir
	Origin data.Source //source used to load the resource

	Destination              string
	Resources                []*data.Resource
	Children                 []*ResourceNode
	PreImportTransformations []byte
	Source                   []data.Source
	ProcessorRepository      *TransformationRepository
	ResourcesDir             string
	Definitions              []string //definitions loaded from this dir
}

func (context *ResourceNode) FindChildResource(path string) (*ResourceNode, error) {
	for _, node := range context.Children {
		if node.Name == path {
			return node, nil
		}
	}
	return nil, errors.New("No such node with import path " + path)
}

type ResourceLocation struct {
	Name   string      //the name to import the location (like 'ozone')
	Source data.Source //the source used to load the resource
	Dir    string      //absolute path of the location
}

func CreateRenderContextFromDir(inputDir string) *RenderContext {
	return CreateRenderContext("k8s", inputDir, "-")
}

func CreateRenderContext(mode string, inputDir string, outputDir string) *RenderContext {
	res := &RenderContext{
		OutputDir:    outputDir,
		Mode:         mode,
		Registry:     NewRegistry(),
		RootResource: CreateResourceNode("<PROJECT_DIR>", inputDir, "", &data.LocalSource{Dir: inputDir}),
	}
	ActivateRun(res.Registry)
	ActivateK8sWriter(res.Registry)
	ActivateEnv(res.Registry)
	ActivateReplace(res.Registry)
	ActivateSubstitute(res.Registry)
	ActivateInit(res.Registry)
	ActivateNamespace(res.Registry)
	ActivatePrefix(res.Registry)
	ActivatePipe(res.Registry)
	Activate(res.Registry)
	ActivateImageSet(res.Registry)
	ActivatePublishStatefulset(res.Registry)
	ActivateDaemonToStateful(res.Registry)
	ActivateMount(res.Registry)
	ActivateCleanup(res.Registry)
	ActivateConfigHash(res.Registry)
	ActivateRemove(res.Registry)
	ActivatePublishService(res.Registry)
	ActivateNameFilter(res.Registry)
	ActivateAdd(res.Registry)
	ActivateKustomize(res.Registry)
	ActivateMerge(res.Registry)
	ActivatePatch(res.Registry)
	return res

}

func (context *RenderContext) ListResourceNodes() []*ResourceNode {
	return listResourceNodesInt(context.RootResource)
}

func (context *RenderContext) AddAdHocTransformations(transformations []string) error {

	proc, err := context.Registry.createTransformation(transformations)
	if err != nil {
		return err
	}
	context.RootResource.ProcessorRepository.AppendAll(proc)
	return nil
}

func (registry *ProcessorTypes) parseTransformation(trafoDef string) (Processor, error) {
	parts := strings.SplitN(trafoDef, ":", 2)
	name := parts[0]
	parameterMap := make(map[string]string)
	if len(parts) > 1 {
		transofmationsString := strings.ReplaceAll(parts[1], "\\,", "__NON_SEPARATOR_COMA__")
		for _, rawParam := range strings.Split(transofmationsString, ",") {
			parameter := strings.ReplaceAll(rawParam, "__NON_SEPARATOR_COMA__", ",")
			paramParts := strings.SplitN(parameter, "=", 2)
			if len(paramParts) < 2 {
				return nil, errors.New("Parameters should be defined in the form key=value and not like " + parameter)
			}
			parameterMap[paramParts[0]] = paramParts[1]
		}
	}
	proc, err := registry.Create(name, parameterMap)
	if err != nil {
		return nil, errors.Wrap(err, "Can't create transformation based on the string "+trafoDef)
	}
	return proc, nil
}

// parse one-liner transformation definition
func (registry *ProcessorTypes) createTransformation(transformationsDefinitions []string) ([]Processor, error) {
	result := make([]Processor, 0)
	for _, trafoDef := range strings.Split(os.Getenv("FLEKSZIBLE_TRANSFORMATION"), ";") {
		if len(strings.TrimSpace(trafoDef)) > 0 {
			transformation, err := registry.parseTransformation(trafoDef)
			if err != nil {
				return result, errors.Wrap(err, "Can't parse transformation defined by FLEKSZIBLE_TRANSFORMATION: "+trafoDef)
			}
			result = append(result, transformation)
		}
	}
	for _, transformationsDefinition := range transformationsDefinitions {
		transformation, err := registry.parseTransformation(transformationsDefinition)
		if err != nil {
			return result, errors.Wrap(err, "Can't parse transformation defined by cli arg: "+transformationsDefinition)
		}
		result = append(result, transformation)
	}
	return result, nil
}

func listResourceNodesInt(node *ResourceNode) []*ResourceNode {
	result := make([]*ResourceNode, 0)
	result = append(result, node)
	for _, child := range node.Children {
		result = append(result, listResourceNodesInt(child)...)
	}
	return result
}

func (context *RenderContext) LoadResourceTree() error {
	data.Generators = append(data.Generators, &data.OutputGenerator{})
	data.Generators = append(data.Generators, &data.ConfigImporter{})
	data.Generators = append(data.Generators, &data.SecretGenerator{})
	data.Generators = append(data.Generators, &data.SecretImporter{})
	cacheManager := data.NewSourceCacheManager(context.RootResource.Dir)
	return context.RootResource.InitFromDir(&cacheManager, context.OutputDir)
}

// List all the resources from the resource tree.
func (context *RenderContext) Resources() []*data.Resource {
	result := make([]*data.Resource, 0)
	result = append(result, context.RootResource.AllResources()...)
	return result
}

func (node *ResourceNode) AllResources() []*data.Resource {
	result := make([]*data.Resource, 0)
	for _, child := range node.Children {
		result = append(result, child.AllResources()...)
	}
	result = append(result, node.Resources...)
	return result
}

func CreateResourceNode(name string, dir string, destination string, source data.Source) *ResourceNode {
	node := ResourceNode{
		Name:                name,
		Dir:                 dir,
		ProcessorRepository: CreateProcessorRepository(),
		Children:            make([]*ResourceNode, 0),
		Destination:         destination,
		Origin:              source,
	}
	return &node
}
func (context *RenderContext) AddResources(resources ...*data.Resource) {
	context.RootResource.Resources = append(context.RootResource.Resources, resources...)
}

func (context *RenderContext) Init() error {
	//load flekszible.yaml (recursive)
	err := context.LoadResourceTree()
	if err != nil {
		return err
	}
	//load 'definitions' subdirs
	err = context.LoadDefinitions()
	if err != nil {
		return err
	}
	//load 'transformations' subdirs
	return context.InitializeTransformations()
}

func (context *RenderContext) InitializeTransformations() error {
	return context.RootResource.InitializeTransformations(context)
}

func (node *ResourceNode) InitializeTransformations(context *RenderContext) error {
	if node.PreImportTransformations != nil {
		processors, err := context.Registry.ReadProcessorDefinition(node.PreImportTransformations)
		if err != nil {
			return errors.Wrapf(err, "error in %s", node.Dir)
		}
		node.ProcessorRepository.AppendAll(processors)
	}

	processors, e := context.Registry.ParseTransformations(node.Dir)
	if e != nil {
		return errors.Wrap(e, "Can't read transformations from "+node.Dir)
	} else {
		for _, processor := range processors {
			if processor.GetScope() == "global" {
				context.RootResource.ProcessorRepository.Append(processor)
			} else {
				node.ProcessorRepository.InsertToBeginning(processor)
			}
		}
	}

	for _, child := range node.Children {
		err := child.InitializeTransformations(context)
		if err != nil {
			return err
		}
	}
	return nil
}

func (context *RenderContext) AppendCustomProcessor(name string, parameters map[string]string) error {
	if definition, ok := context.Registry.TypeMap[strings.ToLower(name)]; ok {
		config := yaml.MapSlice{}
		for k, v := range parameters {
			config.Put(k, v)
		}
		processor, err := definition.Factory(&config)
		if err != nil {
			return err
		}
		context.RootResource.ProcessorRepository.Append(processor)
		return nil
	} else {
		return errors.New("No such processor definition")
	}
}

func (context *RenderContext) AppendProcessor(processor Processor) {
	root := context.RootResource
	repo := root.ProcessorRepository
	repo.Append(processor)
}

type execute func(transformation Processor, context *RenderContext, node *ResourceNode) error

func (node *ResourceNode) Execute(context *RenderContext, functionCall execute) error {
	//execute on child
	for _, child := range node.Children {
		err := child.Execute(context, functionCall)
		if err != nil {
			return errors.Wrap(err, "The function can't be executed on one of the childResources")
		}
	}

	//execute on this resource
	for _, transformation := range node.ProcessorRepository.Processors {
		err := functionCall(transformation, context, node)
		if err != nil {
			return errors.Wrap(err, "The function can't be executed on all the resources. Transformation: ")
		}
	}
	return nil

}

// execute all the transformations in the rendercontext.
func (ctx *RenderContext) Render() error {

	before := func(processor Processor, context *RenderContext, node *ResourceNode) error {
		err := processor.Before(context, node)
		if err != nil {
			return errors.Wrap(err, "Before execution phase is failed")
		}
		return nil
	}
	after := func(processor Processor, context *RenderContext, node *ResourceNode) error {
		err := processor.After(context, node)
		if err != nil {
			return errors.Wrap(err, "After execution phase is failed")
		}
		return nil
	}
	process := func(processor Processor, context *RenderContext, node *ResourceNode) error {
		for _, resource := range node.AllResources() {
			err := processor.BeforeResource(resource)
			if err != nil {
				return errors.Wrap(err, "Applying transformation BeforeResource "+resource.Filename+" is failed")
			}
			resource.Content.Accept(processor)
			err = processor.AfterResource(resource)
			if err != nil {
				return errors.Wrap(err, "Applying transformation AfterResource "+resource.Filename+" is failed")
			}

		}
		return nil
	}
	err := ctx.RootResource.Execute(ctx, func(processor Processor, context *RenderContext, node *ResourceNode) error {
		err := processor.RegisterResources(context, node)
		if err != nil {
			return errors.Wrap(err, "RegisterResource execution phase is failed")
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = ctx.RootResource.Execute(ctx, before)
	if err != nil {
		return err
	}
	err = ctx.RootResource.Execute(ctx, process)
	if err != nil {
		return err
	}
	err = ctx.RootResource.Execute(ctx, after)
	if err != nil {
		return err
	}
	return nil
}

// InitFromDir parses the directory structure and the flekszible configs from the dirs.
func (node *ResourceNode) InitFromDir(sourceCache *data.SourceCacheManager, outputDir string) error {
	conf, _, err := data.ReadConfiguration(node.Dir)
	if err != nil {
		return errors.Wrap(err, "Can't parse flekszible.yaml/Flekszible descriptor from  "+node.Dir)
	}
	return node.InitFromConfig(conf, sourceCache, outputDir)
}

// InitFromConfig initializes and reads all dependencies based on parsed config.
func (node *ResourceNode) InitFromConfig(conf data.Configuration, sourceCache *data.SourceCacheManager, outputDir string) error {
	//output dir should never be read
	absNodeDir, _ := filepath.Abs(node.Dir)
	absDestDir, _ := filepath.Abs(outputDir)

	node.Resources = data.ReadResourcesFromDir(path.Join(node.Dir, "resources"))

	if absNodeDir != absDestDir {
		node.Resources = append(node.Resources, data.ReadResourcesFromDir(node.Dir)...)
	}

	for _, generator := range data.Generators {
		dirs, err := os.ReadDir(node.Dir)
		if err == nil {
			for _, dir := range dirs {
				managedDir := path.Join(node.Dir, dir.Name())
				if dir.IsDir() && generator.IsManagedDir(managedDir) {
					resources, err := generator.Generate(managedDir, outputDir)
					if err != nil {
						return errors.Wrap(err, "Can't generate resources from the dir "+managedDir)
					}
					node.Resources = append(node.Resources, resources...)

				}
			}
		}
	}

	node.Source = make([]data.Source, 0)
	for _, definedSource := range conf.Source {
		if definedSource.Url != "" {
			node.Source = append(node.Source, &data.RemoteSource{Url: definedSource.Url})
		} else if definedSource.Path != "" {
			var sourceDir string
			if path.IsAbs(definedSource.Path) {
				sourceDir = definedSource.Path
			} else {
				sourceDir = path.Join(node.Dir, definedSource.Path)
			}
			node.Source = append(node.Source, &data.LocalSource{Dir: sourceDir})
		}
	}
	//update destinations of the direct k8s resources
	for ix := range node.Resources {
		node.Resources[ix].Destination = node.Destination
	}

	//read inline transformations
	if len(conf.Transformations) > 0 {
		bytes, err := yaml.Marshal(conf.Transformations)
		if err != nil {
			return err
		}
		node.PreImportTransformations = bytes
	}

	//import _global directories
	for _, source := range node.Source {
		sourceDir, err := source.GetPath(sourceCache)
		if err != nil {
			return errors.Wrap(err, "Source definition defined in "+node.Dir+" couldn't be loaded")
		}
		globalDirs := []string{path.Join(sourceDir, "flekszible", "_global"), path.Join(sourceDir, "_global")}
		for _, globalDir := range globalDirs {
			if stat, err := os.Stat(globalDir); err == nil && stat.IsDir() {
				childNode := CreateResourceNode("_global", globalDir, node.Destination, source)
				childNode.Origin = source
				err = childNode.InitFromDir(sourceCache, outputDir)
				if err != nil {
					return errors.Wrap(err, "Couldn't load _global dir from "+globalDir)
				}
				node.Children = append(node.Children, childNode)
				break
			}
		}
	}

	destinations := make([]string, 0)
	//ignore, if it's a destination directory
	for _, oneImport := range conf.Import {
		absImportDest, err := filepath.Abs(path.Join(node.Dir, oneImport.Destination))
		if err == nil {
			destinations = append(destinations, absImportDest)
		}
	}

	//read imported directories
	for _, importDefinition := range conf.Import {
		source, err := LocateImportedDir(node.Dir, importDefinition.Path, node.Source, sourceCache, destinations)
		if err != nil {
			return err
		}
		if source == nil {
			return errors.New("Directory dependency `" + importDefinition.Path + "` defined in " + node.Dir + " can't be found")
		}
		sourceDir, _ := source.GetPath(sourceCache)
		resourceDir := checkPath(sourceDir, importDefinition.Path)

		childNode := CreateResourceNode(importDefinition.Path, resourceDir, importDefinition.Destination, source)
		if importDefinition.Destination == "" {
			childNode.Destination = node.Destination
		}
		childNode.Origin = source
		if node.Dir == childNode.Dir {
			panic("Recursive directory parser " + node.Dir + " loads" + childNode.Dir)
		}
		err = childNode.InitFromDir(sourceCache, outputDir)
		if err != nil {
			return err
		}
		if len(importDefinition.Transformations) > 0 {
			bytes, err := yaml.Marshal(importDefinition.Transformations)
			if err != nil {
				return err
			}
			childNode.PreImportTransformations = bytes
		}
		node.Children = append(node.Children, childNode)
	}
	return nil
}

// LoadDefinitions load. transformation definitions from ./definitions dir (all dir).
func (ctx *RenderContext) LoadDefinitions() error {
	return ctx.RootResource.LoadDefinitions(ctx.Registry)
}

func (node *ResourceNode) LoadDefinitions(registry *ProcessorTypes) error {
	defDir := path.Join(node.Dir, "definitions")
	if _, err := os.Stat(defDir); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(defDir)
		if err != nil {
			logrus.Warningf("Can't load definition directory %s: %s", defDir, err.Error())
		}
		for _, file := range files {
			definitionFile := path.Join(defDir, file.Name())
			name, err := registry.parseDefintion(definitionFile)
			if err != nil {
				return errors.Wrap(err, "Can't parse the definition file "+definitionFile)
			}
			if node.Definitions == nil {
				node.Definitions = make([]string, 0)
			}
			node.Definitions = append(node.Definitions, name)
		}
	}

	for _, child := range node.Children {
		err := child.LoadDefinitions(registry)
		if err != nil {
			return err
		}
	}
	return nil

}

// try to find the first Source which contains the dir
func LocateImportedDir(basedir string, dir string, sources []data.Source, cacheManager *data.SourceCacheManager, excludedDirs []string) (data.Source, error) {
	allSources := make([]data.Source, 0)
	allSources = append(allSources, data.LocalSourcesFromEnv()...)
	allSources = append(allSources, &data.LocalSource{Dir: basedir})
	allSources = append(allSources, sources...)

outside:
	for _, source := range allSources {
		resourcePath, err := source.GetPath(cacheManager)
		if err != nil {
			value := source.ToString()
			logrus.Error("Can't check dir from the source " + value + err.Error())
		} else if resourcePath != "" {
			path := checkPath(resourcePath, dir)

			if path != "" {
				// is it excluded?
				for _, exclude := range excludedDirs {
					if path == exclude {
						continue outside
					}
				}

				return source, nil
			}
		}
	}
	return nil, errors.New("Couldn't find dir: " + dir)
}

func checkPath(baseDir string, subdir string) string {
	optionalSubDir := path.Join(baseDir, "flekszible")
	if stat, err := os.Stat(optionalSubDir); !os.IsNotExist(err) && stat.IsDir() {
		result := path.Join(baseDir, "flekszible", subdir)
		if _, err := os.Stat(result); !os.IsNotExist(err) {
			return result
		}
	}
	result := path.Join(baseDir, subdir)
	if _, err := os.Stat(result); !os.IsNotExist(err) {
		return result
	}
	return ""
}
