package gateway

import "github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"

type BalanceGateway interface {
	FindLastByAccountID(accountId string) (*entity.Balance, error)
}
