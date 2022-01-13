package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

func main() {
	a := App{}

	a.Initialize(
		GetEnv("APP_DB_USERNAME"),
		GetEnv("APP_DB_PASSWORD"),
		GetEnv("APP_DB_NAME"),
	)

	a.Run(":3000")
}
