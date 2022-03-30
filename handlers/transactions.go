package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/brxyxn/go_mps_redcage/data"
	u "github.com/brxyxn/go_mps_redcage/utils"
	"github.com/gorilla/mux"
	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

/*
swagger:route POST /clients/{client_id}/accounts/{account_id}/transactions transactions createTransaction
Return the id of the new transaction including the details of it, include receiverand sender _accountId only with tranfer of funds
responses:
	201: transactionResponse
	400: badRequestErrorResponse
	500: internalErrorResponse

CreateTransaction handles POST request and returns the TransactionID of the new Transaction
*/
func (t *Handlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	u.LogInfo("Handling POST Transactions /clients/:client_id/accounts/:account_id/transactions")

	var v *data.Transaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		u.LogDebug("CreateTransaction handler:", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	accountID, err := strconv.Atoi(mux.Vars(r)["account_id"])
	if err != nil {
		u.LogDebug("CreateTransactions Handler:", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	if err := createTransaction(t.db, v, uint64(accountID)); err != nil {
		u.LogDebug("CreateTransaction handler:", err)
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
	u.LogInfo("Handling GET Transactions /clients/:client_id/accounts/:account_id/transactions")

	accountId, err := strconv.Atoi(mux.Vars(r)["account_id"])
	if err != nil {
		u.LogDebug("GetTransactions Handler:", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	transactions, err := data.GetTransactions(t.db, accountId)
	if err != nil {
		u.LogDebug("GetTransactions Handler:", err)
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
func createTransaction(db *sql.DB, p *data.Transaction, headAccountNumber uint64) error {
	pAmount, err := decimal.NewFromString(p.Amount)
	if err != nil {
		return err
	}

	// When a deposit is processed the account balance is updated
	oldBalance, _, err := data.GetBalanceAccount(db, headAccountNumber)
	if err != nil {
		return err
	}

	previousBalance, err := decimal.NewFromString(oldBalance)

	// Validate what type of transaction is being performed and process the request.
	switch p.TransactionType {
	case data.DictTransactionType.Deposit:
		// If a deposit is being processed the receiverAccount will register the transaction
		p.ReceiverAccountId = headAccountNumber
		p.SenderAccountId = 0

		// Formatting currency as string
		newBalance := previousBalance.Add(pAmount)
		strNewBalance := formatAmountDecimal(newBalance)

		// Updating Account and registering transaction
		err = data.UpdateAccount(db, p, strNewBalance)
		if err != nil {
			return err
		}

		// Formatting currency as string
		p.Amount = formatAmountString(p.Amount)
		err = data.RegisterTransaction(db, p)
		if err != nil {
			return err
		}
		return nil

	case data.DictTransactionType.Withdraw:
		// If a withdraw is being processed the senderAccount will register the transaction
		p.ReceiverAccountId = headAccountNumber
		p.SenderAccountId = 0

		// Formatting currency as string
		newBalance := previousBalance.Sub(pAmount)
		strNewBalance := formatAmountDecimal(newBalance)

		// Updating Account and registering transaction
		err = data.UpdateAccount(db, p, strNewBalance)
		if err != nil {
			return err
		}

		// Formatting currency as string
		p.Amount = formatAmountString(p.Amount)
		err = data.RegisterTransaction(db, p)
		if err != nil {
			return err
		}
		return nil

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

		// Validating if the amount is zero
		if balanceSender.IsZero() {
			u.LogDebug("Sender's balance is zero", balanceSender)
			return errors.New("balance of sender's account has no funds")
		}

		balanceReceiver, err := decimal.NewFromString(oldBalanceReciverAccount)

		newBalanceSenderAccount := balanceSender.Sub(pAmount)
		newBalanceReceiverAccount := balanceReceiver.Add(pAmount)

		// Formatting currency as string
		newBalanceSenderAccountStr := formatAmountDecimal(newBalanceSenderAccount)
		newBalanceReceiverAccountStr := formatAmountDecimal(newBalanceReceiverAccount)

		err = data.UpdateAccountsOnTransfer(db, p, newBalanceSenderAccountStr, newBalanceReceiverAccountStr)
		if err != nil {
			return err
		}

		p.Amount = formatAmountString(p.Amount)
		err = data.RegisterTransaction(db, p)
		if err != nil {
			return err
		}
	}

	return nil
}

// This function returns the input as string with the format: #.##
func formatAmountDecimal(a decimal.Decimal) string {
	parsed := accounting.FormatNumberDecimal(a, 2, ",", ".")
	re := regexp.MustCompile(`([0-9]{1,3})+`)

	// execute regexp
	miles := re.FindAllString(parsed, 8)

	// Format decimal, in the last position
	leng := len(miles) - 1
	miles[leng] = "." + miles[leng]

	// return the slice joined as string
	return strings.Join(miles[:], "")
}

// This function returns the input as string with the format: #.##
func formatAmountString(a string) string {
	amount, _ := decimal.NewFromString(a)

	return formatAmountDecimal(amount)
}
