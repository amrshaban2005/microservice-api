package service

import (
	"context"

	"github.com/amrshaban2005/microservice-api/domain"
	"github.com/amrshaban2005/microservice-api/dto"
	"github.com/amrshaban2005/microservice-api/errs"
)


type AccountService interface{
	NewAccount(ctx context.Context,account dto.NewAccountRequest) (*dto.NewAccountResposne,*errs.AppError)
}

type DefaultAccountService struct{
	repo domain.AccountRepositoryDB
}

func NewAccountService(repo domain.AccountRepositoryDB)DefaultAccountService{
	return DefaultAccountService{repo}
}

func(d DefaultAccountService) NewAccount(ctx context.Context,account dto.NewAccountRequest) (*dto.NewAccountResposne,*errs.AppError){
    accountResponse,err :=d.repo.Save(ctx,domain.NewAccount(account.CustomerId,account.AccountType,account.Amount))
	if err!=nil{
		return nil,err
	}

	return accountResponse.ToNewAccountResponseDto(),nil
}