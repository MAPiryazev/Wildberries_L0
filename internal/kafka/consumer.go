package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"

	"github.com/MAPiryazev/Wildberries_L0/internal/config"
	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/service"
	"github.com/MAPiryazev/Wildberries_L0/internal/validation"
)

type OrderConsumer struct {
	kafkaBroker string
	topic       string
	DLQTopic    string
	groupID     string
	orderSvc    service.OrderService
	batchSize   int
}

func NewOrderConsumer(svc service.OrderService) (*OrderConsumer, error) {
	kafkaConfig, err := config.LoadKafkaConfig()
	if err != nil {
		//в данном случае обе ошибки которые могут прилететь из пакета config критические,
		// мы просто должны передать любую из них наверх
		return nil, err
	}

	return &OrderConsumer{
		kafkaBroker: kafkaConfig.KafkaHost + ":" + strconv.Itoa(kafkaConfig.KafkaPort),
		topic:       kafkaConfig.KafkaTopicName,
		DLQTopic:    kafkaConfig.KafkaTopicDLQName,
		groupID:     kafkaConfig.KafkaGroupID,
		orderSvc:    svc,
		batchSize:   5000,
	}, nil
}

func (consumer *OrderConsumer) ensureTopicExists(topic string) error {
	conn, err := kafka.Dial("tcp", consumer.kafkaBroker)
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}
	for _, p := range partitions {
		if p.Topic == topic {
			return nil
		}
	}
	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}

func (consumer *OrderConsumer) fetchMessageBatch(ctx context.Context, reader *kafka.Reader, batchSize int) ([]kafka.Message, error) {
	messages := make([]kafka.Message, 0, batchSize)

	for i := 0; i < batchSize; i++ {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				break
			}
			return messages, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (consumer *OrderConsumer) Start(ctx context.Context) {
	if err := consumer.ensureTopicExists(consumer.topic); err != nil {
		log.Fatalf("Не удалось создать топик %s: %v", consumer.topic, err)
	}

	if err := consumer.ensureTopicExists(consumer.DLQTopic); err != nil {
		log.Fatalf("Не удалось создать топик DLQ %s: %v", consumer.DLQTopic, err)
	}

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
			msgs, err := consumer.fetchMessageBatch(ctx, reader, consumer.batchSize)
			if err != nil {
				log.Printf("Ошибка получения батча: %v", err)
				continue
			}

			if len(msgs) == 0 {
				continue
			}

			var orders []*models.Order
			for _, msg := range msgs {
				var order models.Order
				if err := json.Unmarshal(msg.Value, &order); err != nil {
					log.Printf("Ошибка десериализации json сообщения: %v", err)
					_ = consumer.orderSvc.SendToDLQ(ctx, msg.Value, err.Error())
					continue
				}

				if err := validation.ValidateOrder(&order); err != nil {
					log.Printf("Ошибка валидации заказа: %v", err)
					_ = consumer.orderSvc.SendToDLQ(ctx, msg.Value, err.Error())
					continue
				}

				orders = append(orders, &order)
			}

			if len(orders) > 0 {
				if err := consumer.orderSvc.SaveOrdersBatch(orders); err != nil {
					log.Printf("Ошибка при батчевом сохранении %d заказов: %v", len(orders), err)
					continue
				}
				log.Printf("Успешно сохранено %d заказов в БД", len(orders))
			}

			if err := reader.CommitMessages(ctx, msgs...); err != nil {
				log.Printf("Ошибка коммита сообщений: %v", err)
			}
		}
	}
}
