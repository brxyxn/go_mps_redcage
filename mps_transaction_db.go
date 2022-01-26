package main

import (
	"database/sql"
	"errors"
	"log"
)

var TransactionQuery = `CREATE TABLE IF NOT EXISTS public."transaction" (
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
	amount numeric(13,2) NOT NULL,
	transaction_type smallint NOT NULL,
	description varchar NOT NULL,
	receiver_account_id integer NOT NULL,
	sender_account_id integer,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT transaction_pk PRIMARY KEY (id),
	CONSTRAINT transaction_fk FOREIGN KEY (receiver_account_id) REFERENCES public.account(id) ON DELETE RESTRICT
);`

/*
This function returns a list of all transactions owned by an {account_id}
*/
func (p *Transaction) GetTransactions(db *sql.DB) (Transactions, error) {
	rows, err := db.Query(
		"SELECT * FROM public.transaction WHERE receiver_account_id=$1",
		&p.ReceiverAccountId,
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
This function validates the type of transaction that is requested
and perform the addition or substaction needed into the account's
table and creates the transaction registry.
*/
func (p *Transaction) CreateTransaction(db *sql.DB) error {
	// Validate what type of transaction is being performed and process the request.
	switch p.TransactionType {
	case DictTransactionType.Deposit:
		// When a deposit is processed the account balance is updated
		oldBalance, _, err := p.getBalanceAccount(db, p.ReceiverAccountId)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// Deposit represents an addition
		newBalance := oldBalance + p.Amount

		// Update Account
		p.updateAccount(db, newBalance)
		p.registerTransaction(db)
		return nil

	case DictTransactionType.Withdraw:
		// When a withdraw is processed the account balance is updated
		oldBalance, _, err := p.getBalanceAccount(db, p.ReceiverAccountId)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// Deposit represents an substraction
		newBalance := oldBalance - p.Amount

		// Updating Account
		p.updateAccount(db, newBalance)
		err = p.registerTransaction(db)
		return err

	case DictTransactionType.Transfer:
		oldBalanceSenderAccount, currencySender, err := p.getBalanceAccount(db, p.SenderAccountId)
		if err != nil {
			log.Println(err)
			return err
		}
		oldBalanceReciverAccount, currencyReceiver, err := p.getBalanceAccount(db, p.ReceiverAccountId)
		if err != nil {
			log.Println(err)
			return err
		}
		if currencySender != currencyReceiver {
			log.Println("Currencies don't match.")
			return errors.New("currencies don't match")
		}

		newBalanceSenderAccount := oldBalanceSenderAccount - p.Amount
		newBalanceReceiverAccount := oldBalanceReciverAccount + p.Amount

		p.updateAccountsOnTransfer(db, newBalanceSenderAccount, newBalanceReceiverAccount)
		p.registerTransaction(db)
	}

	return nil
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
func (p *Transaction) registerTransaction(db *sql.DB) error {
	log.Println("Registering transaction...")
	return db.QueryRow(
		"INSERT INTO public.transaction(amount, transaction_type, description, receiver_account_id, sender_account_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		&p.Amount, &p.TransactionType, &p.Description, &p.ReceiverAccountId, &p.SenderAccountId,
	).Scan(&p.Id)
}

/*
With this function you can request the OLD_BALANCE stored in the account.

The balance will be calculated based on two params: ACCOUNT_BALANCE & CURRENCY.
These params are returned to the main function to be validated.
*/
func (p *Transaction) getBalanceAccount(db *sql.DB, accountId uint64) (accountBalance float64, currency string, err error) {
	err = db.QueryRow(
		"SELECT balance, currency FROM public.account WHERE id=$1",
		&accountId,
	).Scan(&accountBalance, &currency)

	return accountBalance, currency, err
}

/*
Used this function to update account's amount when a DEPOSIT or WITHDRAW event
is processed. New balance is calculated before the amount field is updated.

This functions panics if there is any error while processing the update.
*/
func (p *Transaction) updateAccount(db *sql.DB, newBalance float64) {
	if _, err := db.Exec(
		"UPDATE public.account SET balance=$1 WHERE id=$2;",
		newBalance, &p.ReceiverAccountId,
	); err != nil {
		log.Fatal(err)
	}
}

/*
Used this function to update accounts' amount when a TRANSFER (EFT) event
is processed. New balance is calculated before the amount field is updated.
*/
func (p *Transaction) updateAccountsOnTransfer(db *sql.DB, balanceSender, balanceReceiver float64) error {
	_, err := db.Exec(
		"UPDATE public.account SET balance=$1 WHERE id=$2",
		&balanceSender, &p.SenderAccountId,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE public.account SET balance=$1 WHERE id=$2",
		&balanceReceiver, &p.ReceiverAccountId,
	)
	if err != nil {
		return err
	}

	return nil
}
