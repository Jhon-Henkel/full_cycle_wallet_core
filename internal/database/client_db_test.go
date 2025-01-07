package database

import (
	"database/sql"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	clientDB *ClientDB
}

func (s *ClientDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	s.clientDB = NewClientDB(db)
}

func (s *ClientDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (s *ClientDBTestSuite) TestGet() {
	client, _ := entity.NewClient("John Doe", "go@go.com")
	s.clientDB.Save(client)

	clientGet, err := s.clientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, clientGet.ID)
	s.Equal(client.Name, clientGet.Name)
	s.Equal(client.Email, clientGet.Email)
}

func (s *ClientDBTestSuite) TestSave() {
	client, _ := entity.NewClient("John Doe", "go@go.com")
	err := s.clientDB.Save(client)
	s.Nil(err)
}
