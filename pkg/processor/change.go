package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"regexp"
)

type Change struct {
	DefaultProcessor
	Pattern     string
	Replacement string
}

func (change *Change) ProcessKey(path data.Path, value interface{}) interface{} {
	re, err := regexp.Compile(change.Pattern)
	if err != nil {
		panic(err)
	}
	return re.ReplaceAllString(value.(string), change.Replacement)

}

func init() {
	prototype := Change{}
	ProcessorTypeRegistry.Add(&prototype)
}
