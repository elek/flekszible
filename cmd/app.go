package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/public/processor"
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
	var add = &cobra.Command{
		Use:   "add",
		Short: "Add (import) new app to the flekszible.yaml.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.AddApp(context, findInputDir(inputDir), args[0])
		},
	}
	sourceDestFlags(search, &inputDir, &outputDir)
	sourceDestFlags(add, &inputDir, &outputDir)
	app.AddCommand(search)
	app.AddCommand(add)
	sourceDestFlags(app, &inputDir, &outputDir)
	rootCmd.AddCommand(app)
}
