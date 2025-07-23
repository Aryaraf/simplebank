package db

import (
	"context"
	"fmt"
	
	"github.com/jackc/pgx/v5"
)

type Store struct {
	*Queries
	db DBTX
	realDB *pgx.Conn
}

func NewStore(db DBTX) *Store {
	return &Store{
		db: db,
		realDB: db.(*pgx.Conn), 
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.realDB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit(ctx)
}

type TranferTxParams struct {
	FromAccountID 	int64 `json:"from_account_id"`
	ToAccountID 	int64 `json:"to_account_id"`
	Amount 			int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer 		Transfer	`json:"transfer"`
	FromAccountID 	Account 	`json:"from_account"`
	ToAccountID 	Account 	`json:"to_account"`
	FromEntry 		Entry 		`json:"from_entry"`
	ToEntry 		Entry 		`json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TranferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:  arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: 	arg.FromAccountID,
			Amount: 	-arg.Amount,
		})

		if err != nil{
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: 	arg.ToAccountID,
			Amount: 	arg.Amount,
		})

		if err != nil{
			return err
		}

		return nil
	})

	return result, err
}