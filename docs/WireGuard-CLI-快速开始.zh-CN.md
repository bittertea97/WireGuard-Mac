# WireGuard macOS CLI 快速开始

本指南介绍如何使用仓库中的 `scripts/wireguard-cli` 在 macOS 上检查配置、启动 WireGuard 隧道、测试内网服务并关闭隧道。

这是一个本机 CLI 原型，不是图形化应用或正式发行包。它调用 Homebrew 提供的 `wg-quick`、`wireguard-go` 和 `wg`。运行它不需要 Apple Developer Program；启动和关闭隧道需要本机管理员密码。

## 安装依赖

安装 Homebrew 后运行：

```sh
brew install wireguard-tools wireguard-go
```

确认工具可用：

```sh
command -v wg
command -v wg-quick
command -v wireguard-go
```

## 准备配置

向 VPN 服务管理员获取 WireGuard `.conf` 配置文件，并将其保存在仓库目录之外，例如 `~/Downloads/your-tunnel.conf`。限制本机文件权限：

```sh
chmod 600 "$HOME/Downloads/your-tunnel.conf"
```

配置通常含有私钥。不要把它粘贴到聊天、截图、公开 issue 或 Git 仓库，也不要把私钥写进命令行参数。CLI 会直接从原路径读取配置，不会复制到仓库。

## 检查配置

```sh
scripts/wireguard-cli check "$HOME/Downloads/your-tunnel.conf"
```

此命令进行本地结构检查，不会启动隧道，也不会输出配置内容。它不能代替连接测试；配置能通过检查不代表服务端可达。

## 启动和检查状态

启动隧道：

```sh
sudo scripts/wireguard-cli up "$HOME/Downloads/your-tunnel.conf"
```

在 macOS 密码提示中输入本机管理员密码。输入时终端不会显示字符。启动会创建 `utun` 接口，并可能修改路由和 DNS 设置。

查看接口、对端握手和传输计数：

```sh
sudo scripts/wireguard-cli status "$HOME/Downloads/your-tunnel.conf"
```

`latest handshake` 显示最近一次握手；`transfer` 显示收发字节数。只有接口处于 `UP` 不足以证明隧道可用。WireGuard 通常会在有流量需要发送给对端时发起握手，因此测试时要访问落在配置 `AllowedIPs` 范围内的实际目标。

## 测试内网服务

以下地址仅为示例。将 `192.168.10.20` 和 `8080` 替换为 VPN 服务管理员提供的内网主机地址和服务端口。

确认路由经过隧道接口：

```sh
route -n get 192.168.10.20 | grep interface
```

测试 TCP 端口是否接受连接：

```sh
nc -vz -w 5 192.168.10.20 8080
```

`succeeded` 表示 TCP 连接已建立；这不代表应用层登录或服务功能也正常。根据服务协议继续访问，例如 HTTP：

```sh
curl -v http://192.168.10.20:8080/
```

如果服务使用 SSH，则使用对应用户名和端口：

```sh
ssh -p 8080 username@192.168.10.20
```

上面的 `nc -vz` 是 TCP 测试。UDP 没有 TCP 式的连接握手，单次 UDP 探测通常不能确认端口是否开放；应使用该服务自身的客户端或协议级请求验证。

## 关闭隧道

```sh
sudo scripts/wireguard-cli down "$HOME/Downloads/your-tunnel.conf"
```

关闭后，`wg-quick` 会清理本次隧道创建的接口和路由，并恢复它备份的 DNS 设置。若网络状态异常，先确认关闭命令是否成功，再检查系统网络设置。

## 常见问题

- **配置检查通过但没有握手：** 对 `AllowedIPs` 中一个实际可达的内网地址发起请求；仍无握手时，联系服务管理员核对 Endpoint、双方公钥、服务端 peer 配置及防火墙。
- **端口测试超时：** 确认目标地址属于 `AllowedIPs`、路由经过 `utun` 接口，并确认目标主机及对应 TCP 服务正在运行。
- **提示找不到工具：** 确认已安装 `wireguard-tools` 和 `wireguard-go`，并检查 Homebrew 的 `bin` 目录在 `PATH` 中。
- **配置文件安全：** 不要提交 `.conf`、私钥、个人签名设置或其他凭证。推送前检查 `git status` 和暂存文件列表。

本项目是非官方分支，不由 WireGuard LLC 发布或支持。
