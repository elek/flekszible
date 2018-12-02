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
