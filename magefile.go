// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

//Test executes all unit and integration tests
func Test() error {
	fmt.Println("Executing tests")
	return sh.Run("go", "test", "./...")
}

//Test executes all unit and integration tests
func Coverage() error {
	fmt.Println("Executing tests and generate coverate information")
	err := sh.Run("go", "test", "-coverprofile=/tmp/coverage.out", "./...")
	if err != nil {
		return err
	}
	return sh.Run("go", "tool", "cover", "-html=/tmp/coverage.out")
}

//Lint executes all the linters with golangci-lint
func Lint() error {
	return sh.Run("golangci-lint", "run")
}

//Build do a standard go build
func Build() error {
	return sh.Run("go", "build", "./...")
}

//Format reformat code automatically
func Format() error {
	err := sh.Run("gofmt", "-w", ".")
	if err != nil {
		return err
	}
	return sh.Run("goimports", "-w", ".")

}

//update `./vendor` dependencies
func GenVendor() error {
	err := sh.Run("go", "mod", "tidy")
	if err != nil {
		return err
	}
	return sh.Run("go", "mod", "vendor")
}

//RegenerateProto regenerates all the protobuf related go files
func GenProto() error {
	fmt.Println("Regenerating protobuf files")
	return sh.Run("protoc", "--go_out=.", "--go_opt=paths=source_relative", "--go-grpc_out=.", "--go-grpc_opt=paths=source_relative", "proto/doit.proto")
}

//RegenerateProto regenerates all the protobuf related go files
func GenCi() error {
	fmt.Println("Regenerates ./ci file (with mage)")
	return sh.Run("mage", "-compile", "./ci")
}
