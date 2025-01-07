package gateway

import "github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"

type AccountGateway interface {
	Save(client *entity.Account) error
	FindByID(id string) (*entity.Account, error)
}
