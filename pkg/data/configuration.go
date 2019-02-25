package data

import (
	"github.com/elek/flekszible/pkg/yaml"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Source []ConfigSource
	Import []ImportConfiguration
}

type ConfigSource struct {
	Url  string
	Path string
}

type ImportConfiguration struct {
	Path            string
	Destination     string
	Transformations []yaml.MapSlice
}

//read configuration from flekszible.yaml
func ReadConfiguration(path string) (Configuration, error) {

	conf := Configuration{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return conf, nil
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil

}
