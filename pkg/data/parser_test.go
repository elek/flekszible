package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadFile(t *testing.T) {
	node, err := ReadFile("../../testdata/parser/datanode.yaml")
	assert.Nil(t, err)
	node.Accept(PrintVisitor{})
}
