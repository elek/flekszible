package data

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

type OutputGenerator struct {
}

func (*OutputGenerator) IsManagedDir(dir string) bool {
	return path.Base(dir) == "output"
}

func (*OutputGenerator) Generate(sourceDir string, destinationDir string) ([]*Resource, error) {
	resources := make([]*Resource, 0)
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return resources, err
	}
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			sourceFile := path.Join(sourceDir, filename)
			content, err := ioutil.ReadFile(sourceFile)
			if err != nil {
				return resources, errors.Wrap(err, "Can't read the file to copy to the destination: "+sourceFile)
			}
			stat, err := os.Stat(sourceFile)

			destinationFile := path.Join(destinationDir, filename)
			err = ioutil.WriteFile(destinationFile, content, stat.Mode())
			if err != nil {
				return resources, errors.Wrap(err, "Couldn't copy the file from "+sourceFile+" to "+destinationFile)
			}

		}
	}
	return resources, nil
}
