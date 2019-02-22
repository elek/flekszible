package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/spf13/cobra"
	"os"
)

func init() {

	var resources = &cobra.Command{
		Use:   "resource [sourceDir] [destDir]",
		Short: "List processed k8s resources",
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", defaultInputDir(), defaultInputDir())

			pkg.ListResources(context)
		},
	}
	rootCmd.AddCommand(resources)
}

func defaultInputDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}
