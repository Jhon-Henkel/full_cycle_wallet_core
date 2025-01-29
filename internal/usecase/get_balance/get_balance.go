package get_balance

import (
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/gateway"
)

type GetBalanceInputDTO struct {
	AccountID string
}

type GetBalanceOutputDTO struct {
	AccountID string
	Amount    float64
}

type GetBalanceUseCase struct {
	BalanceGateway gateway.BalanceGateway
}

func NewGetBalanceUseCase(balanceGateway gateway.BalanceGateway) *GetBalanceUseCase {
	return &GetBalanceUseCase{BalanceGateway: balanceGateway}
}

func (uc *GetBalanceUseCase) Execute(inputDTO GetBalanceInputDTO) (*GetBalanceOutputDTO, error) {
	balance, err := uc.BalanceGateway.FindLastByAccountID(inputDTO.AccountID)
	if err != nil {
		return nil, err
	}
	balanceEntity := entity.NewBalance(balance)
	return &GetBalanceOutputDTO{
		AccountID: balanceEntity.AccountID,
		Amount:    balanceEntity.Amount,
	}, nil
}
