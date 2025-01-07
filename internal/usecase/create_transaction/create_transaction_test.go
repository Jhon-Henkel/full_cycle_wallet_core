package create_transaction

import (
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
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

func TestCreateTransactionUseCaseExecute(t *testing.T) {
	clientOne, _ := entity.NewClient("John Doe", "go@go.com")
	accountOne := entity.NewAccount(clientOne)
	accountOne.Credit(1000)

	clientTwo, _ := entity.NewClient("Jane Doe Doe", "go@horse.com")
	accountTwo := entity.NewAccount(clientTwo)
	accountTwo.Credit(1000)

	mockAccount := &AccountGatewayMock{}
	mockAccount.On("FindByID", accountOne.ID).Return(accountOne, nil)
	mockAccount.On("FindByID", accountTwo.ID).Return(accountTwo, nil)

	mockTransaction := &TransactionGatewayMock{}
	mockTransaction.On("Create", mock.Anything).Return(nil)

	input := CreateTransactionInputDTO{
		AccountIDFrom: accountOne.ID,
		AccountIDTo:   accountTwo.ID,
		Amount:        100,
	}

	uc := NewCreateTransactionUseCase(mockTransaction, mockAccount)
	output, err := uc.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)

	mockAccount.AssertExpectations(t)
	mockAccount.AssertNumberOfCalls(t, "FindByID", 2)
	mockTransaction.AssertExpectations(t)
	mockTransaction.AssertNumberOfCalls(t, "Create", 1)
}
