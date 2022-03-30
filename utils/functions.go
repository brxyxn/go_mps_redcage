package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func DotEnvGet(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file.", err)
	}

	return os.Getenv(key)
}

/*
This function responds the information successfully as JSON.
*/
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

/*
This function responds an error message to the frontend as JSON.
*/
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

/*
Handling logs
This functions are suitable for logging into the console using some
prefix to identify better the nature of the errors and logs.

This could be improved and initialized in a better way.
*/
type logging struct {
	Print *log.Logger
}

type any = interface{}

var l logging

func InitLogs(prefix string) *logging {
	l.Print = log.New(os.Stdout, prefix, log.LstdFlags)
	return &l
}

func LogInfo(x ...interface{}) {
	s := "[INFO]"
	args := append([]interface{}{s}, x...)
	l.Print.Println(args...)
}

func LogError(x ...interface{}) {
	s := "[ERROR]"
	args := append([]interface{}{s}, x...)
	l.Print.Println(args...)
}

func LogDebug(x ...interface{}) {
	s := "[DEBUG]"
	args := append([]interface{}{s}, x...)
	l.Print.Println(args...)
}
