package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// connectToHost 建立支持密码和密钥智能切换的 SSH 连接
func connectToHost(host HostInfo) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	// 1. 优先选用私钥认证
	keyPath := host.PrivateKeyPath
	if keyPath != "" {
		if strings.HasPrefix(keyPath, "~/") {
			home, _ := os.UserHomeDir()
			keyPath = filepath.Join(home, keyPath[2:])
		}
		keyBytes, err := os.ReadFile(keyPath)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(keyBytes)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
	}

	// 2. 其次使用密码认证
	if host.Password != "" {
		authMethods = append(authMethods, ssh.Password(host.Password))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("没有提供有效的密码或私钥认证信息")
	}

	config := &ssh.ClientConfig{
		User:            host.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 自动化运维脚本常忽略，生产环境若有需要可强化校验
		Timeout:         7 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host.Host, host.Port)
	return ssh.Dial("tcp", addr, config)
}

// ExecuteBatchCommand 批量并发执行命令并美化输出
func ExecuteBatchCommand(cfg *Config, cmdStr string) {
	var wg sync.WaitGroup
	startTime := time.Now()

	color.Cyan("🚀 正在向 %d 台远程主机分发执行命令: %s", len(cfg.Hosts), cmdStr)
	fmt.Println(strings.Repeat("-", 60))

	for _, h := range cfg.Hosts {
		wg.Add(1)
		go func(host HostInfo) {
			defer wg.Done()

			client, err := connectToHost(host)
			if err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 建立 SSH 连接失败: %v\n", host.Alias, err)
				return
			}
			defer client.Close()

			session, err := client.NewSession()
			if err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 创建 SSH 会话失败: %v\n", host.Alias, err)
				return
			}
			defer session.Close()

			output, err := session.CombinedOutput(cmdStr)

			if err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 执行失败 (Err: %v):\n%s\n", host.Alias, err, string(output))
			} else {
				color.New(color.FgGreen).Printf("[%s] 主机命令执行成功:\n", host.Alias)
				fmt.Println(string(output))
			}
			fmt.Println(strings.Repeat("-", 40))
		}(h)
	}

	wg.Wait()
	color.Cyan("✨ 所有命令执行完毕，总耗时: %v", time.Since(startTime))
}

// ExecuteBatchUpload 批量并发分发文件
func ExecuteBatchUpload(cfg *Config, localPath, remotePath string) {
	var wg sync.WaitGroup
	startTime := time.Now()

	color.Cyan("📂 正在向 %d 台主机并发分发文件: [本地] %s -> [远程] %s", len(cfg.Hosts), localPath, remotePath)
	fmt.Println(strings.Repeat("-", 60))

	for _, h := range cfg.Hosts {
		wg.Add(1)
		go func(host HostInfo) {
			defer wg.Done()

			client, err := connectToHost(host)
			if err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 连接失败: %v\n", host.Alias, err)
				return
			}
			defer client.Close()

			sftpClient, err := sftp.NewClient(client)
			if err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 启动 SFTP 子系统失败: %v\n", host.Alias, err)
				return
			}
			defer sftpClient.Close()

			srcFile, err := os.Open(localPath)
			if err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 打开本地文件错误: %v\n", host.Alias, err)
				return
			}
			defer srcFile.Close()

			remoteDir := filepath.Dir(remotePath)
			if err := sftpClient.MkdirAll(remoteDir); err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 远程创建目录失败: %v\n", host.Alias, err)
				return
			}

			dstFile, err := sftpClient.Create(remotePath)
			if err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 远程创建文件失败: %v\n", host.Alias, err)
				return
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				color.New(color.FgRed).Printf("[%s] ❌ 传输数据中断: %v\n", host.Alias, err)
			} else {
				color.New(color.FgGreen).Printf("[%s] 📂 文件成功同步至 -> %s\n", host.Alias, remotePath)
			}
		}(h)
	}

	wg.Wait()
	color.Cyan("✨ 所有文件分发完毕，总耗时: %v", time.Since(startTime))
}
