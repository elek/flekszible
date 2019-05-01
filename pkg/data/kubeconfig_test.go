package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadNamespace(t *testing.T) {
	kubeConf := CreateKubeConfig("../../testdata/config")
	namespace, err := kubeConf.readCurrentNamespace()
	assert.Nil(t, err)
	assert.Equal(t, "qwe", namespace)
}
