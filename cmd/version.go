/*
Copyright © 2023 bright.ma <bright.ma@magesfc.com>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// VERSION go build -ldflags "-X cmd.VERSION=x.x.x"
var Version = "not specified"
var Commit = "not specified"

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
	fmt.Printf("%s", output)
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Commit: %s\n", Commit)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
