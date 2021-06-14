package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetName(t *testing.T) {
	node := NewMapNode(NewPath())
	metadata := node.CreateMap("metadata")
	metadata.PutValue("name", "{{ something}}-ok")

	r := Resource{Content: &node}

	assert.Equal(t, "ok", r.Name())
}
