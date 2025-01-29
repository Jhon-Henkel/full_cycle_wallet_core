package web

import (
	"encoding/json"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/get_balance"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type BalanceHandler struct {
	GetBalanceUseCase get_balance.GetBalanceUseCase
}

func NewBalanceHandler(getBalanceUseCase get_balance.GetBalanceUseCase) *BalanceHandler {
	return &BalanceHandler{
		GetBalanceUseCase: getBalanceUseCase,
	}
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dto := get_balance.GetBalanceInputDTO{AccountID: id}
	output, err := h.GetBalanceUseCase.Execute(dto)
	log.Println("GetBalance: ", output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
