package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type parsedConfig struct {
	privateKey []byte
}

func parseConfig(path string) (parsedConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return parsedConfig{}, err
	}
	defer file.Close()

	var result parsedConfig
	keepKey := false
	defer func() {
		if !keepKey {
			wipe(result.privateKey)
		}
	}()
	section := ""
	interfaces, peers := 0, 0
	privateKeys, currentPeerPublicKeys := 0, 0
	lineNumber := 0
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 4096), 1024*1024)

	finishPeer := func() error {
		if section == "peer" && currentPeerPublicKeys != 1 {
			return fmt.Errorf("peer section must contain exactly one PublicKey")
		}
		return nil
	}

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if comment := strings.IndexByte(line, '#'); comment >= 0 {
			line = strings.TrimSpace(line[:comment])
		}
		if line == "" {
			continue
		}
		switch line {
		case "[Interface]":
			if err := finishPeer(); err != nil {
				return parsedConfig{}, fmt.Errorf("line %d: %w", lineNumber, err)
			}
			section = "interface"
			interfaces++
			continue
		case "[Peer]":
			if err := finishPeer(); err != nil {
				return parsedConfig{}, fmt.Errorf("line %d: %w", lineNumber, err)
			}
			section = "peer"
			peers++
			currentPeerPublicKeys = 0
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || key == "" || value == "" || section == "" {
			return parsedConfig{}, fmt.Errorf("line %d: malformed setting", lineNumber)
		}
		if section == "interface" && key == "PrivateKey" {
			privateKeys++
			wipe(result.privateKey)
			result.privateKey = []byte(value)
		}
		if section == "peer" && key == "PublicKey" {
			currentPeerPublicKeys++
		}
	}
	if err := scanner.Err(); err != nil {
		return parsedConfig{}, fmt.Errorf("read config: %w", err)
	}
	if err := finishPeer(); err != nil {
		return parsedConfig{}, err
	}
	if interfaces != 1 || peers < 1 || privateKeys != 1 {
		return parsedConfig{}, fmt.Errorf("expected one Interface with one PrivateKey and at least one Peer with a PublicKey")
	}
	keepKey = true
	return result, nil
}
