package data

import (
	"database/sql"
	"errors"
	"log"
)

/*
	Transaction Structures
*/

type Transactions []Transaction

type Transaction struct {
	Id                uint64          `json:"transactionId"`
	Amount            float64         `json:"amount"`
	TransactionType   TransactionType `json:"transactionType"`
	Description       string          `json:"description"`
	ReceiverAccountId uint64          `json:"receiverAccountId"`
	SenderAccountId   uint64          `json:"senderAccountId"`
	CreatedAt         string          `json:"createdAt,omitempty"`
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
This function validates the type of transaction that is requested
and perform the addition or substaction needed into the account's
table and creates the transaction registry.
*/
func CreateTransaction(db *sql.DB, p *Transaction) error {
	// Validate what type of transaction is being performed and process the request.
	switch p.TransactionType {
	case DictTransactionType.Deposit:
		// When a deposit is processed the account balance is updated
		oldBalance, _, err := getBalanceAccount(db, p, p.ReceiverAccountId)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// Deposit represents an addition
		newBalance := oldBalance + p.Amount

		// Update Account
		updateAccount(db, p, newBalance)
		registerTransaction(db, p)
		return nil

	case DictTransactionType.Withdraw:
		// When a withdraw is processed the account balance is updated
		oldBalance, _, err := getBalanceAccount(db, p, p.ReceiverAccountId)
		if err != nil {
			log.Fatal(err)
			return err
		}

		// Deposit represents an substraction
		newBalance := oldBalance - p.Amount

		// Updating Account
		updateAccount(db, p, newBalance)
		err = registerTransaction(db, p)
		return err

	case DictTransactionType.Transfer:
		oldBalanceSenderAccount, currencySender, err := getBalanceAccount(db, p, p.SenderAccountId)
		if err != nil {
			log.Println(err)
			return err
		}
		oldBalanceReciverAccount, currencyReceiver, err := getBalanceAccount(db, p, p.ReceiverAccountId)
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

		updateAccountsOnTransfer(db, p, newBalanceSenderAccount, newBalanceReceiverAccount)
		registerTransaction(db, p)
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
func registerTransaction(db *sql.DB, p *Transaction) error {
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
func getBalanceAccount(db *sql.DB, p *Transaction, accountId uint64) (accountBalance float64, currency string, err error) {
	err = db.QueryRow(
		"SELECT balance, currency FROM public.accounts WHERE id=$1",
		&p.ReceiverAccountId,
	).Scan(&accountBalance, &currency)

	return accountBalance, currency, err
}

/*
Used this function to update account's amount when a DEPOSIT or WITHDRAW event
is processed. New balance is calculated before the amount field is updated.

This functions panics if there is any error while processing the update.
*/
func updateAccount(db *sql.DB, p *Transaction, newBalance float64) {
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
func updateAccountsOnTransfer(db *sql.DB, p *Transaction, balanceSender, balanceReceiver float64) error {
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
