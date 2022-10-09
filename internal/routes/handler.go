package routes

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/seth2810/money-processing/internal/storage"
	"github.com/seth2810/money-processing/internal/storage/queries"
)

var (
	errTransactionTypeNotSupported = errors.New("transaction type not supported")
	errTransferDifferentCurrencies = errors.New("transfer between accounts with different currencies not supported")
)

type handler struct {
	storage *storage.Storage
}

func (h *handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var body CreateClientRequest

	if err := parseBody(r, &body); err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := validateBody(&body); err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	client, err := h.storage.CreateClient(r.Context(), body.Email)
	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	replyWithJSON(w, http.StatusOK, client)
}

func (h *handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body CreateAccountRequest

	if err := parseBody(r, &body); err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := validateBody(&body); err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	acc, err := h.storage.CreateAccount(r.Context(), body.ClientID, body.Currency)
	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	replyWithJSON(w, http.StatusOK, acc)
}

func (h *handler) GetClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	client, err := h.storage.GetClient(r.Context(), int32(id))
	if errors.Is(err, storage.ErrClientNotFound) {
		replyWithError(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	replyWithJSON(w, http.StatusOK, client)
}

func (h *handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	acc, err := h.storage.GetAccount(r.Context(), int32(id))
	if errors.Is(err, storage.ErrAccountNotFound) {
		replyWithError(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	replyWithJSON(w, http.StatusOK, acc)
}

func (h *handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	acc, err := h.storage.GetAccount(r.Context(), int32(id))
	if errors.Is(err, storage.ErrAccountNotFound) {
		replyWithError(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	txs, err := h.storage.GetTransactions(r.Context(), acc.ID)
	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	replyWithJSON(w, http.StatusOK, txs)
}

func (h *handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var body CreateTransactionRequest

	if err := parseBody(r, &body); err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	if err := validateBody(&body); err != nil {
		replyWithError(w, http.StatusBadRequest, err)
		return
	}

	switch body.Type {
	case queries.TransactionTypeDeposit:
		h.deposit(r.Context(), &body, w)
	case queries.TransactionTypeWithdraw:
		h.withdraw(r.Context(), &body, w)
	case queries.TransactionTypeTransfer:
		h.transfer(r.Context(), &body, w)
	default:
		replyWithError(w, http.StatusBadRequest, errTransactionTypeNotSupported)
		return
	}
}

func (h *handler) deposit(ctx context.Context, body *CreateTransactionRequest, w http.ResponseWriter) {
	toAccount, err := h.storage.GetAccount(ctx, body.ToAccountID)
	if errors.Is(err, storage.ErrClientNotFound) {
		replyWithError(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	tx, err := h.storage.Deposit(ctx, body.Amount, toAccount.ID)
	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	replyWithJSON(w, http.StatusOK, tx)
}

func (h *handler) withdraw(ctx context.Context, body *CreateTransactionRequest, w http.ResponseWriter) {
	fromAccount, err := h.storage.GetAccount(ctx, body.FromAccountID)
	if errors.Is(err, storage.ErrClientNotFound) {
		replyWithError(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	tx, err := h.storage.Withdraw(ctx, body.Amount, fromAccount.ID)
	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	replyWithJSON(w, http.StatusOK, tx)
}

func (h *handler) transfer(ctx context.Context, body *CreateTransactionRequest, w http.ResponseWriter) {
	fromAccount, err := h.storage.GetAccount(ctx, body.FromAccountID)
	if errors.Is(err, storage.ErrClientNotFound) {
		replyWithError(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	toAccount, err := h.storage.GetAccount(ctx, body.ToAccountID)
	if errors.Is(err, storage.ErrClientNotFound) {
		replyWithError(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	if fromAccount.Currency != toAccount.Currency {
		replyWithError(w, http.StatusBadRequest, errTransferDifferentCurrencies)
		return
	}

	tx, err := h.storage.Transfer(ctx, body.Amount, body.FromAccountID, body.ToAccountID)
	if err != nil {
		replyWithError(w, http.StatusInternalServerError, err)
		return
	}

	replyWithJSON(w, http.StatusOK, tx)
}
