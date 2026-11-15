package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestParseConfigValid(t *testing.T) {
	path := writeTestConfig(t, "[Interface]\nPrivateKey = private\n\n[Peer]\nPublicKey = public\nAllowedIPs = 192.168.10.0/24\n")
	parsed, err := parseConfig(path)
	if err != nil {
		t.Fatalf("parseConfig() error = %v", err)
	}
	wipe(parsed.privateKey)
}

func TestParseConfigRejectsMissingPeerKey(t *testing.T) {
	path := writeTestConfig(t, "[Interface]\nPrivateKey = private\n\n[Peer]\nAllowedIPs = 192.168.10.0/24\n")
	if _, err := parseConfig(path); err == nil {
		t.Fatal("parseConfig() accepted a peer without a public key")
	}
}

func TestRunCheckDoesNotNeedPlatformBackend(t *testing.T) {
	path := writeTestConfig(t, "[Interface]\nPrivateKey = private\n[Peer]\nPublicKey = public\n")
	var stdout, stderr bytes.Buffer
	code := run(context.Background(), []string{"check", path}, &stdout, &stderr, nil)
	if code != 0 {
		t.Fatalf("run(check) exit code = %d; stderr: %s", code, stderr.String())
	}
}

func TestRunDispatchesTunnelActions(t *testing.T) {
	path := writeTestConfig(t, "[Interface]\nPrivateKey = private\n[Peer]\nPublicKey = public\n")
	for _, action := range []string{"up", "down", "status"} {
		t.Run(action, func(t *testing.T) {
			fake := &fakeBackend{}
			var stdout, stderr bytes.Buffer
			code := run(context.Background(), []string{action, path}, &stdout, &stderr, fake)
			if code != 0 {
				t.Fatalf("run(%s) exit code = %d; stderr: %s", action, code, stderr.String())
			}
			if fake.action != action {
				t.Fatalf("backend action = %q, want %q", fake.action, action)
			}
			if fake.configPath != path {
				t.Fatalf("backend config path = %q, want %q", fake.configPath, path)
			}
		})
	}
}

type fakeBackend struct {
	action     string
	configPath string
}

func (fake *fakeBackend) Up(_ context.Context, path string) error {
	fake.action, fake.configPath = "up", path
	return nil
}

func (fake *fakeBackend) Down(_ context.Context, path string) error {
	fake.action, fake.configPath = "down", path
	return nil
}

func (fake *fakeBackend) Status(_ context.Context, path string, _ io.Writer) error {
	fake.action, fake.configPath = "status", path
	return nil
}

func writeTestConfig(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.conf")
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}
