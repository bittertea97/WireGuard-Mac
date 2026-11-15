# WireGuard CLI (experimental)

This repository provides an experimental, terminal-first, cross-platform CLI prototype. The active implementation is in [`cli/`](cli/); the original Apple application sources remain in the repository but are not an active product target. No GUI is planned.

The CLI shares command handling and configuration checks, then selects an OS-specific backend at build time. Linux and macOS currently use `wg-quick` and `wg`; Windows uses the official WireGuard tunnel service tools. The cross-platform adapters compile, but runtime testing is still required on Linux and Windows. See [README.zh-CN.md](README.zh-CN.md) and the [cross-platform CLI guide](docs/WireGuard-CLI-跨平台指南.zh-CN.md).

## Build and test

Install Go 1.22 or newer, then:

```sh
cd cli
go test ./...
go build -o wgctl .
```

On Windows, build with `go build -o wgctl.exe .`.

On Linux and macOS, install the platform's WireGuard tools and run tunnel-changing commands with root privileges. On Windows, install the official WireGuard for Windows package and use an elevated PowerShell session.

The prior macOS shell prototype remains at [`scripts/wireguard-cli`](scripts/wireguard-cli) as a tested fallback during migration. It is macOS-only and is not the cross-platform entry point.

Upstream source: <https://git.zx2c4.com/wireguard-apple>

## MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
