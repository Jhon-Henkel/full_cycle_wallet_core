package entity

import (
	"github.com/google/uuid"
	"time"
)

type Balance struct {
	ID        string
	AccountID string
	Amount    float64
	CreatedAt time.Time
}

func NewBalance(balance *Balance) *Balance {
	if balance == nil {
		return nil
	}
	return &Balance{
		ID:        uuid.New().String(),
		AccountID: balance.AccountID,
		Amount:    balance.Amount,
		CreatedAt: time.Now(),
	}
}
