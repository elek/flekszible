package processor

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	input, err := ioutil.ReadFile("../../testdata/definition/prometheus.yaml")
	assert.Nil(t, err)
	headExpected, err := ioutil.ReadFile("../../testdata/definition/head.yaml")
	assert.Nil(t, err)
	bodyExpected, err := ioutil.ReadFile("../../testdata/definition/body.yaml")
	assert.Nil(t, err)
	head, body := splitDefinitionFile(input)
	assert.Equal(t, string(headExpected), string(head))
	assert.Equal(t, string(bodyExpected), string(body))
}
