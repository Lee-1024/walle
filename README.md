# Walle (华力) - 多主机并发 SSH 执行与文件分发工具

Walle 是一款使用 Go 语言编写的高性能、轻量级分布式 Linux 批量运维工具。它允许你通过简单的配置文件管理数百台 Linux 主机，并通过现代化的命令行交互，并发在所有机器上执行相同的指令或分发相同的文件。

## ✨ 特性

- **⚡ 极致并发**：利用 Go 原生 Goroutines，瞬间处理海量服务器，拒绝单线程串行阻塞。
- **📝 YAML 资产管理**：支持清晰易读的集群配置文件，支持全局默认配置与节点级特定配置覆盖。
- **🔐 智能认证**：支持 **密码认证** 与 **SSH 密钥（Private Key）认证** 自动切换。
- **🛠️ 现代化 CLI**：基于业内标准的 `Cobra` 构建，提供直观的命令别名及健壮的参数校验。
- **🎨 炫彩终端**：集成彩色高亮输出日志，按主机别名（Alias）进行隔离渲染，执行状态一目了然。

## 📦 快速开始

### 1. 本地环境准备
确保你的本地开发环境已安装 Go 1.18+。

```bash
mkdir -p walle-ssh && cd walle-ssh

# 1. 批量在全量主机上运行指令
go run . run "uname -a && df -h /"

# 2. 批量分发本地文件至全量主机指定目录
go run . push "./my_nginx.conf" "/etc/nginx/nginx.conf"

# 3. 如果配置文件不在当前目录或改了名字，使用 -c 参数指定
go run . -c /path/to/my_hosts.yaml run "systemctl restart nginx"

# 编译出可直接丢进系统 PATH 的独立可执行文件
go build -o walle .

# 投入实用
./walle run "uptime"