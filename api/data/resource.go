package data

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Resource struct {
	Path                string
	Filename            string
	Content             *MapNode
	Destination         string
	DestinationFileName string
}

func (r *Resource) Name() string {
	cleanUpPattern := regexp.MustCompile(`{.*}-?`)
	rawName := r.Get(NewPath("metadata", "name"))
	return cleanUpPattern.ReplaceAllString(rawName, "")
}

func (r *Resource) Kind() string {
	return r.Get(NewPath("kind"))
}

func (r *Resource) Get(path Path) string {
	get := Get{Path: path}
	r.Content.Accept(&get)
	if get.Found {
		return get.ReturnValue.(*KeyNode).Value.(string)
	} else {
		return ""
	}
}
func LoadFromFileInfo(dir string, file os.FileInfo) ([]*Resource, error) {
	return LoadFrom(dir, file.Name())
}

//Load k8s resources from one yaml file
func LoadFrom(dir string, file string) ([]*Resource, error) {
	results := make([]*Resource, 0)
	fullPath := path.Join(dir, file)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return results, err
	}

	fragments := regexp.MustCompile("---\n").Split(string(content), -1)

	for _, fragment := range fragments {
		fragmentContent := strings.TrimSpace(fragment)
		if len(fragmentContent) > 0 {

			parsedFragment, err := ReadManifestString([]byte(fragmentContent))
			if err == nil {
				r := Resource{}
				r.Content = parsedFragment
				r.Filename = file
				results = append(results, &r)
			} else {
				return results, fmt.Errorf("Can't parse the resource file %s: %s", fullPath, err)
			}
		}
	}
	return results, nil
}

//read all the resources from a directory
func ReadResourcesFromDir(dir string) []*Resource {
	logrus.Infof("Reading resources from %s", dir)
	resources := make([]*Resource, 0)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return resources
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() && file.Name() != "flekszible.yaml" && (
			filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml") {
			resource, err := LoadFromFileInfo(dir, file)
			if err != nil {
				panic(err)
			}
			resources = append(resources, resource...)
		}
	}

	return resources
}

func ReadConfigMaps(dir string) ([]*Resource, error) {
	resources := make([]*Resource, 0)
	configDirPath := path.Join(dir, "configmaps")
	if _, err := os.Stat(configDirPath); err != nil {
		if os.IsNotExist(err) {
			return resources, nil
		} else {
			return resources, err
		}
	}
	files, err := ioutil.ReadDir(configDirPath)
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
			data, err := ioutil.ReadFile(path.Join(configDirPath, filename))
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
		r := Resource{}
		r.Content = &rootNode
		//r.Filename = filename
		resources = append(resources, &r)
	}

	return resources, nil
}
