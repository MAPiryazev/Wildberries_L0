package service

import (
	"errors"
	"log"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/repository"
)

const MAX_CACHE_SIZE = 1000

type OrderService interface {
	GetOrderByID(id string) (*models.Order, error)
	SaveOrder(order *models.Order) error
	SaveOrdersBatch(orders []*models.Order) error
}

type orderService struct {
	repo  repository.OrderRepository
	cache *LRUCache
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

func (orderService *orderService) SaveOrder(order *models.Order) error {
	err := orderService.repo.SaveOrder(order)
	if err != nil {
		return err
	}

	orderService.cache.Put(order)

	return nil
}

func (orderService *orderService) SaveOrdersBatch(orders []*models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	if err := orderService.repo.SaveOrdersBatch(orders); err != nil {
		return err
	}

	for _, order := range orders {
		if order == nil {
			continue
		}
		orderService.cache.Put(order)
	}

	return nil
}

func NewOrderService(repo repository.OrderRepository, preloadCount int) (OrderService, error) {
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
		repo:  repo,
		cache: LRUCache,
	}, nil
}
