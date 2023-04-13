/*
Copyright © 2023 bright.ma <bright.ma@magesfc.com>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-gerrit",
	Short: "git gerrit cli",
	Long:  `git gerrit cli`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var defCmd string
	var cmdFound bool
	var cmdHelp bool
	defCmd = "push"
	cmd := rootCmd.Commands()

	for _, a := range cmd {
		for _, b := range os.Args[1:] {
			if a.Name() == b {
				cmdFound = true
				break
			}
		}
	}
	// 这里判断 参数中是否 有help这个子命令, 如果有help 就还需要执行 原生的
	for _, b := range os.Args[1:] {
		if "help" == b {
			cmdHelp = true
			break
		}
	}
	// 如果参数中没有help，也没有其他 子命令，那么久执行默认的 push 子命令
	if !cmdFound && !cmdHelp {
		args := append([]string{defCmd}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git-gerrit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

}
