package cmd

import (
	"github.com/elek/flekszible/pkg"
	"github.com/elek/flekszible/pkg/processor"
	"github.com/spf13/cobra"
	"os"
	"path"
)

func init() {
	var imageOverride string
	var namespaceOverride string
	var minikube bool;
	var inputDir string;
	var outputDir string;
	var k8sCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate k8s resource files",
		Run: func(cmd *cobra.Command, args []string) {
			context := processor.CreateRenderContext("k8s", findInputDir(inputDir), findOutputDir(outputDir))
			context.ImageOverride = imageOverride
			context.Namespace = namespaceOverride
			pkg.Run(context, minikube)
		},
	}
	sourceDestFlags(k8sCmd, &inputDir, &outputDir)
	k8sCmd.Flags().StringVarP(&imageOverride, "image", "i", "", "docker image name override")
	k8sCmd.Flags().StringVarP(&namespaceOverride, "namespace", "n", "<none>", "kubernetes namespace override. With empty value (\"\") the current namespace will "+
		"be used. With exact value the namespace will be forced (set even if the namespace is not added to the to the k8s resources")
	k8sCmd.Flags().BoolVarP(&minikube, "minikube", "m", false, "Enable minikube specific defaults (eg. daemonset to statefulset conversion)")
	rootCmd.AddCommand(k8sCmd)
}

func sourceDestFlags(command *cobra.Command, inputDir *string, outputDir *string) {
	command.Flags().StringVarP(inputDir, "source", "s", "", "Source directory to read the resrouces and definitions (default: ./flekszible if exists, current dir if not")
	command.Flags().StringVarP(outputDir, "destination", "d", "", "Destination directory to generate the k8s resrouces (default: current dir)")

}
func findInputDir(argument string) string {
	if argument != "" {
		return argument
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	subdir := path.Join(pwd, "flekszible")
	if _, err := os.Stat(subdir); !os.IsNotExist(err) {
		return subdir
	}
	return pwd
}

func findOutputDir(argument string) string {
	if argument != "" {
		return argument
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}