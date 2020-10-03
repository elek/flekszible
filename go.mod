module github.com/elek/flekszible

require (
	github.com/appscode/jsonpatch v1.0.1
	github.com/brettski/go-termtables v0.0.0-20190907034855-12ddd59af020
	github.com/elek/flekszible/api/v2 v2.0.0-00010101000000-000000000000
	github.com/fatih/color v1.9.0
	github.com/gin-gonic/gin v1.6.3
	github.com/hashicorp/go-getter v1.1.0
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.1
	golang.org/x/crypto v0.0.0-20190927123631-a832865fa7ad // indirect
	golang.org/x/net v0.0.0-20190926025831-c00fd9afed17 // indirect
)

replace github.com/elek/flekszible/api/v2 => ./api

go 1.13
