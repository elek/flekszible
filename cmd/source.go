package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/public/processor"
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

	var add = &cobra.Command{
		Use:   "add",
		Short: "Add source to the flekszible.yaml definition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.AddSource(context, findInputDir(inputDir), args[0])
		},
	}

	var search = &cobra.Command{
		Use:   "search",
		Short: "Search for the available source repositoris (github repositories with flekszible:topic)",
		Run: func(cmd *cobra.Command, args []string) {
			pkg.SearchSource()
		},
	}

	sourceDestFlags(sources, &inputDir, &outputDir)
	sourceDestFlags(add, &inputDir, &outputDir)
	rootCmd.AddCommand(sources)
	sources.AddCommand(add)
	sources.AddCommand(search)
}
