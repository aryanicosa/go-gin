package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountA := createRandomAccount(t)
	accountB := createRandomAccount(t)

	// waiting database transaction needed to be careful. it's easy to write but also easily become a nightmare
	// handle it with go routine/concurrent

	// run n concurrent transfer transaction
	n := 5
	amount := int64(10)

	// verified error and result, send them back to go routine
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() { // start go routine
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountA.ID,
				ToAccountID:   accountB.ID,
				Amount:        amount,
			})

			// send to channel
			errs <- err
			results <- result
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		// recipe from channel
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, accountA.ID, transfer.FromAccountID)
		require.Equal(t, accountB.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, accountA.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, accountB.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: check account balance
	}
}
