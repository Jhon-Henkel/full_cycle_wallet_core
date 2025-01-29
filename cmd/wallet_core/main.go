package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/database"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/event"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/event/handdler"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_account"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_client"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/create_transaction"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/web"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/web/webserver"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/events"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/kafka"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec("CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("CREATE TABLE IF NOT EXISTS accounts (id varchar(255), client_id varchar(255), balance float, created_at date)")
	db.Exec("CREATE TABLE IF NOT EXISTS transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)")
	db.Exec("CREATE TABLE IF NOT EXISTS balances (id varchar(255), account_id varchar(255), amount float, created_at datetime)")

	db.Exec("TRUNCATE clients")
	db.Exec("TRUNCATE accounts")
	db.Exec("TRUNCATE transactions")
	db.Exec("TRUNCATE balances")

	db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES ('043e233b-13ee-44de-9364-bd4df35439e9', 'John Doe', 'go@go.com', '2025-01-01')")
	db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES ('b0ea1f9a-b671-4862-97bd-1bb433f7412a', 'Gina', 'go-horse@go-horse.com', '2025-01-01')")

	db.Exec("INSERT INTO accounts (id, client_id, balance, created_at) VALUES ('569dbb3d-5bd2-44f0-89e7-5f8d80738491', '043e233b-13ee-44de-9364-bd4df35439e9', 900, '2025-01-01')")
	db.Exec("INSERT INTO accounts (id, client_id, balance, created_at) VALUES ('927c20ba-2e83-44bc-aa32-32fd594ff61d', 'b0ea1f9a-b671-4862-97bd-1bb433f7412a', 1100, '2025-01-01')")

	db.Exec("INSERT INTO transactions (id, account_id_from, account_id_to, amount, created_at) VALUES ('9ea1bfbe-9ffb-4506-a211-1d56d1bda103', '569dbb3d-5bd2-44f0-89e7-5f8d80738491', '927c20ba-2e83-44bc-aa32-32fd594ff61d', 100, '2025-01-01')")

	db.Exec("INSERT INTO balances (id, account_id, amount, created_at) VALUES ('3e53df66-ddbe-11ef-b099-0242ac120003', '569dbb3d-5bd2-44f0-89e7-5f8d80738491', 900, '2025-01-01 00:00:00')")
	db.Exec("INSERT INTO balances (id, account_id, amount, created_at) VALUES ('927c20ba-2e83-44bc-aa32-32fd594ff61e', '927c20ba-2e83-44bc-aa32-32fd594ff61d', 1100, '2025-01-01 00:00:00')")

	configMap := ckafka.ConfigMap{"bootstrap.servers": "kafka:29092"}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handdler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handdler.NewBalanceUpdatedKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})
	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDB)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDB, clientDB)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	ws := webserver.NewWebserver(":8080")

	clientHandler := web.NewClientHandler(*createClientUseCase)
	accountHandler := web.NewAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewTransactionHandler(*createTransactionUseCase)

	ws.AddHandler("/client", clientHandler.CreateClient)
	ws.AddHandler("/accounts", accountHandler.CreateAccount)
	ws.AddHandler("/transactions", transactionHandler.CreateTransaction)

	log.Println("Starting web server on port 8080")
	ws.Start()
}
