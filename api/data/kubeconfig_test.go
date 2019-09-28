package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadNamespace(t *testing.T) {
	kubeConf := CreateKubeConfigFromFile("../../testdata/config")
	namespace, err := kubeConf.ReadCurrentNamespace()
	assert.Nil(t, err)
	assert.Equal(t, "qwe", namespace)
}
