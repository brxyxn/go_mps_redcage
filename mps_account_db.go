package main

import (
	"database/sql"
	"errors"
)

var accountQuery = `CREATE TABLE IF NOT EXISTS public.account (
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	balance numeric(13,2) NOT NULL,
	currency varchar NOT NULL,
	account_type varchar NOT NULL,
	active boolean NOT NULL DEFAULT True,
	client_id integer NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	CONSTRAINT account_pk PRIMARY KEY (id),
	CONSTRAINT account_fk FOREIGN KEY (client_id) REFERENCES public.client(id) ON DELETE CASCADE
);`

/*
This function inserts the information into the table in the database.
*/
func (p *Account) CreateAccount(db *sql.DB) error {
	if p.Currency != "USD" && p.Currency != "MXN" && p.Currency != "COP" {
		return errors.New("the currency field is not valid; USD, MXN, or COP are valid only")
	}
	return db.QueryRow(
		"INSERT INTO public.account(balance, currency, account_type, client_id) VALUES($1, $2, $3, $4) RETURNING id",
		&p.Balance, &p.Currency, &p.AccountType, &p.ClientId,
	).Scan(
		&p.Id,
	)
}

/*
This function returns the information of an account.
*/
func (p *Account) GetAccount(db *sql.DB) error {
	return db.QueryRow(
		"SELECT * FROM public.account WHERE id=$1",
		&p.Id,
	).Scan(
		&p.Id, &p.Balance, &p.Currency, &p.AccountType, &p.Active, &p.ClientId, &p.CreatedAt,
	)
}

/*
This function returns the information of all accounts owned by a {client_id}.
*/
func (p *Account) GetAccounts(db *sql.DB) (Accounts, error) {
	rows, err := db.Query(
		"SELECT * FROM public.account WHERE client_id=$1",
		&p.ClientId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := Accounts{}

	for rows.Next() {
		var row Account

		if err := rows.Scan(
			&row.Id,
			&row.Balance,
			&row.Currency,
			&row.AccountType,
			&row.Active,
			&row.ClientId,
			&row.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, row)
	}

	return accounts, nil
}

/*
	rows, err := db.Query("SELECT * FROM public.comments WHERE restaurant_id=$1 LIMIT $2 OFFSET $3",
		id, count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.Id, &c.Body, &c.Rate, &c.Rate_avg, &c.Restaurant_id); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
*/

/*
type Account struct {
    Id          uint64
    Balance     float64
    Currency    string
    AccountType uint16
    Active      bool
    ClientId    uint64
}

CREATE TABLE IF NOT EXISTS public.account (
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	balance varchar NOT NULL,
	currency varchar NOT NULL,
	account_type varchar NOT NULL UNIQUE,
	active boolean NOT NULL DEFAULT True,
	client_id integer NOT NULL,

	CONSTRAIN fk_client
		FOREIGN KEY(client_id)
			REFERENCES client(id)
);
*/
