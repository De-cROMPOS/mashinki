package envhandler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getEnvPath() string {
	// Getting path to current directory
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}

	// Looking for the .env file in current directory
	envPath := filepath.Join(workDir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		return envPath
	}

	// If not found, looking for the .env file in parent directory
	parentDir := filepath.Dir(workDir)
	envPath = filepath.Join(parentDir, ".env")
	if _, err := os.Stat(envPath); err == nil {
		return envPath
	}

	fmt.Println("Error getting working directory:", err)
	os.Exit(1)
	return ""
}

func GetEnv(key string) string {
	envPath := getEnvPath()

	envData, err := os.ReadFile(envPath)
	if err != nil {
		log.Fatalf("Error while opening .env file: %s", err)
	}

	for _, line := range strings.Split(string(envData), "\n") {
		if strings.HasPrefix(line, key+"=") {
			return strings.TrimPrefix(line, key+"=")
		}
	}
	return ""
}
