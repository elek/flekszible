package data

import (
	"github.com/pkg/errors"
	"os"
	"path"

	"github.com/elek/flekszible/api/v2/yaml"
)

type FlekszibleResource struct {
	ApiVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Spec       Configuration          `yaml:"spec"`
}

type Configuration struct {
	Name            string                `yaml:",omitempty"`
	Description     string                `yaml:",omitempty"`
	Source          []ConfigSource        `yaml:",omitempty"`
	Import          []ImportConfiguration `yaml:",omitempty"`
	Transformations []yaml.MapSlice       `yaml:",omitempty"`
	ResourcesDir    string                `yaml:"resources,omitempty"`
	Standalone      bool                  `yaml:"-"`
	Header          string                `yaml:"header"`
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

// read flekszible.yaml/Flekszible configuration from one file
func readFromFile(file string, conf *Configuration) (bool, error) {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false, nil
	}
	bytes, err := os.ReadFile(file)
	if err != nil {
		return false, err
	}

	rawMap := map[string]interface{}{}
	err = yaml.Unmarshal(bytes, &rawMap)
	if err != nil {
		return true, errors.Wrapf(err, "%s is not a YAML file", file)
	}

	if rawMap["spec"] != nil && rawMap["kind"] != nil {
		res := FlekszibleResource{}
		err = yaml.UnmarshalStrict(bytes, &res)
		if err != nil {
			return false, errors.Wrapf(err, "%s is couldn't be parsed as flekszible K8s resource file", file)
		}
		// todo: any better way to copy Spec to Conf?
		temp, err := yaml.Marshal(res.Spec)
		if err != nil {
			return false, err
		}
		err = yaml.Unmarshal(temp, &conf)
		if err != nil {
			return false, err
		}
	} else {
		err = yaml.UnmarshalStrict(bytes, &conf)
		if err != nil {
			return false, errors.Wrapf(err, "%s is couldn't be parsed as flekszible YAML file", file)
		}
	}
	return true, nil

}

// ReadConfiguration read configuration from flekszible.yaml or Flekszible file.
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
