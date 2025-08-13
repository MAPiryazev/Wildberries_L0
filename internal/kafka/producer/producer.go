package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker = "localhost:29092" //поменять если будет запускаться внутри контейнера
	topicName   = "orders"
)

func createTopic() error {
	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	// Собираем правильный адрес контроллера
	address := controller.Host + ":" + strconv.Itoa(int(controller.Port))

	connController, err := kafka.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer connController.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topicName,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	return connController.CreateTopics(topicConfigs...)
}

func produceOrder(order *models.Order) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: orderJSON,
		Time:  time.Now(),
	}

	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	log.Println("заказ отправлен в kafka", order.OrderUID)
	return nil
}

func main() {
	if err := createTopic(); err != nil {
		log.Fatalf("Не удалось создать топик: %v", err)
	}

	dir := "./test_jsons"

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".json" {
			log.Println("Processing file:", path)

			data, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("Failed to read file:", err)
				return nil
			}

			var order models.Order
			err = json.Unmarshal(data, &order)
			if err != nil {
				log.Println("Failed to unmarshal JSON:", err)
				return nil
			}

			err = produceOrder(&order)
			if err != nil {
				log.Println("Failed to produce order:", err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal("Error walking the directory:", err)
	}
}
