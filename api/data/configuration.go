package data

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/elek/flekszible/api/yaml"
)

type Configuration struct {
	Source          []ConfigSource        `yaml:",omitempty"`
	Import          []ImportConfiguration `yaml:",omitempty"`
	Transformations []yaml.MapSlice       `yaml:",omitempty"`
	ResourcesDir    string                `yaml:"resources,omitempty"`
	Standalone      bool                  `yaml:"-"`
}

type ConfigSource struct {
	Url  string `yaml:",omitempty"`
	Path string `yaml:",omitempty"`
}

type ImportConfiguration struct {
	Path            string
	Destination     string          `yaml:",omitempty"`
	Transformations []yaml.MapSlice `yaml:",omitempty"`
}

//read flekszible.yaml/Flekszible configuration from one file
func readFromFile(file string, conf *Configuration) (bool, error) {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return false, err
		}

		err = yaml.Unmarshal(bytes, &conf)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		return false, nil
	}

}

//read configuration from flekszible.yaml or Flekszible file
func ReadConfiguration(dir string) (Configuration, string, error) {

	conf := Configuration{}
	conf.ResourcesDir = "resources"
	configFilePath := path.Join(dir, "flekszible.yaml")
	loaded, err := readFromFile(configFilePath, &conf)
	if err != nil {
		return conf, "", err
	}
	if loaded {
		conf.ResourcesDir = ""
		return conf, configFilePath, nil
	}

	configFilePath = path.Join(dir, "Flekszible")
	loaded, err = readFromFile(configFilePath, &conf)
	if err != nil {
		return conf, "", err
	}
	if loaded {
		return conf, configFilePath, nil
	}
	conf.ResourcesDir = ""
	return conf, "", nil

}
