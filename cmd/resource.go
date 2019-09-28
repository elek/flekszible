package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/api/processor"
	"github.com/spf13/cobra"
)

func init() {
	var inputDir string;
	var outputDir string;
	var resources = &cobra.Command{
		Use:   "resource",
		Short: "List processed k8s resources",
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.ListResources(context)
		},
	}
	sourceDestFlags(rootCmd, &inputDir, &outputDir)
	rootCmd.AddCommand(resources)
}

