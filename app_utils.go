package main

import (
	"encoding/json"
	"log"
	"net/http"
)

/*
	HTTP Responses
*/

/*
This function responds an error message to the frontend as JSON.
*/
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
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
	Error Handling
*/

/*
Use this func to panic the application if need to exit
in case there's any unexpected output in the program.
Ref.: https://gobyexample.com/panic
*/
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

/*
Use this func to log the error only.
*/
func LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
Use this func to log the error and set a custom message.
*/
func LogErrorMsg(err error, msg string) {
	if err != nil {
		log.Fatalf("%s - %s", err, msg)
	}
}

type responder func(w http.ResponseWriter, code int, message string)

func Responder(err error, fn responder) responder {
	if err != nil {
		return fn
	}
	return nil
}
