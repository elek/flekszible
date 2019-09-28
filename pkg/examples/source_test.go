package examples

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/processor"
	"github.com/elek/flekszible/pkg"
	"testing"
)

func TestSource(t *testing.T) {
	data.DownloaderPlugin = pkg.GoGetterDownloader{}
	processor.TestFromDir(t, "source")
}
