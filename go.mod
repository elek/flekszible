module github.com/elek/flekszible

require (
	github.com/appscode/jsonpatch v1.0.1
	github.com/brettski/go-termtables v0.0.0-20190907034855-12ddd59af020
	github.com/elek/flekszible/api/v2 v2.0.0-20201003124246-f897a6d7ee91
	github.com/fatih/color v1.9.0
	github.com/gin-gonic/gin v1.7.0
	github.com/hashicorp/go-getter v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.1
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
)

replace github.com/elek/flekszible/api/v2 => ./api/v2

go 1.13
