package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T, accountA, accountB Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: accountA.ID,
		ToAccountID:   accountB.ID,
		Amount:        10,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
	require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
	require.Equal(t, transfer.Amount, arg.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	accountA := createRandomAccount(t)
	accountB := createRandomAccount(t)
	createRandomTransfer(t, accountA, accountB)
}

func TestGetTransfer(t *testing.T) {
	accountA := createRandomAccount(t)
	accountB := createRandomAccount(t)
	transferCreated := createRandomTransfer(t, accountA, accountB)

	transfer, err := testQueries.GetTransfer(context.Background(), transferCreated.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.ID, transferCreated.ID)
	require.Equal(t, transfer.FromAccountID, transferCreated.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transferCreated.ToAccountID)
	require.Equal(t, transfer.Amount, transferCreated.Amount)
	require.WithinDuration(t, transfer.CreatedAt, transferCreated.CreatedAt, time.Second)
}

func TestGetListTransfer(t *testing.T) {
	accountA := createRandomAccount(t)
	accountB := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, accountA, accountB)
	}

	arg := ListTransfersParams{
		FromAccountID: accountA.ID,
		ToAccountID:   accountB.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	}
}
