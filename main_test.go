package main

import (
	"log"
	"os"
	"testing"
)

var a *App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("DB_HOST_TEST"),
		os.Getenv("DB_PORT_TEST"),
		os.Getenv("DB_USER_TEST"),
		os.Getenv("DB_PASSWORD_TEST"),
		os.Getenv("DB_NAME_TEST"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.db.Exec("DELETE FROM products")
	a.db.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS public.clients (
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	firstname varchar NOT NULL,
	lastname varchar NOT NULL,
	username varchar NOT NULL UNIQUE,
	active boolean NOT NULL DEFAULT TRUE,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	CONSTRAINT client_pk PRIMARY KEY (id),
	CONSTRAINT client_username_key UNIQUE (username)
);`
