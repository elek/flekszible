package cmd

import (
        "fmt"
        "github.com/spf13/cobra"
)

var version string
var commit string
var date string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
        Use:   "version",
        Short: "Print the version number of flekszible",
        Long:  `All software has versions. This is flekszible`,
        Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("Version:", version)
                fmt.Println("Build Date:", date)
                fmt.Println("Git Commit:", commit)
        },
}

func init() {
        rootCmd.AddCommand(versionCmd)
}