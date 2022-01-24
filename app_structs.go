package main

import (
	"database/sql"
	"time"

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
	Active    bool   `json:"active"`
}

type Profile struct {
	Id                uint64      `json:"profileId"`
	DateOfBirth       time.Time   `json:"dateOfBirth"`
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

/*
	Account Structures
*/

type Account struct {
	Id          uint64
	Currency    string
	AccountType uint16
	Balance     uint32
	Active      bool
	ClientId    uint64
}

/*
	Transaction Structures
*/

type Transaction struct {
	Id              uint64
	Amount          uint32
	TransactionType TransactionType
	Description     string
	DateTime        time.Time
	FromAccountId   uint64
	ToAccountId     uint64
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
