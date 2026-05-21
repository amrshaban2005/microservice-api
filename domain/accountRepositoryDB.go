package domain

import (
	"context"
	"errors"
	"strconv"

	"github.com/amrshaban2005/microservice-api/errs"
	"github.com/amrshaban2005/microservice-api/logger"
	"github.com/jackc/pgx/v5"
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

func (a AccountRepositoryDB) SaveTransaction(ctx context.Context, transaction Transaction) (*Transaction, *errs.AppError) {

	tx, err := a.pool.Begin(ctx)
	if err != nil {
		logger.Error("Error while open a sql transaction " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	insertSql := `INSERT INTO transactions(
	              account_id, amount, transaction_type)
	              VALUES ($1, $2, $3) RETURNING *`

	rows, err := tx.Query(ctx, insertSql, transaction.AccountId, transaction.Amount, transaction.TransactionType)

	if err != nil {
		logger.Error("Error while save new transaction " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	var transactionResonse Transaction
	transactionResonse, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Transaction])
	if err != nil {
		logger.Error("error while mapping transction record " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	defer func() {
		rows.Close()
		tx.Rollback(ctx)
	}()

	var updateSql string
	if transaction.IsWithdrawal() {
		updateSql = `UPDATE accounts
					 SET  amount= amount - $1 where account_id=$2`
	} else {
		updateSql = `UPDATE accounts
					 SET  amount= amount + $1 where account_id=$2`
	}

	_, err = tx.Exec(ctx, updateSql, transaction.Amount, transaction.AccountId)

	if err != nil {
		logger.Error("Error while update account balance " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error("error while commit the transaction " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	account, appErr := a.FindById(ctx, strconv.Itoa(int(transaction.AccountId)))
	if appErr != nil {
		return nil, appErr
	}
	transactionResonse.Amount = account.Amount

	return &transactionResonse, nil

}

func (a AccountRepositoryDB) FindById(ctx context.Context, id string) (*Account, *errs.AppError) {
	sqlSelect := `SELECT account_id, customer_id, opening_date, account_type, amount, status
	             FROM accounts where account_id = $1`

	rows, err := a.pool.Query(ctx, sqlSelect, id)
	if err != nil {
		logger.Error("Error while fetching account record " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected databae error")
	}
	var account Account

	defer rows.Close()

	account, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewNotFoundError("Account not found")
		}
		logger.Error("Error while fetching account record " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected databae error")
	}

	return &account, nil
}
