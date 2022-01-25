package main

import (
	"database/sql"
	"log"
)

var clientQuery = `CREATE TABLE IF NOT EXISTS public.client (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	firstname varchar NOT NULL,
	lastname varchar NOT NULL,
	username varchar NOT NULL UNIQUE,
	active boolean NOT NULL DEFAULT True,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT client_pk PRIMARY KEY (id),
	CONSTRAINT client_username_key UNIQUE (username)
);`

/*
This function inserts the information into the table in the database.
*/
func (p *Client) CreateClient(db *sql.DB) error {
	log.Println(p)
	return db.QueryRow(
		"INSERT INTO public.client(firstname, lastname, username) VALUES($1, $2, $3) RETURNING id",
		&p.Firstname, &p.Lastname, &p.Username,
	).Scan(
		&p.Id,
	)
}

/*
This function returns a client's information.
*/
func (p *Client) GetClient(db *sql.DB) error {
	return db.QueryRow(
		"SELECT * FROM public.client WHERE id=$1",
		&p.Id,
	).Scan(
		&p.Id, &p.Firstname, &p.Lastname, &p.Username, &p.Active, &p.CreatedAt,
	)
}
