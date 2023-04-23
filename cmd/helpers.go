/*
Copyright © 2023 bright.ma <bright.ma@magesfc.com>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func ExecuteCommand(name string, subname string, args ...string) (string, error) {
	args = append([]string{subname}, args...)

	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()

	return string(bytes), err
}

// Capture output but also show progress #3
// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func CaptureCommand(name string, subname string, args ...string) (string, error) {
	args = append([]string{subname}, args...)

	cmd := exec.Command(name, args...)
	//bytes, err := cmd.CombinedOutput()

	var stdBuffer bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdBuffer)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stdBuffer)

	err := cmd.Run()

	return string(stdBuffer.Bytes()), err
}

func Error(cmd *cobra.Command, args []string, err error) {
	fmt.Printf("execute %s args:%v error:%v\n", cmd.Name(), args, err)
	os.Exit(1)
}
