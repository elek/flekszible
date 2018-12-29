package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "flekszible",
	Short: "A highly flexible kubernetes resource manager",
	Long:  ``,
}

func Execute() {

	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return rootCmd.Help()
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

}
