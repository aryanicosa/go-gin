package db

import (
	"context"
	"database/sql"
	"github.com/aryanicosa/go_gin_simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

//func deleteTestAccount(a Account) {
//	_ = testQueries.DeleteAccount(context.Background(), a.ID)
//}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
	//defer deleteTestAccount(accountCreated)
}

func TestGetAccount(t *testing.T) {
	accountCreated := createRandomAccount(t)
	//defer deleteTestAccount(accountCreated)
	accountGet, err := testQueries.GetAccount(context.Background(), accountCreated.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountGet)

	require.Equal(t, accountCreated.ID, accountGet.ID)
	require.Equal(t, accountCreated.Owner, accountGet.Owner)
	require.Equal(t, accountCreated.Balance, accountGet.Balance)
	require.Equal(t, accountCreated.Currency, accountGet.Currency)
	require.WithinDuration(t, accountCreated.CreatedAt, accountGet.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	accountCreated := createRandomAccount(t)
	//defer deleteTestAccount(accountCreated)

	arg := UpdateAccountParams{
		ID:      accountCreated.ID,
		Balance: util.RandomMoney(),
	}

	accountUpdated, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountUpdated)

	require.Equal(t, accountCreated.ID, accountUpdated.ID)
	require.Equal(t, accountCreated.Owner, accountUpdated.Owner)
	require.Equal(t, arg.Balance, accountUpdated.Balance)
	require.Equal(t, accountCreated.Currency, accountUpdated.Currency)
	require.WithinDuration(t, accountCreated.CreatedAt, accountUpdated.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	accountCreated := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), accountCreated.ID)
	require.NoError(t, err)

	accountCheck, err := testQueries.GetAccount(context.Background(), accountCreated.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, accountCheck)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  10,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 10)

	//for _, account := range accounts {
	//	require.NotEmpty(t, account)
	//	deleteTestAccount(account)
	//}
}
