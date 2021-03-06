package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfiguration(t *testing.T) {
	configuration, _, e := ReadConfiguration("../../testdata/conf")

	assert.Nil(t, e)

	assert.Equal(t, 2, len(configuration.Import))
	assert.EqualValues(t, ImportConfiguration{Path: "../dir1"}, configuration.Import[0])
	assert.EqualValues(t, ImportConfiguration{Path: "../dir2"}, configuration.Import[1])
}
