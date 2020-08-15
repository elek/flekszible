package data

import (
	"github.com/pkg/errors"
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
	Metadata            map[string]string
}

func NewResource() Resource {
	return Resource{
		Metadata: make(map[string]string),
	}
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
func LoadResourceFromFileInfo(dir string, file os.FileInfo) ([]*Resource, error) {
	return LoadResourceFromFile(dir, file.Name())
}

func LoadResourceFromByte(data []byte) ([]*Resource, error) {
	results := make([]*Resource, 0)
	fragments := regexp.MustCompile("---\n").Split(string(data), -1)

	for _, fragment := range fragments {
		fragmentContent := strings.TrimSpace(fragment)
		if len(fragmentContent) > 0 {

			parsedFragment, err := ReadManifestString([]byte(fragmentContent))
			if err == nil {
				r := NewResource()
				r.Content = parsedFragment
				results = append(results, &r)
			} else {
				return results, errors.Wrap(err, "Can't parse the resource.")
			}
		}
	}
	return results, nil
}

//Load k8s resources from one yaml file
func LoadResourceFromFile(dir string, file string) ([]*Resource, error) {

	fullPath := path.Join(dir, file)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	resources, err := LoadResourceFromByte(content)
	if err != nil {
		return nil, errors.Wrap(err, "Can't load resource from file "+fullPath)
	}
	for _, resource := range resources {
		resource.Filename = file
	}
	return resources, nil
}

//read all the resources from a directory
func ReadResourcesFromDir(dir string) []*Resource {
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
			resource, err := LoadResourceFromFileInfo(dir, file)
			if err != nil {
				panic(err)
			}
			resources = append(resources, resource...)
		}
	}

	return resources
}

