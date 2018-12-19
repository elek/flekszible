package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfiguration(t *testing.T) {
	configuration, e := ReadConfiguration("../../testdata/conf/flekszible.yaml")

	assert.Nil(t, e)

	assert.Equal(t, 2, len(configuration.Import))
	assert.EqualValues(t, ImportConfiguration{Path: "../dir1"}, configuration.Import[0])
	assert.EqualValues(t, ImportConfiguration{Path: "../dir2"}, configuration.Import[1])
}
