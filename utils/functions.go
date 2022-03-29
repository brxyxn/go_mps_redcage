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
