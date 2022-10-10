//go:build integration
// +build integration

package tests

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/seth2810/money-processing/internal/storage"
	"github.com/seth2810/money-processing/internal/storage/queries"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	suite.Suite
	db      *sql.DB
	storage *storage.Storage
}

func (s *StorageTestSuite) SetupSuite() {
	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	require.NoError(s.T(), err)

	s.db = db
	s.storage = storage.New(db)

	s.T().Name()
}

func (s *StorageTestSuite) TearDownSuite() {
	err := s.db.Close()

	require.NoError(s.T(), err)
}

func (s *StorageTestSuite) TestAccountsTransactionsConcurrency() {
	ctx := context.TODO()

	c, err := s.storage.CreateClient(ctx, fmt.Sprintf("%s@example.com", time.Now()))
	require.NoError(s.T(), err)

	acc1, err := s.storage.CreateAccount(ctx, c.ID, queries.CurrencyTickerUSD)
	require.NoError(s.T(), err)

	acc2, err := s.storage.CreateAccount(ctx, c.ID, queries.CurrencyTickerUSD)
	require.NoError(s.T(), err)

	_, err = s.storage.Deposit(ctx, decimal.NewFromFloat32(200.02), acc1.ID)
	require.NoError(s.T(), err)

	var wg sync.WaitGroup

	wg.Add(300)

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			s.storage.Deposit(ctx, decimal.NewFromFloat32(3.0003), acc1.ID)
		}()

		go func() {
			defer wg.Done()

			s.storage.Transfer(ctx, decimal.NewFromFloat32(1.0001), acc1.ID, acc2.ID)
		}()

		go func() {
			defer wg.Done()

			s.storage.Withdraw(ctx, decimal.NewFromFloat32(1.0001), acc1.ID)
		}()
	}

	wg.Wait()

	acc1, err = s.storage.GetAccount(ctx, acc1.ID)
	require.NoError(s.T(), err)

	acc2, err = s.storage.GetAccount(ctx, acc2.ID)
	require.NoError(s.T(), err)

	require.EqualValues(s.T(), "300.0300", acc1.Balance.StringFixed(4))
	require.EqualValues(s.T(), "100.0100", acc2.Balance.StringFixed(4))
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
