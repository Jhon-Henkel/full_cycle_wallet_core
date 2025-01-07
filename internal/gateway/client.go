package gateway

import "github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
