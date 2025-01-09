package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnv will load the environment variable from the .env file.
func LoadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		split := strings.Split(line, "=")
		if len(split) != 2 {
			return fmt.Errorf("invalid env line: %s", line)
		}

		if err := os.Setenv(split[0], split[1]); err != nil {
			return err
		}

	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
