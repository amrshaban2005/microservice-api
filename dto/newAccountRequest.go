package dto

import (
	"strings"

	"github.com/amrshaban2005/banking-lib/errs"
)

type NewAccountRequest struct{
	CustomerId string  `json:"customer_id"`
	AccountType string `json:"account_type"`
	Amount float64 `json:"amount"`
}

func (r NewAccountRequest) Validate() *errs.AppError{
	if r.Amount < 5000{
		return errs.NewValidationError("Amount should be => 5000")
	}

	if strings.ToLower(r.AccountType) != "saving" && strings.ToLower(r.AccountType) != "checking"{
		return errs.NewValidationError("Account type should be saving or checking")
	} 
	return nil
}