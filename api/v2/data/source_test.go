package data

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalPath(t *testing.T) {

	source := LocalSource{
		Dir: "/tmp/test/sub/../..",
	}
	mgr := SourceCacheManager{}
	result, err := source.GetPath(&mgr)
	assert.Nil(t, err)
	assert.Equal(t, "/tmp", result)
}

func TestDoOnce(t *testing.T) {
	mgr := SourceCacheManager{}
	f1 := func() error {
		return nil
	}
	f2 := func() error {
		return errors.New("Bad")
	}
	assert.Nil(t, mgr.DoOnce("c1", f1))
	assert.Nil(t, mgr.DoOnce("c2", f1))
	assert.Nil(t, mgr.DoOnce("c1", f2))
	assert.NotNil(t, mgr.DoOnce("cx", f2))
}
