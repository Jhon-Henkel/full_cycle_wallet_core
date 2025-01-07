package gateway

import "github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
