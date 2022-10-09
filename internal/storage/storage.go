package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/seth2810/money-processing/internal/storage/queries"
	"github.com/shopspring/decimal"
)

var (
	ErrClientNotFound  = errors.New("client not found")
	ErrAccountNotFound = errors.New("account not found")
)

type Storage struct {
	db *sql.DB
	qs *queries.Queries
}

func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
		qs: queries.New(db),
	}
}

func (s *Storage) CreateClient(ctx context.Context, email string) (*queries.Client, error) {
	return s.qs.CreateClient(ctx, email)
}

func (s *Storage) GetClient(ctx context.Context, id int32) (*queries.Client, error) {
	client, err := s.qs.GetClient(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrClientNotFound
	}

	return client, err
}

func (s *Storage) CreateAccount(
	ctx context.Context,
	clientID int32,
	currency queries.CurrencyTicker,
) (*queries.Account, error) {
	return s.qs.CreateAccount(ctx, queries.CreateAccountParams{
		ClientID: clientID,
		Currency: currency,
	})
}

func (s *Storage) GetAccount(ctx context.Context, id int32) (*queries.Account, error) {
	client, err := s.qs.GetAccount(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrAccountNotFound
	}

	return client, err
}

func (s *Storage) Deposit(
	ctx context.Context,
	amount decimal.Decimal,
	accountID int32,
) (*queries.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	qtx := s.qs.WithTx(tx)

	err = qtx.DepositToAccount(ctx, queries.DepositToAccountParams{
		Amount:    amount,
		AccountID: accountID,
	})
	if err != nil {
		return nil, err
	}

	transaction, err := qtx.CreateDepositTransaction(ctx, queries.CreateDepositTransactionParams{
		Amount:    amount,
		AccountID: accountID,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Storage) Withdraw(
	ctx context.Context,
	amount decimal.Decimal,
	accountID int32,
) (*queries.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	qtx := s.qs.WithTx(tx)

	err = qtx.WithdrawFromAccount(ctx, queries.WithdrawFromAccountParams{
		Amount:    amount,
		AccountID: accountID,
	})
	if err != nil {
		return nil, err
	}

	transaction, err := qtx.CreateWithdrawTransaction(ctx, queries.CreateWithdrawTransactionParams{
		Amount:    amount,
		AccountID: accountID,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Storage) Transfer(
	ctx context.Context,
	amount decimal.Decimal,
	fromAccountID, toAccountID int32,
) (*queries.Transaction, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	qtx := s.qs.WithTx(tx)

	err = qtx.WithdrawFromAccount(ctx, queries.WithdrawFromAccountParams{
		Amount:    amount,
		AccountID: fromAccountID,
	})
	if err != nil {
		return nil, err
	}

	err = qtx.DepositToAccount(ctx, queries.DepositToAccountParams{
		Amount:    amount,
		AccountID: toAccountID,
	})
	if err != nil {
		return nil, err
	}

	transaction, err := qtx.CreateTransferTransaction(ctx, queries.CreateTransferTransactionParams{
		Amount:        amount,
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Storage) GetTransactions(ctx context.Context, accountID int32) ([]*queries.Transaction, error) {
	qs := queries.New(s.db)

	txs, err := qs.ListTransactions(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// create empty array to avoid nil values when no rows found
	if len(txs) == 0 {
		return []*queries.Transaction{}, nil
	}

	return txs, nil
}
