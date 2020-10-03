package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetName(t *testing.T) {
	node := NewMapNode(NewPath())
	metadata := node.CreateMap("metadata")
	metadata.PutValue("name", "{{ something}}-ok")

	r := Resource{Content: &node}

	assert.Equal(t, "ok", r.Name())
}
