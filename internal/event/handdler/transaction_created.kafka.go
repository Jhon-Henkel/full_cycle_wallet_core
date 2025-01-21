package handdler

import (
	"fmt"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/events"
	"github.com/Jhon-Henkel/full_cycle_wallet_core/pkg/kafka"
	"sync"
)

type TransactionCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewTransactionCreatedKafkaHandler(kafka *kafka.Producer) *TransactionCreatedKafkaHandler {
	return &TransactionCreatedKafkaHandler{
		Kafka: kafka,
	}
}

func (h *TransactionCreatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(message, nil, "transactions")
	fmt.Println("TransactionCreatedKafkaHandler: ", message)
}
