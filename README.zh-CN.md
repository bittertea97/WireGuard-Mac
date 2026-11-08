# WireGuard-Mac

这是 WireGuard 官方 Apple 客户端的非官方 macOS 分支，基于上游 MIT 许可源码。当前版本保留上游功能和简体中文界面，尚未更改隧道实现或官方品牌素材。

## 构建 macOS 版本

需要安装完整 Xcode、Go 和 SwiftLint。仅安装 Xcode Command Line Tools 不够。

```sh
brew install go swiftlint
cp Sources/WireGuardApp/Config/Developer.xcconfig.template \
  Sources/WireGuardApp/Config/Developer.xcconfig
```

编辑 `Sources/WireGuardApp/Config/Developer.xcconfig`，填写自己的 Apple Developer Team ID，以及已配置 Network Extensions 能力的 macOS App ID。此文件已加入 `.gitignore`，不要提交个人签名配置。

然后用 Xcode 打开 `WireGuard.xcodeproj`，选择 macOS 的 `WireGuard` scheme 和目标 Mac，完成签名并构建。Go bridge 会在构建过程中按项目配置编译。

首次安装和启动隧道需要有效的 Network Extension entitlement 与签名。未配置个人 Team ID 和相应能力时，项目可以检出，但无法得到可用的独立 VPN 客户端。

## 上游与许可

- 上游项目：<https://git.zx2c4.com/wireguard-apple>
- 上游构建说明：<https://git.zx2c4.com/wireguard-apple/tree/README.md>
- 许可：MIT，详见 [COPYING](COPYING)

本分支不是 WireGuard LLC 发布或支持的官方发行版。WireGuard 名称和相关商标仍归其权利人所有。
