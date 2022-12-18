package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provided all functions to execute db queries and transactions
type Store struct {
	*Queries // composition, prefered way to extend struct functionality in Go instead of inheritance by embedding query inside the store
	// or individual query provided by Queries will be available to store.
	db *sql.DB // provide new DB transaction
}

// NewStore create a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer Transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"FromAccountID"`
	ToAccountID   int64 `json:"ToAccountID"`
	Amount        int64 `json:"Amount"`
}

// TransferTxResult is the result of the transfer Transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"Transfer"`
	FromAccount Account  `json:"FromAccount"`
	ToAccount   Account  `json:"ToAccount"`
	FromEntry   Entry    `json:"FromEntry"`
	ToEntry     Entry    `json:"ToEntry"`
}

// TransferTx performs amoney transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: update accounts' balance

		return nil
	})

	return result, err
}
