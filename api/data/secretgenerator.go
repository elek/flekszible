package data

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type SecretGenerator struct {
}

func (*SecretGenerator) IsManagedDir(dir string) bool {
	descriptorFile := path.Join(dir, ".secretgen")
	if _, err := os.Stat(descriptorFile); !os.IsNotExist(err) {
		return true
	}
	return false
}

func (kt *SecretGenerator) Generate(sourceDir string, destinationDir string) ([]*Resource, error) {
	resources := make([]*Resource, 0)
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return resources, err
	}
	for _, file := range files {
		if !file.IsDir() {
			if !file.IsDir() && file.Name() != ".secretgen" {
				descriptor := path.Join(sourceDir, file.Name())
				ext := path.Ext(file.Name())
				name := file.Name()[0 : len(file.Name())-len(ext)]
				descriptors, err := GetEncodedSecrets("getkeystore.sh", destinationDir, name, descriptor)
				if err != nil {
					return resources, errors.Wrap(err, "Can't generate encoded secret for "+descriptor)
				}
				resource, err := GenerateSecret(destinationDir, name, descriptors)
				if err != nil {
					return resources, errors.Wrap(err, "Can't generate secret from file "+descriptor)
				}
				resources = append(resources, resource)

			}
		}
	}
	return resources, nil
}

func GetEncodedSecrets(scriptName string, destinationDir string, name string, descriptorPath string) (map[string]string, error) {
	result := make(map[string]string)
	script := path.Join(destinationDir, scriptName)
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return result, errors.New("Secret generator bash script must exist at " + script)
	}
	cmd := exec.Command(script, name, descriptorPath)
	cmd.Env = os.Environ()
	cmd.Dir = destinationDir
	output, err := cmd.Output()
	if err != nil {
		return result, errors.Wrap(err, "Secret generation is failed for  "+path.Join(destinationDir, scriptName)+" "+name+" "+descriptorPath)
	}
	for _, line := range strings.Split(string(output), "\n") {
		kv := strings.SplitN(line, " ", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result, nil
}

func GenerateSecret(destinationDir string, name string, secrets map[string]string) (*Resource, error) {
	rootNode := NewMapNode(NewPath())
	rootNode.PutValue("apiVersion", "v1")
	rootNode.PutValue("kind", "Secret")
	metadata := rootNode.CreateMap("metadata")
	metadata.PutValue("name", name)
	data := rootNode.CreateMap("data")
	for secretName, secretValue := range secrets {
		data.PutValue(secretName, secretValue)
	}
	r := Resource{}
	r.Content = &rootNode
	return &r, nil
}
