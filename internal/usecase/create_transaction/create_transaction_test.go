package create_transaction

import (
	"context"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/event"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/mocks"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/events"
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

	mockUow := mocks.UowMock{}
	mockUow.On("Do", mock.Anything).Return(nil)

	input := CreateTransactionInputDTO{
		AccountIDFrom: accountOne.ID,
		AccountIDTo:   accountTwo.ID,
		Amount:        100,
	}

	dispatcher := events.NewEventDispatcher()
	eventTransactionCreated := event.NewTransactionCreated()
	eventBalanceUpdated := event.NewBalanceUpdated()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(&mockUow, dispatcher, eventTransactionCreated, eventBalanceUpdated)
	output, err := uc.Execute(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, output)

	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
