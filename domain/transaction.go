package domain

import "time"

type Transaction struct{
	TransactioId int32 `db:"transaction_id"`
	AccountId int32 `db:"account_id"`
	Amount float32 `db:"amount"`
	TransactionType string `db:"transaction_type"`
	TransactionDate time.Time `db:"transaction_date"`
}