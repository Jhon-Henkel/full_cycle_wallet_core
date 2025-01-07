package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	client, err := NewClient("John Doe", "go@go.com")
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "John Doe", client.Name)
	assert.Equal(t, "go@go.com", client.Email)
	assert.NotEmpty(t, client.ID)
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	client, err := NewClient("", "")
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("John Doe", "go@go.com")
	err := client.Update("Jane Doe Updated", "go@horse.com")
	assert.Nil(t, err)
	assert.Equal(t, "Jane Doe Updated", client.Name)
	assert.Equal(t, "go@horse.com", client.Email)
}

func TestUpdateClientWhenArgsAreInvalid(t *testing.T) {
	client, _ := NewClient("John Doe", "test@test.com")
	err := client.Update("", "")
	assert.Error(t, err, "email is required")
}
