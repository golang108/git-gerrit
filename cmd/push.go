/*
Copyright © 2023 bright.ma <bright.ma@magesfc.com>

*/
package cmd

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

var Branch string

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push to gerrit",
	Long:  `push to gerrit`,
	Run: func(cmd *cobra.Command, args []string) {
		push(cmd, args)
	},
}

func push(cmd *cobra.Command, args []string) {
	output, err := ExecuteCommand("git", "remote", "-v")
	if err != nil {
		Error(cmd, args, err)
	}
	if output == "" {
		err = errors.New("no remote")
		Error(cmd, args, err)
	}

	output, err = ExecuteCommand("git", "branch", "-v")
	if err != nil {
		Error(cmd, args, err)
	}

	fmt.Println("push called", args, Branch, output)
}

func init() {
	pushCmd.Flags().StringVarP(&Branch, "Branch", "b", "", "the remote git branch")

	rootCmd.AddCommand(pushCmd)

}
