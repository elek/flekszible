package processor

import (
	"fmt"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
)

func TestFromDir(t *testing.T, dir string) *RenderContext {
	outputDir := path.Join("../../target", dir)
	inputDir := path.Join("../../testdata", dir)
	expectedDir := path.Join(inputDir, "expected")
	return TestDirAndCompare(t, inputDir, outputDir, expectedDir)
}

func TestExample(t *testing.T, name string) *RenderContext {
	outputDir := path.Join("../../target/examples", name)
	inputDir := path.Join("../../examples", name)
	expectedDir := path.Join("../../testdata/examplesresults", name)
	return TestDirAndCompare(t, inputDir, outputDir, expectedDir)
}

func TestDirAndCompare(t *testing.T, inputDir string, outputDir string, expectedDir string) *RenderContext {

	context := CreateRenderContext("k8s", inputDir, outputDir)

	err := context.Init()
	if err != nil {
		println(err.Error())
	}
	assert.Nil(t, err)

	context.AppendProcessor(&K8sWriter{})
	context.Render()
	compareDir(t, expectedDir, outputDir)
	return context
}

func compareDir(t *testing.T, expected string, result string) {
	exp := readDir(t, expected)
	res := readDir(t, result)

	assert.Equal(t, keysFromMap(exp), keysFromMap(res))
	for key := range res {
		assert.Equal(t, exp[key], res[key], "File "+key+" is different")
	}
}

func keysFromMap(nodes map[string]*data.MapNode) []string {
	keys := make([]string, 0)
	for key := range nodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func readDir(t *testing.T, dirName string) map[string]*data.MapNode {
	dirContent := make(map[string]*data.MapNode)
	readOneDir(t, dirName, &dirContent)
	return dirContent
}

func readOneDir(t *testing.T, dirName string, dirContent *map[string]*data.MapNode) {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filename := path.Join(dirName, f.Name())
		if info, err := os.Stat(filename); err == nil && !info.IsDir() {
			fileContent, err := ioutil.ReadFile(filename)
			assert.Nil(t, err)
			parsedFragment, err := data.ReadManifestString(fileContent)
			assert.Nil(t, err)
			(*dirContent)[f.Name()] = parsedFragment
		} else {
			readOneDir(t, filename, dirContent)
		}
	}
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
