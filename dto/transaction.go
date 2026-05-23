package dto

import "github.com/amrshaban2005/banking-lib/errs"

const DEPOSIT = "deposit"
const WITHDRAW = "withdraw"

type NewTransactionRequest struct {
	AccountId       string  `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

func (t NewTransactionRequest) IsTransactTypeDeposit() bool {
	return t.TransactionType == DEPOSIT
}

func (t NewTransactionRequest) IsTransactTypeWithdraw() bool {
	return t.TransactionType == WITHDRAW
}

func (t NewTransactionRequest) Validate() *errs.AppError {
	if !t.IsTransactTypeDeposit() && !t.IsTransactTypeWithdraw() {
		return errs.NewValidationError("Transactions can be only deposit or withdraw")
	}

	if t.Amount <= 0 {
		return errs.NewValidationError("Amount can not be zero or less")
	}
	return nil
}

type TransactionResponse struct {
	TransactioId    string  `json:"transaction_id"`
	AccountId       string  `json:"account_id"`
	Amount          float64 `json:"new_balance"`
	TransactionType string  `json:"transaction_type"`
	TransactionDate string  `json:"transaction_date"`
}
