package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/seth2810/money-processing/internal/storage"
)

func NewRouter(s *storage.Storage) *mux.Router {
	r := mux.NewRouter()
	h := &handler{storage: s}

	r.HandleFunc("/clients", h.CreateClient).Methods(http.MethodPost)
	r.HandleFunc("/clients/{id:[0-9]+}", h.GetClient).Methods(http.MethodGet)

	r.HandleFunc("/accounts", h.CreateAccount).Methods(http.MethodPost)
	r.HandleFunc("/accounts/{id:[0-9]+}", h.GetAccount).Methods(http.MethodGet)
	r.HandleFunc("/accounts/{id:[0-9]+}/transactions", h.GetTransactions).Methods(http.MethodGet)

	r.HandleFunc("/transactions", h.CreateTransaction).Methods(http.MethodPost)

	return r
}
