package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/repository"
)

const MAX_CACHE_SIZE = 1000

type OrderService interface {
	GetOrderByID(id string) (*models.Order, error)
	SaveOrder(order *models.Order) error
	SaveOrdersBatch(orders []*models.Order) error
	SendToDLQ(ctx context.Context, original []byte, errMsg string) error
}

type orderService struct {
	repo      repository.OrderRepository
	cache     *LRUCache
	dlqWriter *kafka.Writer
}

var ErrOrderNotFound = errors.New("заказ не найден")

func (orderService *orderService) GetOrderByID(id string) (*models.Order, error) {
	order, ok := orderService.cache.Get(id)
	if ok {
		log.Println("Заказ нашелся в кэше")
		return order, nil
	}

	order, err := orderService.repo.GetOrderById(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	orderService.cache.Put(order)

	return order, nil
}

func (s *orderService) SaveOrder(order *models.Order) error {
	err := s.repo.SaveOrder(order)
	if err != nil {
		dlqErr := s.SendToDLQ(context.Background(), marshalOrder(order), err.Error())
		if dlqErr != nil {
			log.Printf("Ошибка при отправке заказа в DLQ: %v", dlqErr)
		}
		return err
	}

	s.cache.Put(order)
	return nil
}

func (s *orderService) SaveOrdersBatch(orders []*models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	err := s.repo.SaveOrdersBatch(orders)
	if err != nil {
		for _, order := range orders {
			if order == nil {
				continue
			}
			dlqErr := s.SendToDLQ(context.Background(), marshalOrder(order), err.Error())
			if dlqErr != nil {
				log.Printf("Ошибка при отправке заказа в DLQ: %v", dlqErr)
			}
		}
		return err
	}

	for _, order := range orders {
		if order == nil {
			continue
		}
		s.cache.Put(order)
	}

	return nil
}

// вспомогательная функция для сериализации заказа
func marshalOrder(order *models.Order) []byte {
	data, err := json.Marshal(order)
	if err != nil {
		return []byte("cannot marshal order")
	}
	return data
}

func (s *orderService) SendToDLQ(ctx context.Context, original []byte, errMsg string) error {
	dlqMsg := models.DLQMessage{
		OriginalValue: original,
		Error:         errMsg,
		Timestamp:     time.Now(),
	}

	data, err := json.Marshal(dlqMsg)
	if err != nil {
		return err
	}

	err = s.dlqWriter.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
	if err != nil {
		log.Printf("ошибка при записи в DLQ: %v", err)
		return err
	}

	return nil
}

func NewOrderService(repo repository.OrderRepository, preloadCount int, dlqWriter *kafka.Writer) (OrderService, error) {
	LRUCache := NewLRUCache(MAX_CACHE_SIZE)

	if preloadCount > 0 {
		orders, err := repo.GetLastNOrders(preloadCount)
		if err != nil {
			return nil, err
		}
		for _, order := range orders {
			LRUCache.Put(order)
		}
	}

	return &orderService{
		repo:      repo,
		cache:     LRUCache,
		dlqWriter: dlqWriter,
	}, nil
}
