package create_account

import (
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) Save(client *entity.Account) error {
	args := m.Called(client)
	return args.Error(0)
}

func (m *AccountGatewayMock) UpdateBalance(client *entity.Account) error {
	args := m.Called(client)
	return args.Error(0)
}

func TestCreateAccountUseCaseExecute(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "go@go.com")
	clientMock := &ClientGatewayMock{}
	clientMock.On("Get", client.ID).Return(client, nil)

	accountMock := &AccountGatewayMock{}
	accountMock.On("Save", mock.Anything).Return(nil)

	uc := NewCreateAccountUseCase(accountMock, clientMock)
	inputDTO := CreateAccountInputDTO{ClientID: client.ID}
	outputDTO, err := uc.Execute(inputDTO)

	assert.Nil(t, err)
	assert.NotNil(t, outputDTO)
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Get", 1)
	accountMock.AssertExpectations(t)
	accountMock.AssertNumberOfCalls(t, "Save", 1)
}
