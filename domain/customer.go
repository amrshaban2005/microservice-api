package domain

import (
	"context"
	"strconv"
	"time"

	"github.com/amrshaban2005/microservice-api/dto"
	"github.com/amrshaban2005/banking-lib/errs"
)

type Customer struct {
	Id          int32     `db:"customer_id"`
	Name        string
	City        string
	Zipcode     string
	DateofBirth time.Time `db:"date_of_birth"`
	Status      int16
}

type CustomerRepository interface {
	FindAll(ctx context.Context,status string)([]Customer,*errs.AppError)
	ById(ctx context.Context,id string)(*Customer, *errs.AppError)
}

func (c Customer) StatusAsText() string{
	statusText :="active"
	if c.Status == 0{
		statusText = "inactive"
	}
	return statusText

}

func (c Customer) ToDto() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:          strconv.Itoa(int(c.Id)),
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateofBirth: c.DateofBirth.Format("2006-01-02"),
		Status:      c.StatusAsText(),
	}
}