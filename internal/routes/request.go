package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/seth2810/money-processing/internal/storage/queries"
	"github.com/shopspring/decimal"
)

type CreateClientRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type CreateAccountRequest struct {
	ClientID int32                  `json:"client_id" validate:"required,min=1"` //nolint:tagliatelle
	Currency queries.CurrencyTicker `json:"currency" validate:"required,oneof=USD COP MXN"`
}

type CreateTransactionRequest struct {
	Type          queries.TransactionType `json:"type" validate:"required"`
	Amount        decimal.Decimal         `json:"amount" validate:"required"`
	FromAccountID int32                   `json:"from_account_id" validate:"omitempty,min=1"` //nolint:tagliatelle
	ToAccountID   int32                   `json:"to_account_id" validate:"omitempty,min=1"`   //nolint:tagliatelle
}

func parseBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func validateBody(v interface{}) error {
	validate := validator.New()

	return validate.Struct(v)
}
