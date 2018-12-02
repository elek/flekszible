package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfiguration(t *testing.T) {
	configuration, e := ReadConfiguration("../../testdata/conf/flekszible.yaml")

	assert.Nil(t, e)

	assert.EqualValues(t, []string{"../dir1", "../dir2"}, configuration.Import)
}
