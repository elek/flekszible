package data

import (
	"github.com/elek/flekszible/api/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Keytabs struct {
	Keytabs []Keytab
}

type Keytab struct {
	Name       string
	Principals []string
}

type KeytabGenerator struct {
}

func (*KeytabGenerator) DirName() string {
	return "keytabs"
}

func (kt *KeytabGenerator) Generate(sourceDir string, destinationDir string) ([]*Resource, error) {
	resources := make([]*Resource, 0)
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return resources, err
	}
	for _, file := range files {
		if !file.IsDir() {
			resource, err := kt.GenerateSecret(path.Join(sourceDir, file.Name()), destinationDir)
			if err != nil {
				return resources, errors.Wrap(err, "Can't generate resource from file "+file.Name())
			}
			if resource != nil {
				resources = append(resources, resource)
			}
		}
	}
	return resources, nil
}

func (*KeytabGenerator) GenerateSecret(file string, destinationDir string) (*Resource, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read keytab file: "+file)
	}
	keytabDef := Keytabs{}
	err = yaml.Unmarshal(content, &keytabDef)
	if err != nil {
		return nil, errors.Wrap(err, "Keytab file: "+file+" is not a valid json")
	}

	basename := path.Base(file)
	resourceName := strings.TrimSuffix(basename, filepath.Ext(basename))
	rootNode := NewMapNode(NewPath())
	rootNode.PutValue("apiVersion", "v1")
	rootNode.PutValue("kind", "Secret")
	metadata := rootNode.CreateMap("metadata")
	metadata.PutValue("name", resourceName)
	data := rootNode.CreateMap("data")
	for _, keytab := range keytabDef.Keytabs {

		keytabtabContent, err := encodePrincipals(resourceName+"-"+keytab.Name, keytab.Principals, destinationDir)
		if err != nil {
			return nil, errors.Wrap(err, "Couldn't generate keytab "+keytab.Name+" for resource "+resourceName)
		}
		data.PutValue(keytab.Name, keytabtabContent)
	}
	r := Resource{}
	r.Content = &rootNode
	return &r, nil
}

func encodePrincipals(cacheKey string, principals []string, dir string) (string, error) {
	cacheFile := path.Join(dir, ".state", "keytabs", "principals", cacheKey)
	if _, err := os.Stat(cacheFile); !os.IsNotExist(err) {
		content, err := ioutil.ReadFile(cacheFile)
		if err != nil {
			return "", errors.Wrap(err, "Couldn't read the keytab cache file: "+cacheFile)
		}
		return string(content), nil
	}
	//we need to generate it
	result, err := createKeytab(dir, principals)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't generate the keytab file for "+cacheKey)
	}
	_ = os.MkdirAll(path.Dir(cacheFile), 0700)
	err = ioutil.WriteFile(cacheFile, []byte(result), 0600)
	if err != nil {
		return "", errors.Wrap(err, "Couldn't persist keytab content to the cache file "+cacheFile)
	}
	return result, nil
}

func createKeytab(dir string, principals []string) (string, error) {
	keytabGenScript := path.Join(dir, "getkeytab.sh")
	if _, err := os.Stat(keytabGenScript); os.IsNotExist(err) {
		return "", errors.New("Keytab generator bash script must exist at " + keytabGenScript)
	}
	cmd := exec.Command(keytabGenScript, principals[0])
	cmd.Env = os.Environ()
	output, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "Keytab generation is failed for principals "+strings.Join(principals, " "))
	}
	return string(output), nil
}
