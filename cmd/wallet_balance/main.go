package main

import (
	"database/sql"
	"fmt"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{"bootstrap.servers": "kafka:29092"}
	kafkaProducer := kafka.NewConsumer(&configMap, []string{"balances"})

}
