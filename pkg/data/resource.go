package data

import (
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
	Path     string
	Filename string
	Content  *MapNode
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
func LoadFromFileInfo(dir string, file os.FileInfo) ([]Resource, error) {
	return LoadFrom(dir, file.Name())
}

//Load k8s resources from one yaml file
func LoadFrom(dir string, file string) ([]Resource, error) {
	results := make([]Resource, 0)
	fullPath := path.Join(dir, file)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return results, err
	}

	fragments := regexp.MustCompile("---\n").Split(string(content), -1)

	for _, fragment := range fragments {
		fragmentContent := strings.TrimSpace(fragment)
		if len(fragmentContent) > 0 {

			parsedFragment, err := ReadString([]byte(fragmentContent))
			if err == nil {
				r := Resource{}
				r.Content = parsedFragment
				r.Filename = file
				results = append(results, r)
			} else {
				return results, fmt.Errorf("Can't parse the resource file %s: %s", fullPath, err)
			}
		}
	}
	return results, nil
}

func ReadResourcesFromDir(dir string) []Resource {
	logrus.Infof("Reading resources from %s", dir)
	resources := make([]Resource, 0)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() && file.Name() != "flekszible.yaml" && filepath.Ext(file.Name()) == ".yaml" {
			resource, err := LoadFromFileInfo(dir, file)
			if err != nil {
				panic(err)
			}
			resources = append(resources, resource...)
		}
	}

	return resources
}
