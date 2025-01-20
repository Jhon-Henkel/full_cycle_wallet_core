package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/database"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/event"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_account"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_client"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_transaction"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/web"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/web/webserver"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/events"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "localhost", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	//eventDispatcher.Register("TransactionCreated", handler)

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)

	ctx := context.Background()
	unitOfWork := uow.NewUow(ctx, db)
	unitOfWork.Register("accountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	unitOfWork.Register("transactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDB)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDB, clientDB)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(unitOfWork, eventDispatcher, transactionCreatedEvent)

	ws := webserver.NewWebserver(":3000")

	clientHandler := web.NewClientHandler(*createClientUseCase)
	accountHandler := web.NewAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewTransactionHandler(*createTransactionUseCase)

	ws.AddHandler("/client", clientHandler.CreateClient)
	ws.AddHandler("/accounts", accountHandler.CreateAccount)
	ws.AddHandler("/transactions", transactionHandler.CreateTransaction)

	log.Println("Starting web server on port 3000")
	ws.Start()
}
