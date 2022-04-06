package main

import (
	"flag"
	"log"
	"os"

	u "github.com/brxyxn/go_mps_redcage/utils"
)

func main() {
	flag.Parse()

	a := App{}

	u.InitLogs("go-mps-api ")
	a.l = log.New(os.Stdout, "go-mps-api ", log.LstdFlags)

	port := os.Getenv("PORT")
	if port == "" {
		a.bindAddr = ":" + "3000"
	} else {
		a.bindAddr = ":" + port
	}
	u.LogInfo(a.bindAddr)

	a.Initialize(
		u.DotEnvGet("DB_HOST"),
		u.DotEnvGet("DB_PORT"),
		u.DotEnvGet("DB_USER"),
		u.DotEnvGet("DB_PASSWORD"),
		u.DotEnvGet("DB_NAME"),
	)

	a.Run()
}
