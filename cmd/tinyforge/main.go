package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/tinyforge-cn/cli/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "login":
		err = commands.LoginCommand(os.Args[2:])
	case "submit":
		err = commands.SubmitCommand(os.Args[2:])
	case "stages":
		err = commands.StagesCommand(os.Args[2:])
	case "test":
		err = commands.TestCommand(os.Args[2:])
	case "version", "--version", "-v":
		commands.VersionCommand()
		return
	case "help", "--help", "-h":
		printUsage()
		return
	default:
		fmt.Fprintf(os.Stderr, "未知命令: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		fmt.Fprintf(os.Stderr, "错误: %s\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`tinyforge — Tinyforge 平台 CLI 工具

用法:
  tinyforge <命令> [选项]

命令:
  login     登录 Tinyforge 平台
  submit    提交代码进行评测
  test      本地运行 tester 评测
  stages    列出当前课程的所有关卡 slug
  version   显示版本信息

选项:
  --help    显示帮助信息

示例:
  tinyforge login                     浏览器授权登录
  tinyforge login --token tcs_xxx      使用 Token 登录
  tinyforge submit                    提交当前目录代码
  tinyforge submit --stage softmax    指定评测关卡
  tinyforge submit --dry-run          仅预览打包内容
  tinyforge stages                    列出当前课程所有关卡 slug
  tinyforge test                      本地评测当前关卡
  tinyforge test --stage softmax      本地评测指定关卡
  tinyforge test --all                本地评测所有关卡`)
}
