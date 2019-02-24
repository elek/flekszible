package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/spf13/cobra"
)

func init() {
	var inputDir string;
	var outputDir string;
	var app = &cobra.Command{
		Use:   "app",
		Short: "List active flekszible apps (dirs)",
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.ListApp(context)
		},
	}
	var search = &cobra.Command{
		Use:   "search",
		Short: "Search for available importable apps from the active sources.",
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.SearchComponent(context)
		},
	}
	sourceDestFlags(search, &inputDir, &outputDir)
	app.AddCommand(search)
	sourceDestFlags(app, &inputDir, &outputDir)
	rootCmd.AddCommand(app)
}
