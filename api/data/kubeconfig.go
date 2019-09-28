package data

import (
	"github.com/elek/flekszible/api/yaml"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type Kubeconfig struct {
	Path string
}

func findKubeConfigFile() string {
	file := os.Getenv("KUBECONFIG")
	if _, err := os.Stat(file); err == nil {
		return file
	}
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	file = path.Join(usr.HomeDir, ".kube", "config")
	if _, err := os.Stat(file); err == nil {
		return file
	}
	return ""
}

func CreateKubeConfig() Kubeconfig {
	kubeConfigFile := findKubeConfigFile()
	if kubeConfigFile != "" {
		return CreateKubeConfigFromFile(kubeConfigFile)
	} else {
		return Kubeconfig{
			Path: "",
		}
	}
}

func CreateKubeConfigFromFile(path string) Kubeconfig {
	return Kubeconfig{
		Path: path,
	}
}

func (config *Kubeconfig) ReadCurrentNamespace() (string, error) {
	if config.Path == "" {
		return "default", nil
	}
	content, err := ioutil.ReadFile(config.Path)
	if err != nil {
		return "", err
	}
	parsedConfig := make(map[string]interface{})
	err = yaml.Unmarshal(content, &parsedConfig)
	if err != nil {
		return "", err
	}
	if current, ok := parsedConfig["current-context"]; ok {
		if contexts, ok := parsedConfig["contexts"]; ok {
			for _, context := range contexts.([]interface{}) {
				if name, ok := (context.(yaml.MapSlice)).Get("name"); ok {
					if name != current {
						continue
					}
				} else {
					continue
				}
				if contexMetadata, ok := context.(yaml.MapSlice).Get("context"); ok {
					if namespace, ok := contexMetadata.(yaml.MapSlice).Get("namespace"); ok {
						return namespace.(string), nil
					}
				}
			}
		}
	}
	return "default", nil

}
