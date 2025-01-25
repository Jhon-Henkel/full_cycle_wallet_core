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
