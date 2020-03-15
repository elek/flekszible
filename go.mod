module github.com/elek/flekszible

require (
	github.com/brettski/go-termtables v0.0.0-20190907034855-12ddd59af020
	github.com/elek/flekszible/api v0.0.0-20190928090338-d1d34ef1bc1f
	github.com/hashicorp/go-getter v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.2.0
	github.com/stretchr/testify v1.2.2
	github.com/urfave/cli v1.22.1
	golang.org/x/crypto v0.0.0-20190927123631-a832865fa7ad // indirect
	golang.org/x/net v0.0.0-20190926025831-c00fd9afed17 // indirect
	golang.org/x/sys v0.0.0-20190927073244-c990c680b611 // indirect
	golang.org/x/text v0.3.2 // indirect
)

replace github.com/elek/flekszible/api => ./api

go 1.13
