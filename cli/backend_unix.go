//go:build darwin || linux

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type unixBackend struct{}

func newBackend() backend { return unixBackend{} }

func (unixBackend) Up(ctx context.Context, configPath string) error {
	if err := requireRoot(); err != nil {
		return err
	}
	return runInherited(ctx, "wg-quick", "up", configPath)
}

func (unixBackend) Down(ctx context.Context, configPath string) error {
	if err := requireRoot(); err != nil {
		return err
	}
	return runInherited(ctx, "wg-quick", "down", configPath)
}

func (unixBackend) Status(ctx context.Context, configPath string, stdout io.Writer) error {
	if err := requireRoot(); err != nil {
		return err
	}
	parsed, err := parseConfig(configPath)
	if err != nil {
		return err
	}
	defer wipe(parsed.privateKey)

	keyInput := append(append([]byte(nil), parsed.privateKey...), '\n')
	defer wipe(keyInput)
	pubKeyCommand := exec.CommandContext(ctx, "wg", "pubkey")
	pubKeyCommand.Stdin = bytes.NewReader(keyInput)
	publicKey, err := pubKeyCommand.Output()
	if err != nil {
		return fmt.Errorf("derive interface identity with wg pubkey: %w", err)
	}
	expected := strings.TrimSpace(string(publicKey))

	interfaces, err := exec.CommandContext(ctx, "wg", "show", "interfaces").Output()
	if err != nil {
		return fmt.Errorf("list WireGuard interfaces: %w", err)
	}
	for _, name := range strings.Fields(string(interfaces)) {
		actual, err := exec.CommandContext(ctx, "wg", "show", name, "public-key").Output()
		if err != nil || strings.TrimSpace(string(actual)) != expected {
			continue
		}
		fmt.Fprintf(stdout, "Tunnel is up: %s\n", name)
		command := exec.CommandContext(ctx, "wg", "show", name)
		command.Stdout = stdout
		command.Stderr = os.Stderr
		return command.Run()
	}
	return fmt.Errorf("no running interface matches this config")
}

func requireRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("run this operation as root (for example, with sudo)")
	}
	return nil
}

func runInherited(ctx context.Context, name string, args ...string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("required command %q was not found in PATH", name)
	}
	command := exec.CommandContext(ctx, path, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
