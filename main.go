package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	StartServer()
}

func StartServer() {
	http.HandleFunc("/", landingPageHandler)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Starting server... Requested: %s\n", r.URL.Path)
}
