package domain

import (
	"context"
	"errors"

	"github.com/amrshaban2005/banking-lib/errs"
	"github.com/amrshaban2005/banking-lib/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type CustomerRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewCustomerRepositoryDB(pool *pgxpool.Pool)CustomerRepositoryDB{
	return CustomerRepositoryDB{pool}
}

func (c CustomerRepositoryDB) FindAll(ctx context.Context,status string) ([]Customer,*errs.AppError){
	customers := make([]Customer,0)

	findAllQuery := "SELECT customer_id, name, date_of_birth, city, zipcode, status FROM customers"
	args :=[]any{}
	
	if status != "" {
		findAllQuery += " WHERE status = $1"
		args = append(args, status)
	}

	rows, err := c.pool.Query(ctx, findAllQuery, args...)
	if err != nil {
		logger.Error("error while fetching customer table "+ err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	
	defer rows.Close()
	
	
	customers, err = pgx.CollectRows(rows, pgx.RowToStructByName[Customer])
	if err != nil {
		logger.Error("error while mapping customer table "+ err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return customers, nil
}

func (d CustomerRepositoryDB) ById(ctx context.Context,id string) (*Customer,*errs.AppError){
	
	findQuery := "SELECT customer_id, name, date_of_birth, city, zipcode, status FROM customers where customer_id=$1"
	var customer Customer

	rows,err := d.pool.Query(ctx, findQuery, id)
	if err != nil {
		logger.Error("error while fetching customer by id "+ err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	defer rows.Close()
	
	
	customer, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Customer])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil,errs.NewNotFoundError("Customer not found")
		}
		logger.Error("error while mapping customer record "+ err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return &customer, nil
}