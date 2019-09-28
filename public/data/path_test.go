package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathMatch(t *testing.T) {

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

func TestPathMatchLimited(t *testing.T) {

	m, s := NewPath("a", "b", "c").MatchLimited(NewPath("a", "b"))
	assert.True(t, m)
	assert.Equal(t, s, "c")

	m, s = NewPath("a", "b", "c").MatchLimited(NewPath("a", "b", "c"))
	assert.False(t, m)
}

func TestPathParent(t *testing.T) {
	p := NewPath("list", "0")
	assert.Equal(t, NewPath("list"), p.Parent())
}
