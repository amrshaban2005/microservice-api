package service

import (
	"context"

	"github.com/amrshaban2005/microservice-api/domain"
	"github.com/amrshaban2005/microservice-api/dto"
	"github.com/amrshaban2005/microservice-api/errs"
)

type AccountService interface {
	NewAccount(ctx context.Context, account dto.NewAccountRequest) (*dto.NewAccountResposne, *errs.AppError)
	MakeTransaction(ctx context.Context, transaction dto.NewTransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRepositoryDB
}

func NewAccountService(repo domain.AccountRepositoryDB) DefaultAccountService {
	return DefaultAccountService{repo}
}

func (d DefaultAccountService) NewAccount(ctx context.Context, account dto.NewAccountRequest) (*dto.NewAccountResposne, *errs.AppError) {
	err := account.Validate()
	if err != nil {
		return nil, err
	}

	accountResponse, err := d.repo.Save(ctx, domain.NewAccount(account.CustomerId, account.AccountType, account.Amount))
	if err != nil {
		return nil, err
	}

	return accountResponse.ToNewAccountResponseDto(), nil
}

func (d DefaultAccountService) MakeTransaction(ctx context.Context, transaction dto.NewTransactionRequest) (*dto.TransactionResponse, *errs.AppError) {

	err := transaction.Validate()
	if err != nil {
		return nil, err
	}

	if transaction.IsTransactTypeWithdraw() {
		account, err := d.repo.FindById(ctx, transaction.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.CanWithdraw(transaction.Amount) {
			return nil, errs.NewValidationError("Insufficient amount")
		}
	}

	transactionResonse, err := d.repo.SaveTransaction(ctx, domain.NewTransaction(transaction.AccountId, transaction.Amount, transaction.TransactionType))
	if err != nil {
		return nil, err
	}

	return transactionResonse.ToDto(), nil
}
