package models

import (
	"database/sql"

	model "github.com/brxyxn/go_mps_redcage/structs"
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

type ClientRef struct {
	*model.Client
}

func (p *ClientRef) CreateClient(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO public.clients(firstname, lastname, username) VALUES($1, $2, $3) RETURNING id",
		&p.Client.Firstname, &p.Client.Lastname, &p.Client.Username,
	).Scan(
		&p.Client.Id,
	)
}

func (p *ClientRef) GetClient(db *sql.DB) error {
	return db.QueryRow(
		"SELECT * FROM public.clients WHERE id=$1",
		&p.Client.Id,
	).Scan(
		&p.Client.Id, &p.Client.Firstname, &p.Client.Lastname, &p.Client.Username, &p.Client.Active,
	)
}
