package main

import (
	"database/sql"
)

/*
CREATE TABLE public.clients (
	id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
	firstname varchar NOT NULL,
	lastname varchar NOT NULL,
	username varchar NOT NULL UNIQUE,
	active boolean NOT NULL DEFAULT True
);
*/

func (p *Client) CreateClient(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO public.clients(firstname, lastname, username) VALUES($1, $2, $3) RETURNING id",
		&p.Firstname, &p.Lastname, &p.Username,
	).Scan(
		&p.Id,
	)
}

func (p *Client) GetClient(db *sql.DB) error {
	return db.QueryRow(
		"SELECT * FROM public.clients WHERE id=$1",
		&p.Id,
	).Scan(
		&p.Id, &p.Firstname, &p.Lastname, &p.Username, &p.Active,
	)
}
