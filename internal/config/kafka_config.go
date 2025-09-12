package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ErrKafkaParamNotFound = errors.New(`Один или несолько из критически важных параметров не был найден в env, проверьте:
	KAFKA_HOST
	KAFKA_TOPIC_NAME
	KAFKA_GROUP_ID
	`)
)

type KafkaConfig struct {
	KafkaPort         int
	KafkaHost         string
	KafkaTopicName    string
	KafkaGroupID      string
	KafkaTopicDLQName string
}

func LoadKafkaConfig() (*KafkaConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		if err2 := godotenv.Load("../environment/.env"); err2 != nil {
			return nil, fmt.Errorf("%w: %v", ErrEnvNotFound, err2)
		}
	}

	kafkaPort, err := strconv.Atoi(os.Getenv("KAFKA_PORT"))
	if err != nil {
		kafkaPort = 29092
		log.Println("Не удалось считать порт для кафки, выбран ", kafkaPort)
	}

	KafkaHost := os.Getenv("KAFKA_HOST")
	KafkaTopicName := os.Getenv("KAFKA_TOPIC_NAME")
	KafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	KafkaTopicDLQName := os.Getenv("KAFKA_DLQ_TOPIC_NAME")

	if KafkaHost == "" || KafkaTopicName == "" || KafkaGroupID == "" || KafkaTopicDLQName == "" {
		return nil, fmt.Errorf("%w ", ErrKafkaParamNotFound)
	}

	return &KafkaConfig{
		KafkaPort:         kafkaPort,
		KafkaHost:         KafkaHost,
		KafkaTopicName:    KafkaTopicName,
		KafkaGroupID:      KafkaGroupID,
		KafkaTopicDLQName: KafkaTopicDLQName,
	}, nil
}
