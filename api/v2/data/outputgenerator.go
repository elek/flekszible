package data

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

			if stat.Mode().Perm()&0444 > 0 {
				//executable files are written only if doesn't exist
				if _, err = os.Stat(destinationFile); !os.IsNotExist(err) {
					destinationContent, err := ioutil.ReadFile(destinationFile)
					if err != nil {
						return resources, errors.Wrap(err, "Can't read the file from the destination: "+destinationFile)
					}
					if !bytes.Equal(content, destinationContent) {
						logrus.Warn("Sourcefile is executable and destination file exists. Doesn't overwrite it for security reason. Please delete " + destinationFile)
						return resources, nil
					}
				} else {
					logrus.Warn("Copying executable file from " + sourceFile + " to " + destinationFile)
				}

			}
			err = ioutil.WriteFile(destinationFile, content, stat.Mode())
			if err != nil {
				return resources, errors.Wrap(err, "Couldn't copy the file from "+sourceFile+" to "+destinationFile)
			}

		}
	}
	return resources, nil
}
