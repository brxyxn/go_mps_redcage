package main

import (
	"flag"
	"log"
	"os"

	u "github.com/brxyxn/go_mps_redcage/utils"
)

var bindAddress = flag.String("BIND_ADDRESS", ":3000", "Bind address for the server")

func main() {
	flag.Parse()

	a := App{}

	a.l = log.New(os.Stdout, "go-mps-api ", log.LstdFlags)

	a.bindAddr = *bindAddress

	a.Initialize(
		u.DotEnvGet("DB_HOST"),
		u.DotEnvGet("DB_PORT"),
		u.DotEnvGet("DB_USER"),
		u.DotEnvGet("DB_PASSWORD"),
		u.DotEnvGet("DB_NAME"),
	)

	a.Run()
}
