package examples

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/processor"
	"github.com/elek/flekszible/pkg"
	"testing"
)

func TestSource(t *testing.T) {
	data.DownloaderPlugin = pkg.GoGetterDownloader{}
	processor.TestFromDir(t, "source")
}
