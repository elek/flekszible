package version

import (
        "runtime"
        "fmt"
)

var GitCommit string

var BuildDate = ""

var GoVersion = runtime.Version()

var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)