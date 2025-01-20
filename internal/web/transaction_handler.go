package web

import (
	"encoding/json"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_transaction"
	"net/http"
)

type TransactionHandler struct {
	CreateTransactionUseCase create_transaction.CreateTransactionUseCase
}

func NewTransactionHandler(createTransactionUseCase create_transaction.CreateTransactionUseCase) *TransactionHandler {
	return &TransactionHandler{
		CreateTransactionUseCase: createTransactionUseCase,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var dto create_transaction.CreateTransactionInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	output, err := h.CreateTransactionUseCase.Execute(ctx, dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
