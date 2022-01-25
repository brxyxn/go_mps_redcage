package main

import (
	"database/sql"

	"github.com/gorilla/mux"
)

/*
	Main App Structures
*/

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

/*
	Client Structures
*/

type Client struct {
	Id        uint64 `json:"clientId"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Username  string `json:"username"`
	Active    bool   `json:"active,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

/*
type Profile struct {
	Id                uint64      `json:"profileId"`
	DateOfBirth       string   `json:"dateOfBirth"`
	ProfilePicture    string      `json:"profilePicture"` // base64
	Contact           ContactInfo `json:"contact"`
	PhysicalAddress   Address     `json:"physicalAddress"`
	MailingAddress    Address     `json:"mailingAddress"`
	AdditionalDetails string      `json:"additionalDetails"`
	ClientId          uint64      `json:"clientId"`
}

type ContactInfo struct {
	PrimaryEmail         string      `json:"primaryEmail"`
	SecondaryEmail       string      `json:"secondaryEmail"`
	PrimaryPhoneNumber   PhoneNumber `json:"primaryPhoneNumber"`
	SecondaryPhoneNumber PhoneNumber `json:"secondaryPhoneNumber"`
}

type PhoneNumber struct {
	AreaCode          string `json:"areaCode"`
	Number            string `json:"number"`
	ContactByText     bool   `json:"contactByText"`
	ContactByVoice    bool   `json:"contactByVoice"`
	ContactByWhatsApp bool   `json:"contactByWhatsApp"`
}

type Address struct {
	AddressLineOne string `json:"addressLineOne"`
	AddressLineTwo string `json:"addressLineTwo"`
	City           string `json:"city"`
	State          string `json:"state"`
	Country        string `json:"country"`
	ClientId       uint64 `json:"clientId"`
}
*/

/*
	Account Structures
*/

type Accounts []Account

type Account struct {
	Id          uint64      `json:"accountId"`
	Balance     float32     `json:"balance"`
	Currency    string      `json:"currency"`
	AccountType AccountType `json:"accountType"`
	Active      bool        `json:"active,omitempty"`
	ClientId    uint64      `json:"clientId"`
	CreatedAt   string      `json:"createdAt,omitempty"`
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
