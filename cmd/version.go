/*
Copyright © 2023 bright.ma <bright.ma@magesfc.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show this command version info",
	Long:  `show this command version info`,
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

func showVersion() {
	output, _ := ExecuteCommand("git", "version")
	fmt.Printf("git version %s\n", output)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
