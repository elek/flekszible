package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPath_Match(t *testing.T) {

	if !NewPath("a", "b").Match(NewPath("a", "b")) {
		t.Errorf("Path match is failed")
	}

	if NewPath("a", "b").Match(NewPath("a", "c")) {
		t.Errorf("Path match is failed")
	}

	if NewPath("a", "b").Match(NewPath("a", "b", "c")) {
		t.Errorf("Path match is failed")
	}

	if !NewPath("a", ".*").Match(NewPath("a", "b")) {
		t.Errorf("Path match is failed")
	}
}

func TestPathParent(t *testing.T) {
	p := NewPath("list", "0")
	assert.Equal(t, NewPath("list"), p.Parent())
}