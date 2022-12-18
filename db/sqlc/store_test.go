package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	accountA := createRandomAccount(t)
	accountB := createRandomAccount(t)

	// waiting database transaction needed to be careful. it's easy to write but also easily become a nightmare
	// handle it with go routine/concurrent

	// logging
	fmt.Println("before : ", "accountA", accountA.Balance, "accountB", accountB.Balance)

	// run n concurrent transfer transaction
	n := 5
	amount := int64(10)

	// verified error and result, send them back to go routine
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() { // start go routine
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
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
	existed := make(map[int]bool)
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

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, accountA.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, accountB.ID, toAccount.ID)

		//check account's balance

		// logging
		fmt.Println("transaction : ", "from", fromAccount.Balance, "to", toAccount.Balance)
		diffA := accountA.Balance - fromAccount.Balance
		diffB := toAccount.Balance - accountB.Balance
		require.Equal(t, diffA, diffB)
		require.True(t, diffA > 0)
		require.True(t, diffA%amount == 0) // 1 = amount, 2 * amount .... n * amount

		k := int(diffA / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final updated balance
	updatedAccountA, err := testQueries.GetAccount(context.Background(), accountA.ID)
	require.NoError(t, err)

	updatedAccountB, err := testQueries.GetAccount(context.Background(), accountB.ID)
	require.NoError(t, err)

	// logging
	fmt.Println("after : ", "accountA", updatedAccountA.Balance, "accountB", updatedAccountB.Balance)
	require.Equal(t, accountA.Balance-int64(n)*amount, updatedAccountA.Balance)
	require.Equal(t, accountB.Balance+int64(n)*amount, updatedAccountB.Balance)
}

func TestTransferTxDeacLock(t *testing.T) {
	store := NewStore(testDB)

	accountA := createRandomAccount(t)
	accountB := createRandomAccount(t)

	// waiting database transaction needed to be careful. it's easy to write but also easily become a nightmare
	// handle it with go routine/concurrent

	// logging
	fmt.Println("before : ", "accountA", accountA.Balance, "accountB", accountB.Balance)

	// run n concurrent transfer transaction
	n := 10
	amount := int64(10)

	// verified error and result, send them back to go routine
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := accountA.ID
		toAccountId := accountA.ID

		if i%2 == 1 {
			fromAccountId = accountB.ID
			toAccountId = accountB.ID
		}
		go func() { // start go routine
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			// send to channel
			errs <- err
		}()
	}

	// check results
	for i := 0; i < n; i++ {
		// recipe from channel
		err := <-errs
		require.NoError(t, err)

	}

	// check final updated balance
	updatedAccountA, err := testQueries.GetAccount(context.Background(), accountA.ID)
	require.NoError(t, err)

	updatedAccountB, err := testQueries.GetAccount(context.Background(), accountB.ID)
	require.NoError(t, err)

	// logging
	fmt.Println("after : ", "accountA", updatedAccountA.Balance, "accountB", updatedAccountB.Balance)
	require.Equal(t, accountA.Balance, updatedAccountA.Balance)
	require.Equal(t, accountB.Balance, updatedAccountB.Balance)
}
