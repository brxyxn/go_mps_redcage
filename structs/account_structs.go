package model

type Account struct {
	Id          uint64
	Currency    string
	AccountType uint16
	Balance     uint32
	Active      bool
	ClientId    uint64
}
