package processor

import (
	"fmt"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path"
	"sort"
	"strings"
	"testing"
)

func TestFromDir(t *testing.T, dir string) data.RenderContext {
	outputDir := path.Join("../../build", dir)
	inputDir := path.Join("../../testdata", dir)
	expectedDir := path.Join(inputDir, "expected")
	return TestDirAndCompare(t, inputDir, outputDir, expectedDir)
}

func TestExample(t *testing.T, name string) data.RenderContext {
	outputDir := path.Join("../../build/examples", name)
	inputDir := path.Join("../../examples", name)
	expectedDir := path.Join("../../testdata/examplesresults", name)
	return TestDirAndCompare(t, inputDir, outputDir, expectedDir)
}

func TestDirAndCompare(t *testing.T, inputDir string, outputDir string, expectedDir string) data.RenderContext {

	context := data.RenderContext{
		OutputDir: outputDir,
		Mode:      "k8s",
		InputDir:  []string{inputDir},
	}
	context.ReadConfigs()
	LoadDefinitions(&context)

	repository := CreateProcessorRepository()

	for _, directory := range context.InputDir {
		repository.ParseProcessors(directory)
	}
	repository.Append(&K8sWriter{})
	Generate(repository, &context)
	compareDir(t, expectedDir, outputDir)
	return context
}

func compareDir(t *testing.T, expected string, result string) {
	exp := readDir(t, expected)
	res := readDir(t, result)

	assert.Equal(t, keysFromMap(exp), keysFromMap(res))
	for key, _ := range res {
		assert.Equal(t, exp[key], res[key], "File "+key+" is different")
	}
}

func keysFromMap(nodes map[string]*data.MapNode) []string {
	keys := make([]string, 0)
	for key, _ := range nodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func readDir(t *testing.T, dirName string) map[string]*data.MapNode {
	dirContent := make(map[string]*data.MapNode)
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filename := path.Join(dirName, f.Name())
		fileContent, err := ioutil.ReadFile(filename)
		assert.Nil(t, err)
		parsedFragment, err := data.ReadString(fileContent)
		assert.Nil(t, err)
		dirContent[f.Name()] = parsedFragment
	}
	return dirContent
}

func ExecuteProcessorAndCompare(t *testing.T, dir string, prefix string) data.RenderContext {
	testdir := path.Join("../../testdata", dir)
	//given\
	processors, err := ReadProcessorDefinitionFile(path.Join(testdir, prefix+"_config.yaml"))
	assert.Nil(t, err)

	processor := processors[0]
	resources, err := data.LoadFrom(testdir, prefix+".yaml")
	assert.Nil(t, err)

	ctx := data.RenderContext{
		Resources: resources,
	}
	//when
	processor.Before(&ctx)
	processor.BeforeResource(&ctx.Resources[0])
	ctx.Resources[0].Content.Accept(processor)

	//then
	expected, err := data.LoadFrom(testdir, prefix+"_expected.yaml")
	assert.Nil(t, err)

	assert.EqualValues(t, ToSimpleYaml(&expected[0]), ToSimpleYaml(&ctx.Resources[0]))

	return ctx

}

func ToSimpleYaml(resource *data.Resource) interface{} {
	buffer := strings.Builder{}
	writer := K8sWriter{
		output: &buffer,
	}
	resource.Content.Accept(&writer)
	var res interface{}
	err := yaml.Unmarshal([]byte(buffer.String()), &res)
	fmt.Println(string(buffer.String()))

	if err != nil {
		panic(err)

	}
	return res
}
