package operator

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonProcess(t *testing.T) {
	content, err := ioutil.ReadFile("../../testdata/operator/example.json")
	assert.Nil(t, err)
	response, err := handleRequest("../../testdata/operator", content)
	assert.Nil(t, err)
	assert.Equal(t, "05fecf43-63f5-43cd-8294-c823cd932947", response.Response.Uid)
	decoded, err := base64.StdEncoding.DecodeString(response.Response.Patch)
	assert.Nil(t, err)
	assert.Equal(t, "[{\"op\":\"add\",\"path\":\"/metadata/labels/generated\",\"value\":true}]", string(decoded))
}
