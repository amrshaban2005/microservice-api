package domain

import (
	"strconv"
	"time"

	"github.com/amrshaban2005/microservice-api/dto"
)

type Transaction struct {
	TransactioId    int32     `db:"transaction_id"`
	AccountId       int32     `db:"account_id"`
	Amount          float64   `db:"amount"`
	TransactionType string    `db:"transaction_type"`
	TransactionDate time.Time `db:"transaction_date"`
}

func (t Transaction) ToDto() *dto.TransactionResponse {
	return &dto.TransactionResponse{
		TransactioId:    strconv.Itoa(int(t.TransactioId)),
		AccountId:       strconv.Itoa(int(t.AccountId)),
		Amount:          t.Amount,
		TransactionType: t.TransactionType,
		TransactionDate: t.TransactionDate.Format("2006-01-02"),
	}
}

func (t Transaction) IsWithdrawal() bool{
	if t.TransactionType == "withdraw" {
		return  true
	}
	return false
}

func NewTransaction(accountId string, amount float64, transactiontype string) Transaction {
	accountIdInt, _ := strconv.Atoi(accountId)

	return Transaction{
		TransactioId:    0,
		AccountId:       int32(accountIdInt),
		Amount:          amount,
		TransactionType: transactiontype,
		TransactionDate: time.Now(),
	}
}
