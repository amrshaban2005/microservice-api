package domain

import (
	"context"
	"strconv"
	"time"

	"github.com/amrshaban2005/microservice-api/dto"
	"github.com/amrshaban2005/banking-lib/errs"
)

type Account struct {
	AccountId   int32     `db:"account_id"`
	CusomerId   int32     `db:"customer_id"`
	OpeningDate time.Time `db:"opening_date"`
	AccountType string    `db:"account_type"`
	Amount      float64   `db:"amount"`
	Status      int8      `db:"status"`
}

type AccountRepository interface {
	Save(ctx context.Context, account Account) (*Account, *errs.AppError)
	SaveTransaction(ctx context.Context, transaction Transaction) (*Transaction, *errs.AppError)
	FindById(ctx context.Context, id string) (*Account, *errs.AppError)
}

func (a Account) ToNewAccountResponseDto() *dto.NewAccountResposne {
	return &dto.NewAccountResposne{AccountId: strconv.Itoa(int(a.AccountId))}
}

func (a Account) CanWithdraw(amount float64) bool {
	return a.Amount >= amount
}

func NewAccount(customerId string, accountType string, amount float64) Account {
	customerIdInt, _ := strconv.Atoi(customerId)

	return Account{
		CusomerId:   int32(customerIdInt),
		OpeningDate: time.Now(),
		AccountType: accountType,
		Amount:      amount,
		Status:      1,
	}

}
