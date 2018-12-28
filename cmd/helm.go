package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command

func init() {
	var imageOverride string
	var helmCmd = &cobra.Command{
		Use:   "helm [sourceDir] [destDir]",
		Short: "Generate helm chart",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file := "-"
			if len(args) > 1 {
				file = args[1]
			}
			context := processor.CreateRenderContext("helm", args[0], file)
			context.ImageOverride = imageOverride

			pkg.Run(context)
		},
	}
	helmCmd.Flags().StringVarP(&imageOverride, "image", "i", "", "docker image name override")
	rootCmd.AddCommand(helmCmd)
}
