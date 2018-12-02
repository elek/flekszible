package cmd

import (
        "fmt"
        "github.com/spf13/cobra"
        "os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
        Use:   "flekszible",
        Short: "A brief description of your application",
        Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

        Run: func(cmd *cobra.Command, args []string) {
        },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
        if err := rootCmd.Execute(); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
}

func init() {
        cobra.OnInitialize()

        // Cobra also supports local flags, which will only run
        // when this action is called directly.
        // rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}