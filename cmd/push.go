/*
Copyright © 2023 bright.ma <bright.ma@magesfc.com>
*/
package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

var Branch string
var Topic string
var Hashtags string
var Message string
var Label string

var Wip bool
var Ready bool

var Private bool
var RemovePrivate bool

var PublishComments bool
var NoPublishComments bool

var Edit bool

var RefsHeads bool

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
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "x {{ .Name | red }}\t({{ .URL | blue }})",
			Inactive: "  {{ .Name }}\t({{ .URL  }})",
			Selected: "you select this remote: {{ .Name | green }} ({{ .URL | blue }})",
		}

		prompt := promptui.Select{
			Label:     "Select Remote",
			Items:     remote_options,
			Templates: templates,
			Size:      remote_options_len,
		}
		chooseIndex, _, err := prompt.Run()

		if err != nil {
			Error(cmd, args, err)
		}
		remoteOption = remote_options[chooseIndex]
	} else if remote_options_len == 1 {
		remoteOption = remote_options[0]
	} else {
		err = errors.New("没有任何 remote 名称, 请使用git remote add 添加")
		Error(cmd, args, err)
	}
	return remoteOption
}

type BranchOption struct {
	Name string
	Desc string
}

func getBranch(cmd *cobra.Command, args []string, remoteOption RemoteOption) BranchOption {
	// 分支查找顺序，第一是 本地分支名。第二个是 HEAD指向的，或者 /m/ 指向的。第三是其他remote匹配的分支名
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

	// 本地 新建的分支 解析出来的,如果解析出来的有一个值那么就优先使用这个，然后直接就是 返回了
	local_branch_options := getLocalBranchs(branchs, remoteOption)
	if len(local_branch_options) == 1 {
		return local_branch_options[0]
	}

	// 2个特殊分支解析出来的 这种 特殊的分支名称，只可能有 1个
	special_branch_options := getSpecialBranchs(branchs, remoteOption)
	// 所有 remotes 解析出来的分支 名称
	remtoes_branch_options := getRemotesBranchs(branchs, remoteOption)
	remtoes_branch_options = removeSpecialBranchs(remtoes_branch_options, special_branch_options)

	branch_options := make([]BranchOption, 0, len(branchs))
	branch_options = append(branch_options, special_branch_options...)
	branch_options = append(branch_options, remtoes_branch_options...)

	branch_options_len := len(branch_options)

	var branchOption BranchOption
	if branch_options_len > 1 {
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "x {{ .Name | red }} {{ .Desc | red }}",
			Inactive: "  {{ .Name }} {{ .Desc }}",
			Selected: "you select this branch: {{ .Name | green }}",
		}

		prompt := promptui.Select{
			Label:     "Select Branch",
			Items:     branch_options,
			Templates: templates,
			Size:      branch_options_len,
		}
		chooseIndex, _, err := prompt.Run()

		if err != nil {
			Error(cmd, args, err)
		}
		branchOption = branch_options[chooseIndex]
	} else if branch_options_len == 1 {
		branchOption = branch_options[0]
	} else {
		err = errors.New("没有任何分支名称, 请使用-b选项指定分支名")
		Error(cmd, args, err)
	}
	return branchOption
}

func getSpecialBranchs(branchs []string, remoteOption RemoteOption) []BranchOption {
	// 2个特殊分支解析出来的
	branch_options := make([]BranchOption, 0, len(branchs))
	for _, v := range branchs {
		// 单独克隆的: gerrit上克隆的仓库比较特殊，会有一个 HEAD 指向  remotes/origin/HEAD -> origin/master
		keyword1 := fmt.Sprintf("remotes/%s/HEAD", remoteOption.Name)
		if strings.Contains(v, keyword1) {
			branch := parseSpecRef(v, remoteOption)
			branch_options = append(branch_options, BranchOption{
				Name: branch,
				Desc: v,
			})
		}
		// 使用repo下载的仓库比较特殊，会有一个  remotes/m/  指向  remotes/m/dev -> origin/dev
		keyword2 := "remotes/m/"
		if strings.Contains(v, keyword2) {
			branch := parseSpecRef(v, remoteOption)
			branch_options = append(branch_options, BranchOption{
				Name: branch,
				Desc: v,
			})
		}
	}
	return branch_options
}

func removeSpecialBranchs(firstS, secondS []BranchOption) []BranchOption {
	// 遍历第一个切片，删除 存在于第二个切片中的元素
	branch_options := make([]BranchOption, 0, len(firstS))
	for _, v := range firstS {
		if !inSlice(secondS, v) {
			branch_options = append(branch_options, v)
		}
	}
	return branch_options
}

func inSlice(items []BranchOption, e BranchOption) bool {
	for _, v := range items {
		if v.Name == e.Name {
			return true
		}
	}
	return false
}

