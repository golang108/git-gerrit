# git-gerrit
Git for gerrit 简易命令行工具，辅助提交代码到gerrit上的

## 用法
```
will run: git push origin HEAD:refs/for/main 是否决定执行了？ y/N ?
```
默认会执行类似上面的push命令，输入y就会真正执行push动作了，输入其他就退出。

origin 是默认的 git remote， 如果有多个会提示让用户选择的

main 是 默认的 git branch， 同样的分支名称也是，有多个会提示用户去选择的



## push子命令

默认不加任何子命令的话 就会是 push子命令，等于 git-gerrit push 
```
$ ./git-gerrit -h    
push to gerrit

Usage:
  git-gerrit push [flags]

Flags:
  -b, --branch string   what remote branch want to push
  -D, --draft           push to gerrit as drafts
  -H, --heads           push to gerrit directly
  -h, --help            help for push
  -t, --topic string    push to gerrit with topic


```

## help 子命令
```
$ ./git-gerrit help 
git gerrit cli

Usage:
  git-gerrit [command]

Available Commands:
  help        Help about any command
  push        push to gerrit
  version     show this command version info

Flags:
  -h, --help   help for git-gerrit

Use "git-gerrit [command] --help" for more information about a command.

```

## 作为git子命令来使用

把这个 可执行文件放到环境变量 PATH 路径下面 即可

然后就可以这样执行了
```
$ git gerrit -h 

```










