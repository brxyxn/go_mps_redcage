package data

import (
	"database/sql"
)

/*
	Client Structures
*/

/*
Client structure for the API
swagger:model
*/
type Client struct {
	// The id of the client
	// required: false
	Id uint64 `json:"clientId"`
	// The firstname of the client
	// required: true
	Firstname string `json:"firstName"`
	// The lastname of the client
	// required: true
	Lastname string `json:"lastName"`
	// The username of the client
	// required: true
	Username string `json:"username"`
	// The status of the client, instead of deletion you change the status of the client
	// required: false
	Active bool `json:"active,omitempty"`
	// Timestamp automatically included when the client is created
	// required: false
	CreatedAt string `json:"createdAt,omitempty"`
}

/*
This function inserts the information into the table in the database.
*/
func CreateClient(db *sql.DB, p *Client) error {
	return db.QueryRow(
		"INSERT INTO public.clients(firstname, lastname, username) VALUES($1, $2, $3) RETURNING id",
		&p.Firstname, &p.Lastname, &p.Username,
	).Scan(
		&p.Id,
	)
}

/*
This function returns a client's information.
*/
func GetClient(db *sql.DB, p *Client) error {
	return db.QueryRow(
		"SELECT * FROM public.clients WHERE id=$1",
		&p.Id,
	).Scan(
		&p.Id, &p.Firstname, &p.Lastname, &p.Username, &p.Active, &p.CreatedAt,
	)
}
