package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	v "github.com/gobuffalo/validate"
	"github.com/seth2810/money-processing/internal/storage/queries"
	"github.com/shopspring/decimal"
)

type CreateClientRequest struct {
	Email string `json:"email"`
}

func (r *CreateClientRequest) IsValid(errors *v.Errors) {
	if r.Email == "" {
		errors.Add("email", "should not be blank")
	}
}

type CreateAccountRequest struct {
	ClientID int32                  `json:"client_id"` //nolint:tagliatelle
	Currency queries.CurrencyTicker `json:"currency"`
}

func (r *CreateAccountRequest) IsValid(errors *v.Errors) {
	if r.ClientID == 0 {
		errors.Add("client_id", "should not be blank")
	}

	currencies := []queries.CurrencyTicker{
		queries.CurrencyTickerUSD,
		queries.CurrencyTickerCOP,
		queries.CurrencyTickerMXN,
	}

	var isCurrencyValid bool

	for _, v := range currencies {
		if r.Currency == v {
			isCurrencyValid = true
			break
		}
	}

	if !isCurrencyValid {
		errors.Add("currency", fmt.Sprintf("should be one of: %s", currencies))
	}
}

type CreateTransactionRequest struct {
	Type          queries.TransactionType `json:"type"`
	Amount        decimal.Decimal         `json:"amount"`
	FromAccountID int32                   `json:"from_account_id"` //nolint:tagliatelle
	ToAccountID   int32                   `json:"to_account_id"`   //nolint:tagliatelle
}

func (r *CreateTransactionRequest) IsValid(errors *v.Errors) {
	types := []queries.TransactionType{
		queries.TransactionTypeDeposit,
		queries.TransactionTypeTransfer,
		queries.TransactionTypeWithdraw,
	}

	var isTypeValid bool

	for _, v := range types {
		if r.Type == v {
			isTypeValid = true
			break
		}
	}

	if !isTypeValid {
		errors.Add("type", fmt.Sprintf("should be one of: %s", types))
	}

	if r.Amount.LessThanOrEqual(decimal.NewFromInt(0)) {
		errors.Add("amount", "should not be negative")
	}

	if r.FromAccountID == 0 {
		switch r.Type {
		case queries.TransactionTypeWithdraw, queries.TransactionTypeTransfer:
			errors.Add("from_account_id", "should not be blank")
		case queries.TransactionTypeDeposit:
			// skip
		}
	}

	if r.ToAccountID == 0 {
		switch r.Type {
		case queries.TransactionTypeDeposit, queries.TransactionTypeTransfer:
			errors.Add("to_account_id", "should not be blank")
		case queries.TransactionTypeWithdraw:
			// skip
		}
	}

	if r.Type == queries.TransactionTypeTransfer {
		if r.FromAccountID == r.ToAccountID {
			errors.Add("from_account_id", "should not be equal to_account_id")
			errors.Add("to_account_id", "should not be equal from_account_id")
		}
	}
}

func parseBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
