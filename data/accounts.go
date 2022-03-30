package data

import (
	"database/sql"
	"errors"
)

/*
	Account Structures
*/

type Accounts []Account

/*
Account structure for the API
swagger:model
*/
type Account struct {
	// The id of the account
	// required: false
	Id uint64 `json:"accountId"`
	// The current balance of the account
	// required: true
	Balance string `json:"balance"`
	// The currency symbol or description of the account - USD, MXN, COP
	// required: true
	Currency string `json:"currency"`

	// The account type can be checking, savings, credit card
	// required: true
	AccountType AccountType `json:"accountType"`
	// The status of the account, instead of deletion you change the status of the account
	// required: true
	Active bool `json:"active,omitempty"`
	// The ClientId is considered the FK of this model
	// required: true
	ClientId uint64 `json:"clientId"`
	// Timestamp automatically included when the transaction is created
	// required: false
	CreatedAt string `json:"createdAt,omitempty"`
}

type AccountType string

// Dictionary of AccountType
var DictAccountType = struct {
	Savings    AccountType
	Checking   AccountType
	CreditCard AccountType
}{
	Savings:    "Savings",
	Checking:   "Checking",
	CreditCard: "Credit Card",
}

/*
This function inserts the information into the table in the database.
*/
func CreateAccount(db *sql.DB, v *Account) error {
	if v.Currency != "USD" && v.Currency != "MXN" && v.Currency != "COP" {
		return errors.New("the currency field is not valid; USD, MXN, or COP are valid only")
	}
	return db.QueryRow(
		"INSERT INTO public.accounts(balance, currency, account_type, client_id) VALUES($1, $2, $3, $4) RETURNING id",
		&v.Balance, &v.Currency, &v.AccountType, &v.ClientId,
	).Scan(
		&v.Id,
	)
}

/*
This function returns the information of an account.
*/
func GetAccount(db *sql.DB, p *Account) error {
	return db.QueryRow(
		"SELECT * FROM public.accounts WHERE id=$1",
		&p.Id,
	).Scan(
		&p.Id, &p.Balance, &p.Currency, &p.AccountType, &p.Active, &p.ClientId, &p.CreatedAt,
	)
}

/*
This function returns the information of all accounts owned by a {client_id}.
*/
func GetAccounts(db *sql.DB, clientID int) Accounts {
	rows, err := db.Query(
		"SELECT * FROM public.accounts WHERE client_id=$1",
		clientID,
	)
	if err != nil {
		return nil
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
			return nil
		}
		accounts = append(accounts, row)
	}

	return accounts
}
