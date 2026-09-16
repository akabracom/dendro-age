package dsn

import (
	"fmt"
	"os"
)

func FromEnv() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)
}
