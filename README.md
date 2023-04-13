# git-gerrit
Git for gerrit 简易命令行工具，辅助提交代码到gerrit上的

## 用法
```
will run: git push origin HEAD:refs/for/main 是否决定执行了？ y/N ?
```
默认会执行类似上面的push命令，输入y就会真正执行push动作了，输入其他就好退出。

origin 是默认的 git remote， 如果有多个会提示让用户选择的
main 是 默认的 git branch， 同样的分支名称也是，有多个会提示用户去选择的






