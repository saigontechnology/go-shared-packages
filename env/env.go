package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/datngo2sgtech/go-packages/must"
)

const (
	EnvironmentVariable   = "ENV"
	dotEnvFileDirVariable = "DOT_ENV_FILE_DIR"
	EnvironmentTest       = "test"
	EnvironmentDev        = "dev"
)

// Load load environment variables from .env file by environment.
func Load() {
	envFileDir := os.Getenv(dotEnvFileDirVariable)
	envFile := filepath.Join(envFileDir, ".env."+Environment())
	if _, err := os.Stat(envFile); err == nil {
		fmt.Printf("\n[ENV File] Loading file %s \n", envFile)
		err = godotenv.Load(envFile)
		must.NotFail(err)
	}
}

func Environment() string {
	env := os.Getenv(EnvironmentVariable)
	if env == "" {
		return EnvironmentDev
	}
	return env
}

func IsTestEnv() bool {
	return Environment() == EnvironmentTest
}
