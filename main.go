package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	configFile string
	commandStr string
	localFile  string
	remotePath string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "walle",
		Short: "Walle 是一款高性能的分布式 Linux 批量运维与文件分发工具",
		Long:  `Walle (华力) 基于 Go 语言轻量级并发特性开发，支持通过并发 SSH 批量执行命令、SFTP 并发分发文件，支持密码与 SSH 密钥认证。`,
	}

	// 1. 批量执行命令的子命令定义
	var runCmd = &cobra.Command{
		Use:   "run [command]",
		Short: "在所有远程主机上并发执行指定的命令",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commandStr = args[0]
			cfg, err := LoadConfig(configFile)
			if err != nil {
				color.Red("❌ 加载配置文件失败: %v", err)
				os.Exit(1)
			}
			ExecuteBatchCommand(cfg, commandStr)
		},
	}

	// 2. 批量上传文件的子命令定义
	var pushCmd = &cobra.Command{
		Use:   "push [local_path] [remote_path]",
		Short: "并发分发本地文件到所有远程主机的指定路径",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			localFile = args[0]
			remotePath = args[1]
			cfg, err := LoadConfig(configFile)
			if err != nil {
				color.Red("❌ 加载配置文件失败: %v", err)
				os.Exit(1)
			}
			ExecuteBatchUpload(cfg, localFile, remotePath)
		},
	}

	// 注册全局标志 (绑定默认文件 hosts.yaml)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "hosts.yaml", "指定主机配置文件路径 (YAML 格式)")

	// 组装子命令到总入口
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(pushCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
