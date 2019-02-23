package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/spf13/cobra"
)

func init() {
	var inputDir string;
	var outputDir string;
	var component = &cobra.Command{
		Use:   "component",
		Short: "List active flekszible components (dirs)",
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.ListComponent(context)
		},
	}
	sourceDestFlags(component, &inputDir, &outputDir)
	rootCmd.AddCommand(component)
}
