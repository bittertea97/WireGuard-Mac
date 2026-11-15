package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type backend interface {
	Up(context.Context, string) error
	Down(context.Context, string) error
	Status(context.Context, string, io.Writer) error
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer, system backend) int {
	if len(args) != 2 {
		writeUsage(stderr)
		return 2
	}

	action := args[0]
	if action != "check" && action != "up" && action != "down" && action != "status" {
		writeUsage(stderr)
		return 2
	}

	configPath, err := filepath.Abs(args[1])
	if err != nil {
		return reportError(stderr, "resolve config path", err)
	}
	info, err := os.Stat(configPath)
	if err != nil || !info.Mode().IsRegular() {
		return reportError(stderr, "config file is not readable", err)
	}
	if !strings.EqualFold(filepath.Ext(configPath), ".conf") {
		return reportError(stderr, "config file must end in .conf", nil)
	}

	parsed, err := parseConfig(configPath)
	if err != nil {
		return reportError(stderr, "invalid WireGuard config structure", err)
	}
	wipe(parsed.privateKey)

	if action == "check" {
		fmt.Fprintf(stdout, "Configuration structure is valid: %s\n", filepath.Base(configPath))
		return 0
	}
	if system == nil {
		return reportError(stderr, "platform backend is unavailable", nil)
	}

	switch action {
	case "up":
		err = system.Up(ctx, configPath)
	case "down":
		err = system.Down(ctx, configPath)
	case "status":
		err = system.Status(ctx, configPath, stdout)
	}
	if err != nil {
		return reportError(stderr, action+" failed", err)
	}
	return 0
}

func reportError(stderr io.Writer, message string, err error) int {
	if err != nil {
		fmt.Fprintf(stderr, "wgctl: %s: %v\n", message, err)
	} else {
		fmt.Fprintf(stderr, "wgctl: %s\n", message)
	}
	return 1
}

func writeUsage(stderr io.Writer) {
	fmt.Fprintln(stderr, "Usage: wgctl <check|up|status|down> <config.conf>")
}

func wipe(value []byte) {
	for i := range value {
		value[i] = 0
	}
}
