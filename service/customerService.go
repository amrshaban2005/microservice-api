package service

import (
	"context"

	"github.com/amrshaban2005/microservice-api/domain"
	"github.com/amrshaban2005/microservice-api/dto"
	"github.com/amrshaban2005/banking-lib/errs"
)

type CustomerService interface{
	GetAllCustomers(ctx context.Context,status string)([]dto.CustomerResponse,*errs.AppError)
	GetCustomer(ctx context.Context,id string)(*dto.CustomerResponse,*errs.AppError)
}


type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func NewCustomerService(repo domain.CustomerRepository)DefaultCustomerService{
	return DefaultCustomerService{repo:repo}
}

func (s DefaultCustomerService) GetAllCustomers(ctx context.Context,status string)([]dto.CustomerResponse,*errs.AppError){
  switch status {
	case "active":
		status = "1"
	case "inactive":
		status = "0"
	default:
		status = ""
	}
	

	customers,err := s.repo.FindAll(ctx,status)
	if err !=nil {
		return nil,err
	}

	response := make([]dto.CustomerResponse,0)

	for _, customer := range customers {
		response = append(response, customer.ToDto())
	}

	return response, nil
}

func (s DefaultCustomerService) GetCustomer(ctx context.Context,id string)(*dto.CustomerResponse,*errs.AppError){
	

	customer,err := s.repo.ById(ctx,id)
	if err !=nil {
		return nil,err
	}

    response :=customer.ToDto()

	return &response, nil
}