package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/database"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/usecase/get_balance"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/internal/web"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type BalanceUpdatedPayload struct {
	AccountIDFrom        string `json:"account_id_from"`
	AccountIDTo          string `json:"account_id_to"`
	BalanceAccountIDFrom int    `json:"balance_account_id_from"`
	BalanceAccountIDTo   int    `json:"balance_account_id_to"`
}

type BalanceUpdated struct {
	Name    string                `json:"Name"`
	Payload BalanceUpdatedPayload `json:"Payload"`
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaConsumer := kafka.NewConsumer(&configMap, []string{"balances"})
	msgChan := make(chan *ckafka.Message)

	go func() {
		err := kafkaConsumer.Consume(msgChan)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for msg := range msgChan {
			var balanceUpdated BalanceUpdated
			err := json.Unmarshal(msg.Value, &balanceUpdated)
			if err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			_, err = db.Exec("INSERT INTO balances (id, account_id, amount, created_at) VALUES (?, ?, ?, ?)",
				uuid.New().String(),
				balanceUpdated.Payload.AccountIDFrom,
				balanceUpdated.Payload.BalanceAccountIDFrom,
				time.Now(),
			)

			_, err = db.Exec("INSERT INTO balances (id, account_id, amount, created_at) VALUES (?, ?, ?, ?)",
				uuid.New().String(),
				balanceUpdated.Payload.AccountIDTo,
				balanceUpdated.Payload.BalanceAccountIDTo,
				time.Now(),
			)

			if err != nil {
				log.Printf("[error] - insert balance in database: %v", err)
			} else {
				log.Printf("[success] - Inserted balance into database")
			}
		}
	}()

	balanceDB := database.NewBalanceDB(db)
	getBalanceUseCase := get_balance.NewGetBalanceUseCase(balanceDB)
	balanceHandler := web.NewBalanceHandler(*getBalanceUseCase)

	r := mux.NewRouter()
	r.HandleFunc("/balances/{id}", balanceHandler.GetBalance).Methods("GET")

	http.Handle("/", r)
	log.Println("Starting web server on port 3003")
	log.Fatal(http.ListenAndServe(":3003", nil))
}
