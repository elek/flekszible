package data

import (
	"github.com/elek/flekszible/pkg/yaml"
	"io/ioutil"
	"os"
	"path"
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

//read flekszible.yaml configuration from one file
func readFromFile(file string, conf *Configuration) error {
	if _, err := os.Stat(file); ! os.IsNotExist(err) {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(bytes, &conf)
		if err != nil {
			return err
		}
	}
	return nil

}

//read configuration from flekszible.yaml or Flekszible file
func ReadConfiguration(dir string) (Configuration, error) {

	conf := Configuration{}

	err := readFromFile(path.Join(dir, "flekszible.yaml"), &conf)
	if err != nil {
		return conf, err
	}

	err = readFromFile(path.Join(dir, "Flekszible"), &conf)
	if err != nil {
		return conf, err
	}
	return conf, nil

}
