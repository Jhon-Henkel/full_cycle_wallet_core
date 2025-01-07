package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	clientOne, _ := NewClient("John Doe", "go@go.com")
	accountOne := NewAccount(clientOne)
	accountOne.Credit(1000)

	clientTwo, _ := NewClient("Jane Doe Doe", "go@horse.com")
	accountTwo := NewAccount(clientTwo)
	accountTwo.Credit(1000)

	transaction, err := NewTransaction(accountOne, accountTwo, 100)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(900), accountOne.Balance)
	assert.Equal(t, float64(1100), accountTwo.Balance)
}

func TestCreateTransactionWithInsufficientBalance(t *testing.T) {
	clientOne, _ := NewClient("John Doe", "go@go.com")
	accountOne := NewAccount(clientOne)
	accountOne.Credit(1000)

	clientTwo, _ := NewClient("Jane Doe Doe", "go@horse.com")
	accountTwo := NewAccount(clientTwo)
	accountTwo.Credit(1000)

	transaction, err := NewTransaction(accountOne, accountTwo, 2000)

	assert.NotNil(t, err)
	assert.Error(t, err, "insufficient funds")
	assert.Nil(t, transaction)
	assert.Equal(t, float64(1000), accountOne.Balance)
	assert.Equal(t, float64(1000), accountTwo.Balance)
}

func TestCreateTransactionWithNegativeValue(t *testing.T) {
	clientOne, _ := NewClient("John Doe", "go@go.com")
	accountOne := NewAccount(clientOne)
	accountOne.Credit(1000)

	clientTwo, _ := NewClient("Jane Doe Doe", "go@horse.com")
	accountTwo := NewAccount(clientTwo)
	accountTwo.Credit(1000)

	transaction, err := NewTransaction(accountOne, accountTwo, -500)

	assert.NotNil(t, err)
	assert.Error(t, err, "amount must be greater than 0")
	assert.Nil(t, transaction)
	assert.Equal(t, float64(1000), accountOne.Balance)
	assert.Equal(t, float64(1000), accountTwo.Balance)
}
