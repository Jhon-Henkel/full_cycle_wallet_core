package create_account

import (
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/gateway"
)

type CreateAccountInputDTO struct {
	ClientID string
}

type CreateAccountOutputDTO struct {
	ClientID string
}

type CreateAccountUseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway  gateway.ClientGateway
}

func NewCreateAccountUseCase(accountGateway gateway.AccountGateway, clientGateway gateway.ClientGateway) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountGateway: accountGateway,
		ClientGateway:  clientGateway,
	}
}

func (uc *CreateAccountUseCase) Execute(inputDTO *CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	client, err := uc.ClientGateway.Get(inputDTO.ClientID)
	if err != nil {
		return nil, err
	}
	account := entity.NewAccount(client)
	err = uc.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}
	return &CreateAccountOutputDTO{
		ClientID: account.ID,
	}, nil
}
