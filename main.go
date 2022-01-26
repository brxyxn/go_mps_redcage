package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	a := App{}

	a.Initialize(
		dotEnvGet("DB_HOST"),
		dotEnvGet("DB_PORT"),
		dotEnvGet("DB_USER"),
		dotEnvGet("DB_PASSWORD"),
		dotEnvGet("DB_NAME"),
	)

	// We try to validate that the connection to the DB is correct, otherwise, the app
	// will restart, this is a temporary solution because the postgres image usually
	// is initialized after golang's image.
	if err := a.DB.Ping(); err != nil {
		a.DB.Close()
		log.Fatal(err)
	} else {
		// Executing SQL statements to create tables and seed DB.
		query, err := ioutil.ReadFile("docker_postgres_init.sql")
		LogErrorMsg(err, "Error while reading ./docker_postgres_init.sql file")
		if _, err := a.DB.Exec(string(query)); err != nil {
			panic(err)
		}
	}

	a.Run("8080")
}

func dotEnvGet(key string) string {
	err := godotenv.Load()
	LogErrorMsg(err, "Error loading .env file.")

	return os.Getenv(key)
}
