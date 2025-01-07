package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	client, _ := NewClient("John Doe", "go@go.com")
	account := NewAccount(client)
	assert.NotNil(t, account)
	assert.Equal(t, account.Client.ID, client.ID)
}

func TestCreateAccountWithNilClient(t *testing.T) {
	account := NewAccount(nil)
	assert.Nil(t, account)
}

func TestCreditAccount(t *testing.T) {
	client, _ := NewClient("John Doe", "go@go.com")
	account := NewAccount(client)
	account.Credit(100)
	assert.Equal(t, float64(100), account.Balance)
}

func TestDebitAccount(t *testing.T) {
	client, _ := NewClient("John Doe", "go@go.com")
	account := NewAccount(client)
	account.Credit(100)
	account.Debit(50)
	assert.Equal(t, float64(50), account.Balance)
}

func TestAddAccountToClient(t *testing.T) {
	client, _ := NewClient("John Doe", "go@go.com")
	err := client.AddAccount(NewAccount(client))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(client.Accounts))
}
