# WireGuard CLI（实验性）

项目当前方向是终端优先的跨平台 CLI，不再开发 GUI。仓库保留了上游 WireGuard Apple 应用源代码作为历史基础，但这些 App 源码不是当前产品目标，也不代表 GUI 已构建或验证可用。

公共 Go CLI 位于 [`cli/`](cli/)，包含共享命令层、配置结构检查和按操作系统选择的平台后端。macOS/Linux 后端调用 `wg-quick`/`wg`；Windows 后端调用官方 WireGuard tunnel service 工具。Linux 和 Windows 目前仅完成交叉编译，仍需在对应系统上验证运行行为。

## 快速开始

需要 Go 1.22 或更新版本：

```sh
cd cli
go test ./...
go build -o wgctl .
```

Windows 请使用 `go build -o wgctl.exe .`。

macOS 和 Linux 需安装 `wg`、`wg-quick`；macOS 还需 `wireguard-go`。Windows 需安装官方 WireGuard for Windows。详细依赖、命令和状态见[跨平台 CLI 使用指南](docs/WireGuard-CLI-跨平台指南.zh-CN.md)及[架构说明](docs/WireGuard-CLI-跨平台架构.zh-CN.md)。

旧的 macOS shell 原型 [`scripts/wireguard-cli`](scripts/wireguard-cli) 保留作迁移期回退，不是跨平台入口。项目不需要 Apple Developer Program 来运行终端原型；启动和关闭隧道需要系统管理员权限。

## 上游与许可

- 上游基础：<https://git.zx2c4.com/wireguard-apple>
- 许可：MIT，详见 [COPYING](COPYING)

本项目不是 WireGuard LLC 发布或支持的官方产品。WireGuard 名称和相关商标仍归其权利人所有。
