package model

import (
	"time"
)

type Transactions struct {
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
