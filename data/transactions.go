package data

import (
	"database/sql"
	"log"

	"github.com/shopspring/decimal"
)

/*
	Transaction Structures
*/

type Transactions []Transaction

/*
Transaction structure for the API
swagger:model
*/
type Transaction struct {
	// The Id of the transaction
	// required: false
	Id uint64 `json:"transactionId"`
	// The transaction amount
	// required: true
	Amount string `json:"amount"`
	// The transaction type - withdraw, deposit, transfer
	// required: true
	TransactionType TransactionType `json:"transactionType"`
	// The description of the transaction
	// required: true
	Description string `json:"description"`
	// The account number/id of the receiver or affected account
	// required: true
	ReceiverAccountId uint64 `json:"receiverAccountId"`
	// The account number/id of the sender account if it's a transfer of funds
	// required: true
	SenderAccountId uint64 `json:"senderAccountId"`
	// Timestamp automatically included when the transaction is created
	// required: false
	CreatedAt string `json:"createdAt,omitempty"`
}

type TransactionType int

// Dictionary of TransactionType
var DictTransactionType = struct {
	Deposit  TransactionType
	Withdraw TransactionType
	Transfer TransactionType
}{
	Deposit:  1,
	Withdraw: 2,
	Transfer: 3,
}

/*
This function returns a list of all transactions owned by an {account_id}
*/
func GetTransactions(db *sql.DB, accountId int) (Transactions, error) {
	rows, err := db.Query(
		"SELECT * FROM public.transactions WHERE receiver_account_id=$1",
		&accountId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transactions := Transactions{}

	for rows.Next() {
		var row Transaction

		if err := rows.Scan(
			&row.Id,
			&row.Amount,
			&row.TransactionType,
			&row.Description,
			&row.SenderAccountId,
			&row.ReceiverAccountId,
			&row.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, row)
	}

	return transactions, nil
}

/*
-----------------------------------------------------------------------
Splitted responsibility making this code look clean and understandable.
-----------------------------------------------------------------------
*/

/*
Used this function to register/insert the new registry when a
DEPOSIT, WITHDRAW or TRANSFER transaction is processed.
*/
func RegisterTransaction(db *sql.DB, p *Transaction) error {
	return db.QueryRow(
		"INSERT INTO public.transactions(amount, transaction_type, description, receiver_account_id, sender_account_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		&p.Amount, &p.TransactionType, &p.Description, &p.ReceiverAccountId, &p.SenderAccountId,
	).Scan(&p.Id)
}

/*
With this function you can request the OLD_BALANCE stored in the account.

The balance will be calculated based on two params: ACCOUNT_BALANCE & CURRENCY.
These params are returned to the main function to be validated.
*/
func GetBalanceAccount(db *sql.DB, accountId uint64) (accountBalance string, currency string, err error) {
	err = db.QueryRow(
		"SELECT balance, currency FROM public.accounts WHERE id=$1",
		&accountId,
	).Scan(&accountBalance, &currency)

	return accountBalance, currency, err
}

/*
Used this function to update account's amount when a DEPOSIT or WITHDRAW event
is processed. New balance is calculated before the amount field is updated.

This functions panics if there is any error while processing the update.
*/
func UpdateAccount(db *sql.DB, p *Transaction, newBalance decimal.Decimal) {
	if _, err := db.Exec(
		"UPDATE public.accounts SET balance=$1 WHERE id=$2;",
		newBalance, &p.ReceiverAccountId,
	); err != nil {
		log.Fatal(err)
	}
}

/*
Used this function to update accounts' amount when a TRANSFER (EFT) event
is processed. New balance is calculated before the amount field is updated.
*/
func UpdateAccountsOnTransfer(db *sql.DB, p *Transaction, balanceSender, balanceReceiver decimal.Decimal) error {
	_, err := db.Exec(
		"UPDATE public.accounts SET balance=$1 WHERE id=$2",
		&balanceSender, &p.SenderAccountId,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE public.accounts SET balance=$1 WHERE id=$2",
		&balanceReceiver, &p.ReceiverAccountId,
	)
	if err != nil {
		return err
	}

	return nil
}
