/*
Copyright © 2023 bright.ma <bright.ma@magesfc.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show this command version info",
	Long:  `show this command version info`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := ExecuteCommand("git", "version", args...)
		if err != nil {
			Error(cmd, args, err)
		}

		fmt.Fprint(os.Stdout, output)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
