package kafka

import (
	"context"
	"encoding/json"
	"log"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/service"
	"github.com/segmentio/kafka-go"
)

type OrderConsumer struct {
	kafkaBroker string
	topic       string
	groupID     string
	orderSvc    service.OrderService
}

func NewOrderConsumer(broker, topic, groupID string, svc service.OrderService) *OrderConsumer {
	return &OrderConsumer{
		kafkaBroker: broker,
		topic:       topic,
		groupID:     groupID,
		orderSvc:    svc,
	}
}

func (consumer *OrderConsumer) Start(ctx context.Context) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{consumer.kafkaBroker},
		GroupID: consumer.groupID,
		Topic:   consumer.topic,
	})
	defer reader.Close()

	log.Println("Консъюмер кафки запущен")

	for {
		select {
		case <-ctx.Done():
			log.Println("Консьюмер кафки остановлен")
			return
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Ошибка чтения сообщения: %v", err)
				continue
			}

			var order models.Order
			if err := json.Unmarshal(message.Value, &order); err != nil {
				log.Printf("Ошибка десериализации json сообщения: %v", err)
				continue
			}

			if err := consumer.orderSvc.SaveOrder(&order); err != nil {
				log.Printf("Ошибка при сохранении заказа в БД %s: %v", order.OrderUID, err)
				continue
			}

			log.Printf("Заказ %s успешно сохранен в БД", order.OrderUID)
		}
	}
}
