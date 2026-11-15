//go:build windows

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type windowsBackend struct{}

var validTunnelName = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

func newBackend() backend { return windowsBackend{} }

func (windowsBackend) Up(ctx context.Context, configPath string) error {
	wireguard, err := findWireGuardExe()
	if err != nil {
		return err
	}
	return runInheritedPath(ctx, wireguard, "/installtunnelservice", configPath)
}

func (windowsBackend) Down(ctx context.Context, configPath string) error {
	wireguard, err := findWireGuardExe()
	if err != nil {
		return err
	}
	name, err := tunnelName(configPath)
	if err != nil {
		return err
	}
	return runInheritedPath(ctx, wireguard, "/uninstalltunnelservice", name)
}

func (windowsBackend) Status(ctx context.Context, configPath string, stdout io.Writer) error {
	name, err := tunnelName(configPath)
	if err != nil {
		return err
	}
	wg, err := findWGExe()
	if err != nil {
		return err
	}
	command := exec.CommandContext(ctx, wg, "show", name)
	command.Stdout = stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func tunnelName(configPath string) (string, error) {
	name := strings.TrimSuffix(filepath.Base(configPath), filepath.Ext(configPath))
	if !validTunnelName.MatchString(name) {
		return "", fmt.Errorf("Windows tunnel name must contain only letters, digits, dot, underscore, or hyphen")
	}
	return name, nil
}

func findWireGuardExe() (string, error) {
	return findWindowsTool("wireguard.exe")
}

func findWGExe() (string, error) {
	if path, err := exec.LookPath("wg.exe"); err == nil {
		return path, nil
	}
	wireguard, err := findWireGuardExe()
	if err != nil {
		return "", err
	}
	candidate := filepath.Join(filepath.Dir(wireguard), "wg.exe")
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}
	return "", fmt.Errorf("wg.exe was not found beside wireguard.exe or in PATH")
}

func findWindowsTool(name string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		return path, nil
	}
	for _, variable := range []string{"ProgramFiles", "ProgramFiles(x86)"} {
		if root := os.Getenv(variable); root != "" {
			candidate := filepath.Join(root, "WireGuard", name)
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		}
	}
	return "", fmt.Errorf("%s was not found; install the official WireGuard for Windows package", name)
}

func runInheritedPath(ctx context.Context, path string, args ...string) error {
	command := exec.CommandContext(ctx, path, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
