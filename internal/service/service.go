package service

import (
	"errors"
	"log"
	"sync"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/repository"
)

type OrderService interface {
	GetOrderByID(id string) (*models.Order, error)
	SaveOrder(order *models.Order) error
	SaveOrdersBatch(orders []*models.Order) error
}

type orderService struct {
	repo  repository.OrderRepository
	cache map[string]*models.Order
	mu    sync.RWMutex
}

func (orderService *orderService) GetOrderByID(id string) (*models.Order, error) {
	orderService.mu.RLock()
	if order, ok := orderService.cache[id]; ok {
		orderService.mu.RUnlock()
		log.Println("Заказ ", order.OrderUID, " нашелся в кэше")
		return order, nil
	}
	orderService.mu.RUnlock()

	order, err := orderService.repo.GetOrderById(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	orderService.mu.Lock()
	orderService.cache[id] = order
	orderService.mu.Unlock()

	return order, nil
}

func (orderService *orderService) SaveOrder(order *models.Order) error {
	err := orderService.repo.SaveOrder(order)
	if err != nil {
		return err
	}
	orderService.mu.Lock()
	defer orderService.mu.Unlock()
	orderService.cache[order.OrderUID] = order

	return nil
}

func (orderService *orderService) SaveOrdersBatch(orders []*models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	if err := orderService.repo.SaveOrdersBatch(orders); err != nil {
		return err
	}

	orderService.mu.Lock()
	for _, order := range orders {
		if order == nil {
			continue
		}
		orderService.cache[order.OrderUID] = order
	}
	orderService.mu.Unlock()

	return nil
}

func NewOrderService(repo repository.OrderRepository, preloadCount int) (OrderService, error) {
	cache := make(map[string]*models.Order)

	if preloadCount > 0 {
		orders, err := repo.GetLastNOrders(preloadCount)
		if err != nil {
			return nil, err
		}
		for _, order := range orders {
			cache[order.OrderUID] = order
		}
	}

	return &orderService{
		repo:  repo,
		cache: cache,
	}, nil
}
