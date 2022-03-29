package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/brxyxn/go_mps_redcage/data"
	u "github.com/brxyxn/go_mps_redcage/utils"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

/*
swagger:route POST /clients/{client_id}/accounts/{account_id}/transactions transactions createTransaction
Return the id of the new transaction including the details of it
responses:
	201: transactionResponse
	400: badRequestErrorResponse
	500: internalErrorResponse

CreateTransaction handles POST request and returns the TransactionID of the new Transaction
*/
func (t *Handlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	t.l.Println("Handling POST Transactions")

	var v *data.Transaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		t.l.Println("[Error] CreateTransaction handler:", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := createTransaction(t.db, v); err != nil {
		t.l.Println("[Debug] CreateTransaction handler:", err)
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RespondWithJSON(w, http.StatusOK, v)
}

/*
swagger:route GET /clients/{client_id}/accounts/{account_id}/transactions transactions getTransactions
Return a list of the accounts listed by the {account_id}
responses:
	200: transactionsResponse
	400: badRequestErrorResponse
	404: notFoundErrorResponse

GetTransactions handles GET requests and returns a list of the transactions listed by an account
**{client_id} should be replaced by an authentication token (JWT) and validated with middleware.
*/
func (t *Handlers) GetTransactions(w http.ResponseWriter, r *http.Request) {
	t.l.Println("Handling GET Transactions")

	// var v *data.Transaction

	accountId, err := strconv.Atoi(mux.Vars(r)["account_id"])
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
	}

	// v.ReceiverAccountId = uint64(accountId)

	transactions, err := data.GetTransactions(t.db, accountId)
	if err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RespondWithJSON(w, http.StatusOK, transactions)
}

/*
This function validates the type of transaction that is requested
and perform the addition or substaction needed into the account's
table and creates the transaction registry.
*/
func createTransaction(db *sql.DB, p *data.Transaction) error {
	pAmount, err := decimal.NewFromString(p.Amount)
	if err != nil {
		return err
	}

	// When a deposit is processed the account balance is updated
	oldBalance, _, err := data.GetBalanceAccount(db, p.ReceiverAccountId)
	if err != nil {
		return err
	}

	previousBalance, err := decimal.NewFromString(oldBalance)

	// Validate what type of transaction is being performed and process the request.
	switch p.TransactionType {
	case data.DictTransactionType.Deposit:

		// Deposit represents an addition
		newBalance := previousBalance.Add(pAmount)

		// Update Account
		data.UpdateAccount(db, p, newBalance)
		err = data.RegisterTransaction(db, p)
		return err

	case data.DictTransactionType.Withdraw:
		// Deposit represents an substraction
		newBalance := previousBalance.Sub(pAmount)

		// Updating Account
		data.UpdateAccount(db, p, newBalance)
		err = data.RegisterTransaction(db, p)
		return err

	case data.DictTransactionType.Transfer:
		oldBalanceSenderAccount, currencySender, err := data.GetBalanceAccount(db, p.SenderAccountId)
		if err != nil {
			return err
		}
		oldBalanceReciverAccount, currencyReceiver, err := data.GetBalanceAccount(db, p.ReceiverAccountId)
		if err != nil {
			return err
		}

		// Validating the origin and destination account currencies match
		if currencySender != currencyReceiver {
			return errors.New("currencies don't match")
		}

		// Converting amounts to decimal.Decimal
		balanceSender, err := decimal.NewFromString(oldBalanceSenderAccount)
		if err != nil {
			return err
		}

		balanceReceiver, err := decimal.NewFromString(oldBalanceReciverAccount)

		newBalanceSenderAccount := balanceSender.Sub(pAmount)
		newBalanceReceiverAccount := balanceReceiver.Add(pAmount)

		data.UpdateAccountsOnTransfer(db, p, newBalanceSenderAccount, newBalanceReceiverAccount)
		data.RegisterTransaction(db, p)
	}

	return nil
}
