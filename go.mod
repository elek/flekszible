module github.com/elek/flekszible

require (
	github.com/appscode/jsonpatch v1.0.1
	github.com/brettski/go-termtables v0.0.0-20190907034855-12ddd59af020
	github.com/elek/flekszible/api/v2 v2.0.0-20201003124246-f897a6d7ee91
	github.com/fatih/color v1.9.0
	github.com/gin-gonic/gin v1.6.3
	github.com/hashicorp/go-getter v1.1.0
	github.com/magefile/mage v1.11.0 // indirect
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.1
	golang.org/x/tools v0.1.1 // indirect
)

replace github.com/elek/flekszible/api/v2 => ./api/v2

go 1.13
