package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/spf13/cobra"
)

func init() {
	var inputDir string;
	var outputDir string;
	var list = &cobra.Command{
		Use:     "definition",
		Aliases: []string{"processor", "def"},
		Short:   "List available processor definitions",
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.ListProcessor(context)
		},
	}
	var show = &cobra.Command{
		Use:   "show",
		Short: "Show details of a specific processor definition",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			pkg.ShowProcessor(context, args[0])
		},
	}
	sourceDestFlags(list, &inputDir, &outputDir)
	sourceDestFlags(show, &inputDir, &outputDir)
	list.AddCommand(show)
	rootCmd.AddCommand(list)
}
