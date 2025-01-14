package web

import (
	"encoding/json"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_client"
	"net/http"
)

type ClientHandler struct {
	CreateClientUseCase create_client.CreateClientUseCase
}

func NewClientHandler(createClientUseCase create_client.CreateClientUseCase) *ClientHandler {
	return &ClientHandler{
		CreateClientUseCase: createClientUseCase,
	}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var dto create_client.CreateClientInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.CreateClientUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
