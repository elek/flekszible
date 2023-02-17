module github.com/elek/flekszible

require (
	github.com/appscode/jsonpatch v1.0.1
	github.com/brettski/go-termtables v0.0.0-20190907034855-12ddd59af020
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/elek/flekszible/api/v2 v2.0.0-20210614204704-5d2be8c07dc2
	github.com/fatih/color v1.12.0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-playground/validator/v10 v10.8.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/go-getter v1.7.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/magefile/mage v1.11.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ugorji/go v1.2.6 // indirect
	github.com/urfave/cli v1.22.5
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/elek/flekszible/api/v2 => ./api/v2

go 1.13
