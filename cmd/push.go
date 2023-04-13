/*
Copyright © 2023 bright.ma <bright.ma@magesfc.com>

*/
package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/spf13/cobra"
)

var Branch string
var Topic string
var Draft bool
var Bypass bool

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push to gerrit",
	Long:  `push to gerrit`,
	Run: func(cmd *cobra.Command, args []string) {
		push(cmd, args)
	},
}

type RemoteOption struct {
	Name string
	URL  string
}

func getRemote(cmd *cobra.Command, args []string) RemoteOption {
	output, err := ExecuteCommand("git", "remote", "-v")
	if err != nil {
		Error(cmd, args, err)
	}
	if output == "" {
		err = errors.New("no remote")
		Error(cmd, args, err)
	}
	//解析所有的remote，然后保存成 RemoteOption 切片
	remotes := strings.Split(output, "\n")

	remote_options := make([]RemoteOption, 0, len(remotes))
	for _, v := range remotes {
		if v != "" && strings.Contains(v, "(push)") {
			words := strings.Fields(v)
			remote_options = append(remote_options, RemoteOption{
				Name: words[0],
				URL:  words[1],
			})
		}
	}

	remote_options_len := len(remote_options)

	var remoteOption RemoteOption
	if remote_options_len > 1 {
		for i, v := range remote_options {
			result := fmt.Sprintf("[%d]\t%s\t%s", i, v.Name, v.URL)
			fmt.Println(result)
		}
		fmt.Println("请输入上面数字选择push的远端地址：")
		var index int
		fmt.Scanln(&index)
		if index < 0 || index >= remote_options_len {
			err = errors.New("输入数字非法")
			Error(cmd, args, err)
		}
		remoteOption = remote_options[index]
	} else {
		remoteOption = remote_options[0]
	}
	return remoteOption
}

type BranchOption struct {
	Name string
}

func getBranch(cmd *cobra.Command, args []string, remoteOption RemoteOption) BranchOption {
	if Branch != "" {
		branchOption := BranchOption{
			Name: Branch,
		}
		return branchOption
	}
	output, err := ExecuteCommand("git", "branch", "-a")
	if err != nil {
		Error(cmd, args, err)
	}
	if output == "" {
		err = errors.New("no branch")
		Error(cmd, args, err)
	}
	//解析所有的 branch，然后保存成 BranchOption 切片
	branchs := strings.Split(output, "\n")
	branch_options := make([]BranchOption, 0, len(branchs))

	for _, v := range branchs {
		if v == "" {
			continue
		}
		if strings.Contains(v, "*") {
			words := strings.Fields(v)
			if len(words) == 2 {
				branch_options = append(branch_options, BranchOption{
					Name: words[1],
				})
				break
			}
			keyword := fmt.Sprintf("%s/", remoteOption.Name)
			//英文语言下的 括号
			if len(words) > 2 && strings.Contains(v, keyword) {
				if strings.Contains(v, "(") && strings.Contains(v, ")") {
					branch := parseDetached(v, ")", remoteOption)
					branch_options = append(branch_options, BranchOption{
						Name: branch,
					})
					break
				}
				//中文语言下的 括号
				if len(words) > 2 && strings.Contains(v, "（") && strings.Contains(v, "）") {
					branch := parseDetached(v, "）", remoteOption)
					branch_options = append(branch_options, BranchOption{
						Name: branch,
					})
					break
				}
			}

		}
	}

	branch_options_len := len(branch_options)
	if branch_options_len == 1 {
		return branch_options[0]
	}

	for _, v := range branchs {
		// 单独克隆的: gerrit上克隆的仓库比较特殊，会有一个 HEAD 指向  remotes/origin/HEAD -> origin/master
		keyword1 := fmt.Sprintf("remotes/%s/HEAD", remoteOption.Name)
		if strings.Contains(v, keyword1) {
			branch := parseSpecRef(v, remoteOption)
			branch_options = append(branch_options, BranchOption{
				Name: branch,
			})
			break
		}
		// 使用repo下载的仓库比较特殊，会有一个  remotes/m/  指向  remotes/m/dev -> origin/dev
		keyword2 := "remotes/m/"
		if strings.Contains(v, keyword2) {
			branch := parseSpecRef(v, remoteOption)
			branch_options = append(branch_options, BranchOption{
				Name: branch,
			})
			break
		}
		keyword3 := fmt.Sprintf("remotes/%s/", remoteOption.Name)
		if strings.Contains(v, keyword3) {
			branch := parseSpecRef(v, remoteOption)
			branch = strings.TrimPrefix(branch, keyword3)
			branch_options = append(branch_options, BranchOption{
				Name: branch,
			})
		}
	}

	branch_options_len = len(branch_options)

	var branchOption BranchOption
	if branch_options_len > 1 {
		for i, v := range branch_options {
			result := fmt.Sprintf("[%d]\t%s", i, v.Name)
			fmt.Println(result)
		}
		fmt.Println("请输入上面数字选择push的远端分支名称：")
		var index int
		fmt.Scanln(&index)
		if index < 0 || index >= branch_options_len {
			err = errors.New("输入数字非法")
			Error(cmd, args, err)
		}
		branchOption = branch_options[index]
	} else if branch_options_len == 1 {
		branchOption = branch_options[0]
	} else {
		err = errors.New("没有任何分支名称, 请使用-b选项指定分支名")
		Error(cmd, args, err)
	}
	return branchOption
}

func parseSpecRef(v string, remoteOption RemoteOption) string {
	words := strings.Fields(v)
	prefix := fmt.Sprintf("%s/", remoteOption.Name)
	branch := strings.TrimPrefix(words[len(words)-1], prefix)
	return branch
}

func parseDetached(v string, old string, remoteOption RemoteOption) string {
	words := strings.Fields(v)
	prefix := fmt.Sprintf("%s/", remoteOption.Name)
	branch := strings.TrimPrefix(words[len(words)-1], prefix)
	branch = strings.Replace(branch, old, "", -1)
	return branch
}

func push(cmd *cobra.Command, args []string) {
	remoteOption := getRemote(cmd, args)
	branchOption := getBranch(cmd, args, remoteOption)
	s := fmt.Sprintf("HEAD:refs/for/%s", branchOption.Name)
	if Draft {
		s = fmt.Sprintf("HEAD:refs/drafts/%s", branchOption.Name)
	}
	if Topic != "" {
		s = fmt.Sprintf("%s%%topic=%s", s, Topic)
	}
	if Bypass {
		s = fmt.Sprintf("HEAD:refs/heads/%s", branchOption.Name)
	}

	fmt.Println("will run: git push", remoteOption.Name, s, "是否决定执行了？ y/N ?")
	var yes string
	fmt.Scanln(&yes)
	if yes == "y" || yes == "Y" {
		fmt.Println("will run: git push ")
		output, err := ExecuteCommand("git", "push", remoteOption.Name, s)
		fmt.Println(output)
		if err != nil {
			Error(cmd, args, err)
		}
	}
}

func init() {
	pushCmd.Flags().StringVarP(&Branch, "branch", "b", "", "what remote branch want to push")
	pushCmd.Flags().StringVarP(&Topic, "topic", "t", "", "push to gerrit with topic")
	pushCmd.Flags().BoolVarP(&Draft, "draft", "d", false, "push to gerrit as drafts")
	pushCmd.Flags().BoolVarP(&Bypass, "bypass", "p", false, "push to gerrit directly")

	rootCmd.AddCommand(pushCmd)

}
