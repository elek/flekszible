package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/data"
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
			context := data.RenderContext{
				OutputDir:     file,
				Mode:          "helm",
				ImageOverride: imageOverride,
				InputDir:      []string{args[0],},
			}

			pkg.Run(&context)
		},
	}
	helmCmd.Flags().StringVarP(&imageOverride, "image", "i", "", "docker image name override")
	rootCmd.AddCommand(helmCmd)
}
