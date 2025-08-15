package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"

	"github.com/MAPiryazev/Wildberries_L0/internal/config"
	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/service"
)

type OrderConsumer struct {
	kafkaBroker string
	topic       string
	groupID     string
	orderSvc    service.OrderService
}

func NewOrderConsumer(svc service.OrderService) *OrderConsumer {
	kafkaConfig := config.LoadKafkaConfig()

	return &OrderConsumer{
		kafkaBroker: kafkaConfig.KafkaHost + ":" + strconv.Itoa(kafkaConfig.KafkaPort),
		topic:       kafkaConfig.KafkaTopicName,
		groupID:     kafkaConfig.KafkaGroupID,
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
