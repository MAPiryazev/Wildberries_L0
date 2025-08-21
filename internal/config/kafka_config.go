package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type KafkaConfig struct {
	KafkaPort      int
	KafkaHost      string
	KafkaTopicName string
	KafkaGroupID   string
}

func LoadKafkaConfig() *KafkaConfig {
	if err := godotenv.Load(".env"); err != nil {
		if err2 := godotenv.Load("../environment/.env"); err2 != nil {
			log.Fatalln("../environment/.env файл не найден: ", err2)
		}
	}

	kafkaPort, err := strconv.Atoi(os.Getenv("KAFKA_PORT"))
	if err != nil {
		log.Println("Ошибка при считывании порта для кафки, ставим 29092")
		kafkaPort = 29092
	}
	return &KafkaConfig{
		KafkaPort:      kafkaPort,
		KafkaHost:      os.Getenv("KAFKA_HOST"),
		KafkaTopicName: os.Getenv("KAFKA_TOPIC_NAME"),
		KafkaGroupID:   os.Getenv("KAFKA_GROUP_ID"),
	}
}
