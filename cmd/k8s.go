package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command

func init() {
	var imageOverride string
	var namespaceOverride string
	var minikube bool;
	var k8sCmd = &cobra.Command{
		Use:   "generate [sourceDir] [destDir]",
		Short: "Generate k8s resource files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			file := "-"
			if len(args) > 1 {
				file = args[1]
			} else {
				logrus.SetLevel(logrus.ErrorLevel)
			}
			context := processor.CreateRenderContext("k8s", args[0], file)
			context.ImageOverride = imageOverride
			context.Namespace = namespaceOverride
			pkg.Run(context, minikube)
		},
	}
	k8sCmd.Flags().StringVarP(&imageOverride, "image", "i", "", "docker image name override")
	k8sCmd.Flags().StringVarP(&namespaceOverride, "namespace", "n", "", "kubernetes namespace override")
	k8sCmd.Flags().BoolVarP(&minikube, "minikube", "m", false, "Enable minikube specific defaults (eg. daemonset to statefulset conversion)")
	rootCmd.AddCommand(k8sCmd)
}
