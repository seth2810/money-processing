// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package queries

import (
	"database/sql/driver"
	"fmt"
	"time"

	types "github.com/seth2810/money-processing/internal/storage/types"
	decimal "github.com/shopspring/decimal"
)

type CurrencyTicker string

const (
	CurrencyTickerUSD CurrencyTicker = "USD"
	CurrencyTickerCOP CurrencyTicker = "COP"
	CurrencyTickerMXN CurrencyTicker = "MXN"
)

func (e *CurrencyTicker) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CurrencyTicker(s)
	case string:
		*e = CurrencyTicker(s)
	default:
		return fmt.Errorf("unsupported scan type for CurrencyTicker: %T", src)
	}
	return nil
}

type NullCurrencyTicker struct {
	CurrencyTicker CurrencyTicker
	Valid          bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCurrencyTicker) Scan(value interface{}) error {
	if value == nil {
		ns.CurrencyTicker, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CurrencyTicker.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCurrencyTicker) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.CurrencyTicker, nil
}

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)

func (e *TransactionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionType(s)
	case string:
		*e = TransactionType(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionType: %T", src)
	}
	return nil
}

type NullTransactionType struct {
	TransactionType TransactionType
	Valid           bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTransactionType) Scan(value interface{}) error {
	if value == nil {
		ns.TransactionType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TransactionType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTransactionType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.TransactionType, nil
}

type Account struct {
	ID       int32           `db:"id" json:"id"`
	ClientID int32           `db:"client_id" json:"client_id"`
	Balance  decimal.Decimal `db:"balance" json:"balance"`
	Currency CurrencyTicker  `db:"currency" json:"currency"`
}

type Client struct {
	ID    int32  `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
}

type Transaction struct {
	ID            int32           `db:"id" json:"id"`
	Type          TransactionType `db:"type" json:"type"`
	Amount        decimal.Decimal `db:"amount" json:"amount"`
	FromAccountID types.NullInt32 `db:"from_account_id" json:"from_account_id"`
	ToAccountID   types.NullInt32 `db:"to_account_id" json:"to_account_id"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
}
