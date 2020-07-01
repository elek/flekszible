package data

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
	"strings"
)

type ConfigGenerator struct {
}

func (*ConfigGenerator) IsManagedDir(dir string) bool {
	return path.Base(dir) == "configmaps"
}

func (*ConfigGenerator) Generate(sourceDir string, destinationDir string) ([]*Resource, error) {
	resources := make([]*Resource, 0)
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return resources, err
	}
	configMapWithData := make(map[string]map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			pieces := strings.SplitN(filename, "_", 2)
			if len(pieces) != 2 {
				return resources, errors.New("Filename should be in the format: configmap_key.ext. " + filename)
			}
			data, err := ioutil.ReadFile(path.Join(sourceDir, filename))
			if err != nil {
				return resources, err
			}
			if _, found := configMapWithData[pieces[0]]; !found {
				configMapWithData[pieces[0]] = make(map[string]string)
			}
			configMapWithData[pieces[0]][pieces[1]] = string(data)
		}
	}
	for name, dataMap := range configMapWithData {
		rootNode := NewMapNode(NewPath())
		rootNode.PutValue("apiVersion", "v1")
		rootNode.PutValue("kind", "ConfigMap")
		metadata := rootNode.CreateMap("metadata")
		metadata.PutValue("name", name)
		data := rootNode.CreateMap("data")
		for keyName, rawData := range dataMap {
			data.PutValue(keyName, rawData)
		}
		r := NewResource()
		r.Content = &rootNode
		//r.Filename = filename
		resources = append(resources, &r)
	}

	return resources, nil
}
