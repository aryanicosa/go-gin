package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    10,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, arg.AccountID)
	require.Equal(t, entry.Amount, arg.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entryCreated := createRandomEntry(t, account)

	entry, err := testQueries.GetEntry(context.Background(), entryCreated.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.ID, entryCreated.ID)
	require.Equal(t, entry.AccountID, entryCreated.AccountID)
	require.Equal(t, entry.Amount, entryCreated.Amount)
	require.WithinDuration(t, entry.CreatedAt, entryCreated.CreatedAt, time.Second)
}

func TestGetListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 5; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
