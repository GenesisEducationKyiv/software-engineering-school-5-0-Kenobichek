package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// parseEnvFile reads environment variables from the current process and merges them with variables from the specified file.
// The file is parsed line-by-line, supporting comments and quoted values. Variables from the file override existing ones.
// Returns a map of combined environment variables and any error encountered during file reading. If the file does not exist, only the current environment is returned without error.
func parseEnvFile(envPath string) (map[string]string, error) {
	envVars := make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}

	file, err := os.Open(envPath) //nolint:gosec
	if err != nil {
		if os.IsNotExist(err) {
			return envVars, nil
		}
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		envVars[key] = value
	}

	return envVars, scanner.Err()
}
