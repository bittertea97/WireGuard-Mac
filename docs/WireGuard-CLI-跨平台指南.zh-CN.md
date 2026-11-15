# WireGuard CLI 跨平台使用指南

当前 CLI 是终端工具，不包含 GUI。公共命令和配置结构检查使用 Go 实现；隧道生命周期由操作系统适配层委托给 WireGuard 官方工具。

## 平台状态

| 平台 | 后端 | 当前验证状态 |
| --- | --- | --- |
| macOS | `wg-quick`、`wg` | Go CLI 已构建；旧 shell 原型曾完成实际隧道和 TCP 连通测试。Go CLI 的特权隧道操作仍需单独验证。 |
| Linux | `wg-quick`、`wg` | 当前在 macOS 上通过交叉编译；尚未在 Linux 主机验证。 |
| Windows | `wireguard.exe` 隧道服务、`wg.exe` | 当前在 macOS 上通过交叉编译；尚未在 Windows 主机验证。 |

交叉编译通过只说明代码能为目标平台构建，不等同于隧道运行、路由或 DNS 行为已验证。

## 准备依赖

安装 Go 1.22 或更新版本。

- macOS：安装 Homebrew 工具 `wireguard-tools` 和 `wireguard-go`。
- Linux：安装发行版提供的 WireGuard 工具包，确保 `wg`、`wg-quick` 在 `PATH` 中。
- Windows：安装官方 WireGuard for Windows，并确保 `wireguard.exe`、`wg.exe` 可被 CLI 定位。

配置文件应由 VPN 管理员提供，保存在仓库之外。它通常含有私钥，不要提交、公开或粘贴到聊天中。

## 构建与检查

在仓库根目录执行：

```sh
cd cli
go test ./...
go build -o wgctl .
```

在 Windows PowerShell 中改用 `.exe` 文件名构建：

```powershell
go build -o wgctl.exe .
```

macOS/Linux 检查配置文件的基本结构：

```sh
./wgctl check "$HOME/Downloads/your-tunnel.conf"
```

Windows PowerShell 使用：

```powershell
.\wgctl.exe check "$HOME\Downloads\your-tunnel.conf"
```

结构检查只确认 Interface、PrivateKey、Peer 和 PublicKey 等必需结构，不会建立连接，也不等同于完整的 WireGuard 配置语义校验。

## 管理隧道

macOS/Linux 使用 root 权限：

```sh
sudo ./wgctl up "$HOME/Downloads/your-tunnel.conf"
sudo ./wgctl status "$HOME/Downloads/your-tunnel.conf"
sudo ./wgctl down "$HOME/Downloads/your-tunnel.conf"
```

Windows 在“以管理员身份运行”的 PowerShell 中执行：

```powershell
.\wgctl.exe up "$HOME\Downloads\your-tunnel.conf"
.\wgctl.exe status "$HOME\Downloads\your-tunnel.conf"
.\wgctl.exe down "$HOME\Downloads\your-tunnel.conf"
```

启动和关闭会修改系统网络状态。`status` 可查看接口状态、握手和传输计数；避免把包含公钥或 Endpoint 的完整状态输出公开分享。

## 测试内网服务

将示例地址和端口替换为 `AllowedIPs` 范围内、由 VPN 管理员确认可达的目标：

macOS 查看选路接口：

```sh
route -n get 192.168.10.20 | grep interface
```

Linux 查看选路：

```sh
ip route get 192.168.10.20
```

Windows PowerShell 查看选路：

```powershell
route print 192.168.10.20
```

macOS/Linux 测试 TCP 端口：

```sh
nc -vz -w 5 192.168.10.20 8080
```

Windows PowerShell 测试 TCP 端口：

```powershell
Test-NetConnection -ComputerName 192.168.10.20 -Port 8080
```

`nc` 示例检查 TCP 连接，不验证应用层登录或服务逻辑。若服务是 HTTP，可进一步运行：

```sh
curl -v http://192.168.10.20:8080/
```

## 交叉编译

从 `cli/` 目录可以分别构建目标二进制：

```sh
GOOS=darwin GOARCH=arm64 go build -o wgctl-darwin-arm64 .
GOOS=linux GOARCH=amd64 go build -o wgctl-linux-amd64 .
GOOS=windows GOARCH=amd64 go build -o wgctl-windows-amd64.exe .
```

这些构建不会把 WireGuard 驱动或系统依赖打包进去；目标机器仍需安装对应平台的 WireGuard 工具，并使用相应管理员权限。