func getRemotesBranchs(branchs []string, remoteOption RemoteOption) []BranchOption {
	// 所有 remotes 解析出来的分支 名称
	branch_options := make([]BranchOption, 0, len(branchs))
	for _, v := range branchs {
		keyword3 := fmt.Sprintf("remotes/%s/", remoteOption.Name)
		if strings.Contains(v, keyword3) {
			branch := parseSpecRef(v, remoteOption)
			branch = strings.TrimPrefix(branch, keyword3)
			branch_options = append(branch_options, BranchOption{
				Name: branch,
				Desc: v,
			})
		}
	}
	return branch_options
}

func getLocalBranchs(branchs []string, remoteOption RemoteOption) []BranchOption {
	// 本地 新建的分支 解析出来的
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
	return branch_options
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

	pushArgs := make([]string, 0)
	pushArgs = append(pushArgs, remoteOption.Name)

	if RefsHeads {
		refsPattern := fmt.Sprintf("HEAD:refs/%s/%s", "heads", branchOption.Name)
		pushArgs = append(pushArgs, refsPattern)
	} else {
		refsPattern := fmt.Sprintf("HEAD:refs/%s/%s", "for", branchOption.Name)
		pushArgs = append(pushArgs, refsPattern)

		if Topic != "" {
			s := fmt.Sprintf("topic=%s", Topic)
			pushArgs = append(pushArgs, "-o", s)
		}

		if Hashtags != "" {
			s := fmt.Sprintf("hashtag=%s", Hashtags)
			pushArgs = append(pushArgs, "-o", s)
		}

		if Message != "" {
			s := fmt.Sprintf("message=%s", Message)
			pushArgs = append(pushArgs, "-o", s)
		}

		if Label != "" {
			s := fmt.Sprintf("label=%s", Label)
			pushArgs = append(pushArgs, "-o", s)
		}

		if Private {
			pushArgs = append(pushArgs, "-o", "private")
		} else if RemovePrivate {
			pushArgs = append(pushArgs, "-o", "remove-private")
		}

		if Wip {
			pushArgs = append(pushArgs, "-o", "wip")
		} else if Ready {
			pushArgs = append(pushArgs, "-o", "ready")
		}

		if PublishComments {
			pushArgs = append(pushArgs, "-o", "publish-comments")
		} else if NoPublishComments {
			pushArgs = append(pushArgs, "-o", "no-publish-comments")
		}

		if Edit {
			pushArgs = append(pushArgs, "-o", "edit")
		}
	}

	pushString := strings.Join(pushArgs, " ")

	label := fmt.Sprintf("%s %s %s", "will run: git push", pushString, "是否决定执行了")
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	yes, err := prompt.Run()

	if err != nil {
		Error(cmd, args, err)
	}
	fmt.Printf("You choose %q\n", yes)

	if yes == "y" || yes == "Y" {
		fmt.Println("will run: git push", pushString)
		output, err := CaptureCommand("git", "push", pushArgs...)
		if err != nil {
			ParseOutput(output)
			Error(cmd, args, err)
		}
	}
}

func ParseOutput(output string) {
	//todo
}

func init() {
	pushCmd.Flags().StringVarP(&Branch, "branch", "b", "", "what remote branch want to push")
	pushCmd.Flags().StringVarP(&Topic, "topic", "t", "", "push to gerrit with topic")
	pushCmd.Flags().StringVarP(&Hashtags, "hashtags", "g", "", "push to gerrit with hashtags")
	pushCmd.Flags().StringVarP(&Message, "message", "m", "", "push to gerrit with Patch Set Description")
	pushCmd.Flags().StringVarP(&Label, "label", "l", "", "push to gerrit with Review labels \nex: Code-Review+1,l=Verified+1")

	pushCmd.Flags().BoolVarP(&RefsHeads, "heads", "H", false, "push to gerrit refs/heads/ directly")

	pushCmd.Flags().BoolVarP(&Edit, "edit", "E", false, "push to gerrit Change Edits, \nedit is not supported for new changes")

	pushCmd.Flags().BoolVarP(&Wip, "wip", "W", false, "push a Work-In-Progress change")
	pushCmd.Flags().BoolVarP(&Ready, "remove-wip", "", false, "push to remove the wip flag")

	pushCmd.Flags().BoolVarP(&Private, "private", "P", false, "push to a private change")
	pushCmd.Flags().BoolVarP(&RemovePrivate, "remove-private", "", false, "push to remove the private flag ")

	pushCmd.Flags().BoolVarP(&PublishComments, "publish-comments", "C", false, "push to gerrit with Publish Draft Comments")
	pushCmd.Flags().BoolVarP(&NoPublishComments, "no-publish-comments", "", false, "push to gerrit with No Publish Draft Comments")

	rootCmd.AddCommand(pushCmd)

}
