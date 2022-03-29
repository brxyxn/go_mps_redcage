package data

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"

	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

/* Encoding/Decoding */

func (p *Accounts) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Accounts) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

type Decimal struct {
	R      accounting.Accounting
	Symbol string
	Amount decimal.Decimal
	Money  string
}

const (
	USD = "USD"
	MXN = "MXN"
	COP = "COP"
)

func data() {}

func (v *Decimal) ParseDecimal(str string) error {
	currencyRe := regexp.MustCompile(`^[A-Z]{3}`)
	amountRe := regexp.MustCompile(`([0-9]+\.[0-9]{2})`)

	symbol := currencyRe.FindString(str)
	if symbol == "" {
		return errors.New("Unable to find a currency symbol")
	}
	amountStr := amountRe.FindString(str)
	if amountStr == "" {
		return errors.New("Unable to find an amount value")
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return errors.New("Unable to find an amount value " + err.Error())
	}
	v.Symbol = symbol
	v.Amount = amount

	v.FormatMoney()

	return nil
}

func (v *Decimal) FormatMoney() {
	v.R = *accounting.NewAccounting(v.Symbol, 2, ",", ".", "%s %v", "%s (%v)", "%s 0.00")
	v.Money = v.R.FormatMoneyDecimal(v.Amount)
}
