// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"time"
)

var Default = Build

var packageName = "github.com/elek/flekszible"
var ldflags = "-X $PACKAGE/version.BuildDate=$BUILD_DATE -X $PACKAGE/version.GitCommit=$COMMIT_HASH"

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"PACKAGE":     packageName,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}

func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	fmt.Println(flagEnv())
	return sh.RunWith(flagEnv(), "go", "build", "-ldflags", ldflags, "-o", "target/bin/flekszible", ".")
}

func InstallDeps() error {
	fmt.Println("Installing Deps...")
	return sh.Run("go", "mod", "download")
}

func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("target")
}

func Test() error {
	mg.Deps(Clean)
	fmt.Println("Executing tests")
	return sh.RunV("go", "test", "./...")
}

func UpdateBuilder() error {
	return sh.RunV("mage", "-compile", "build")
}
