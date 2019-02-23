package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/spf13/cobra"
)

func init() {
	var inputDir string;
	var outputDir string;
	var sources = &cobra.Command{
		Use:   "source",
		Short: "List imported sources",
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.ListSources(context)
		},
	}
	sourceDestFlags(sources, &inputDir, &outputDir)
	rootCmd.AddCommand(sources)
}
