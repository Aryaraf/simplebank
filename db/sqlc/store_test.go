package db

import (
	"context"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5 
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func(){
			result, err := store.TransferTx(context.Background(), TranferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})

			errs <- err
			results <- result

		}()
	}		
}