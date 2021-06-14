package data

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type SecretGetConfig struct {
	Script string
}

type SecretGenerator struct {
}

func ReadConfig(file string) (SecretGetConfig, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		errors.Wrap(err, "Couldn't load secret generator file "+file)
	}
	config := SecretGetConfig{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		errors.Wrap(err, "Couldn't parse secret generator file "+file)
	}
	return config, nil
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
	config, err := ReadConfig(path.Join(sourceDir, ".secretgen"))
	if err != nil {
		return resources, err
	}
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
				descriptors, err := GetEncodedSecrets(config.Script, destinationDir, name, descriptor)
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
	logrus.Info("Generating secret " + name + " with executing '" + script + " " + descriptorPath + "'")
	cmd := exec.Command(script, descriptorPath)
	var stdout, stderr bytes.Buffer
	cmd.Env = os.Environ()
	cmd.Dir = destinationDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		command := path.Join(destinationDir, scriptName) + " " + descriptorPath
		logrus.Error("Command is failed " + command)
		logrus.Error("Stdout: " + string(stdout.Bytes()))
		logrus.Error("Stderr: " + string(stderr.Bytes()))
		return result, errors.Wrap(err, "Secret generation is failed for  cd '"+destinationDir+" && "+command+"'")
	}
	for _, line := range strings.Split(string(string(stdout.Bytes())), "\n") {
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
