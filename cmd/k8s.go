package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/data"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command

func init() {
	var imageOverride string
	var namespaceOverride string

	var k8sCmd = &cobra.Command{
		Use:   "k8s [sourceDir] [destDir]",
		Short: "Generate k8s resource files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file := "-"
			if len(args) > 1 {
				file = args[1]
			} else {
				logrus.SetLevel(logrus.ErrorLevel)
			}
			context := data.RenderContext{
				OutputDir:     file,
				Mode:          "k8s",
				ImageOverride: imageOverride,
				Namespace:     namespaceOverride,
				InputDir:      []string{args[0],},
			}
			pkg.Run(&context)
		},
	}
	k8sCmd.Flags().StringVarP(&imageOverride, "image", "i", "", "docker image name override")
	k8sCmd.Flags().StringVarP(&namespaceOverride, "namespace", "n", "", "kubernetes namespace override")

	rootCmd.AddCommand(k8sCmd)
}
