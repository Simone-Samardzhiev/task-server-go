package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnv loads environment variables from a .env file.
func LoadEnv() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %w", err)
	}

	for {
		envFile := filepath.Join(currentDir, ".env")
		if fileExists(envFile) {
			if err := readEnvFile(envFile); err != nil {
				return fmt.Errorf("error reading .env file: %w", err)
			}
			return nil
		}

		currentDir = filepath.Dir(currentDir)
		if fileExists(filepath.Join(currentDir, "go.mod")) {
			break
		}
	}

	return fmt.Errorf(".env file not found")
}

// readEnvFile reads the .env file and loads variables into the environment.
func readEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		split := strings.SplitN(line, "=", 2)
		if len(split) != 2 {
			return fmt.Errorf("invalid env format: %s", line)
		}

		key := strings.TrimSpace(split[0])
		value := strings.TrimSpace(split[1])
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting environment variable %s: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

// fileExists checks if a file exists and is not a directory.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
