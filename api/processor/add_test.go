package processor

import (
	"testing"
)

func TestAddArrayMap(t *testing.T) {
	TestFromDir(t, "add/arraymap")
}

func TestAddMapElement(t *testing.T) {
	TestFromDir(t, "add/mapelement")
}

func TestAddMapElementFiltered(t *testing.T) {
	TestFromDir(t, "add/filtered")
}

func TestAddMapVolume(t *testing.T) {
	TestFromDir(t, "add/volume")
}

func TestYamlize(t *testing.T) {
	TestFromDir(t, "add/yamlize")
}
