package env

import (
	"os"

	"github.com/joho/godotenv"
)

func DotEnv(key string) string {
	godotenv.Load()
	return os.Getenv(key)
}
