package env

import (
	"os"

	"github.com/joho/godotenv"
)

func DotEnv(key string) string {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	return os.Getenv(key)
}
