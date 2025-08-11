package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/segmentio/kafka-go"
)

func produceOrder(order *models.Order) error {
	//создаем писателя в kafka (подключение к брокеру)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return err
	}

	//создаем сообщение
	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: orderJSON,
		Time:  time.Now(),
	}

	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	log.Println("заказ отправлен в kafka ", order.OrderUID)
	return nil
}

func main() {
	dir := "./test_jsons"

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Пропускаем папки, читаем только файлы с расширением .json
		if !d.IsDir() && filepath.Ext(path) == ".json" {
			log.Println("Processing file:", path)

			// Читаем содержимое файла
			data, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println("Failed to read file:", err)
				return nil // не останавливаем весь процесс из-за ошибки одного файла
			}

			// Парсим JSON в структуру заказа
			var order models.Order
			err = json.Unmarshal(data, &order)
			if err != nil {
				log.Println("Failed to unmarshal JSON:", err)
				return nil
			}

			// Отправляем заказ в Kafka
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
