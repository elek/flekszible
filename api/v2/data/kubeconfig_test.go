package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadNamespace(t *testing.T) {
	kubeConf := CreateKubeConfigFromFile("../../testdata/config")
	namespace, err := kubeConf.ReadCurrentNamespace()
	assert.Nil(t, err)
	assert.Equal(t, "qwe", namespace)
}
