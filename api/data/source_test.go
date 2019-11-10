package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalPath(t *testing.T) {

	source := LocalSource{
		Dir:         "/tmp/test/sub",
		RelativeDir: "../../",
	}
	mgr := SourceCacheManager{}
	result, err := source.GetPath(&mgr, "dir")
	assert.Nil(t, err)
	assert.Equal(t, "/tmp/dir", result)
}
