package database

import (
	"database/sql"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/entity"
)

type BalanceDB struct {
	DB *sql.DB
}

func NewBalanceDB(db *sql.DB) *BalanceDB {
	return &BalanceDB{DB: db}
}

func (a *BalanceDB) FindLastByAccountID(id string) (*entity.Balance, error) {
	var balance entity.Balance

	stmt, err := a.DB.Prepare("SELECT id, account_id, amount, created_at FROM balances WHERE account_id = ? ORDER BY created_at DESC LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(&balance.ID, &balance.AccountID, &balance.Amount, &balance.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}
