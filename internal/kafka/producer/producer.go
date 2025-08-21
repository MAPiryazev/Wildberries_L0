package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker = "localhost:29092"
	topicName   = "orders"
	batchSize   = 10000
)

func main() {
	file, err := os.Open("test_jsons/orders.json")
	if err != nil {
		log.Fatal("Ошибка открытия файла:", err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	_, err = dec.Token()
	if err != nil {
		log.Fatal("Ошибка чтения токена начала массива:", err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{kafkaBroker},
		Topic:     topicName,
		Balancer:  &kafka.LeastBytes{},
		BatchSize: batchSize,
	})
	defer writer.Close()

	var batch []kafka.Message
	count := 0
	start := time.Now()

	for dec.More() {
		var order models.Order
		if err := dec.Decode(&order); err != nil {
			log.Fatal("Ошибка декодирования JSON:", err)
		}

		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Println("Ошибка маршалинга:", err)
			continue
		}

		batch = append(batch, kafka.Message{
			Key:   []byte(order.OrderUID),
			Value: orderJSON,
			Time:  time.Now(),
		})

		if len(batch) >= batchSize {
			if err := writer.WriteMessages(context.Background(), batch...); err != nil {
				log.Println("Ошибка отправки батча в Kafka:", err)
			}
			count += len(batch)
			log.Printf("Отправлено %d сообщений (%.2f сек)", count, time.Since(start).Seconds())
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := writer.WriteMessages(context.Background(), batch...); err != nil {
			log.Println("Ошибка отправки последнего батча:", err)
		}
		count += len(batch)
	}

	_, err = dec.Token()
	if err != nil {
		log.Fatal("Ошибка чтения конца массива:", err)
	}

	log.Printf("Готово: отправлено %d заказов за %.2f сек", count, time.Since(start).Seconds())
}
