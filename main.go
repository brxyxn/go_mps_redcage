package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	a := App{}

	fmt.Println(dotEnvGet("TESTING"))

	a.Initialize(
		dotEnvGet("DEV_DB_USERNAME"),
		dotEnvGet("DEV_DB_PASSWORD"),
		dotEnvGet("DEV_DB_NAME"),
	)
}

func dotEnvGet(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
