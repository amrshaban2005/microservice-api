package domain

import (
	"context"

	"github.com/amrshaban2005/microservice-api/errs"
	"github.com/amrshaban2005/microservice-api/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewAccountRepositoryDB(pool *pgxpool.Pool) AccountRepositoryDB {
	return AccountRepositoryDB{pool}
}

func (a AccountRepositoryDB) Save(ctx context.Context, account Account) (*Account, *errs.AppError) {
	sqlInsert := `INSERT INTO accounts (customer_id, account_type, amount)
		VALUES ($1, $2, $3)
		RETURNING account_id`

	err := a.pool.QueryRow(ctx, sqlInsert, account.CusomerId, account.AccountType, account.Amount).Scan(&account.AccountId)

	if err != nil {
		logger.Error("Error while save new account " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return &account, nil
}
