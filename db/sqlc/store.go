package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db     DBTX
	realDB *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		db:      pool,
		realDB:  pool,
		Queries: New(pool),
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
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer      Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry     Entry    `json:"from_entry"`
	ToEntry       Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TranferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Create transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			fmt.Println("ERROR creating CreateTransfer:", err)
			return err
		}

		// Create from entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			fmt.Println("ERROR creating FromEntry:", err)
			return err
		}

		// Create to entry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			fmt.Println("ERROR creating ToEntry:", err)
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				AccountID: arg.FromAccountID,
				Amount:    -arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				AccountID: arg.ToAccountID,
				Amount:    arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				AccountID: arg.ToAccountID,
				Amount:    arg.Amount,
			})
		
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				AccountID: arg.FromAccountID,
				Amount:    -arg.Amount,
			})
			
			if err != nil {
				return err
			}
		}

		// Fetch updated fromAccount
		// result.FromAccount, err = q.GetAccount(ctx, arg.FromAccountID)
		// if err != nil {
		// 	fmt.Println("ERROR fetching FromAccount:", err)
		// 	return err
		// }

		// Fetch updated toAccount
		// result.ToAccount, err = q.GetAccount(ctx, arg.ToAccountID)
		// if err != nil {
		// 	fmt.Println("ERROR fetching ToAccount:", err)
		// }

		return nil
	})

	return result, err
}
