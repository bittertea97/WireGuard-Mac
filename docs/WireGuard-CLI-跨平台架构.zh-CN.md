# WireGuard CLI 跨平台架构

## 分层

- **公共命令层**：`cli/main.go` 和 `cli/cli.go` 处理 `check`、`up`、`status`、`down`，统一参数、错误和退出码。
- **配置结构层**：`cli/config.go` 读取配置并校验必需 section/key；不会打印配置或私钥。
- **平台后端接口**：公共层只依赖 `Up`、`Down`、`Status` 操作。
- **Unix 后端**：`cli/backend_unix.go` 通过 `wg-quick`、`wg` 操作 Linux/macOS，并要求 root 权限。
- **Windows 后端**：`cli/backend_windows.go` 通过官方 `wireguard.exe` 管理 tunnel service，并用 `wg.exe` 查询状态；需要管理员 PowerShell。

各后端使用 Go build tags 在编译时选定。CLI 不实现 WireGuard 协议、加密或 TUN 驱动，数据面仍由操作系统和 WireGuard 官方组件负责。

## 兼容性边界

- 命令名称和配置路径在各平台保持一致。
- 权限提升、接口创建、路由/DNS 管理、隧道服务生命周期由各平台后端处理，因此系统行为不会完全相同。
- `check` 是跨平台的结构检查，不验证对端可达性，也不替代每个平台的完整配置解析器。
- Windows 隧道名从 `.conf` 文件名生成；当前实现仅接受 ASCII 字母、数字、点、下划线和连字符。
- macOS/Linux 的 `status` 通过配置私钥派生本机公钥并匹配实际接口。私钥只进入进程内存和 `wg pubkey` 标准输入，不放在命令参数或日志中。

## 验证矩阵

当前已完成公共层单元测试，以及 darwin/arm64、linux/amd64、windows/amd64 交叉编译。旧 macOS shell 原型完成过实际隧道连通测试；新 Go 后端的特权操作、Linux/Windows 网络行为仍待对应主机验证。
