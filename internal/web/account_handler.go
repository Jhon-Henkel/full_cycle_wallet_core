package web

import (
	"encoding/json"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_account"
	"net/http"
)

type AccountHandler struct {
	CreateAccountUseCase create_account.CreateAccountUseCase
}

func NewAccountHandler(createAccountUseCase create_account.CreateAccountUseCase) *AccountHandler {
	return &AccountHandler{
		CreateAccountUseCase: createAccountUseCase,
	}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var dto create_account.CreateAccountInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	output, err := h.CreateAccountUseCase.Execute(dto)
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
