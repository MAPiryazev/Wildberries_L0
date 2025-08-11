package service

import (
	"errors"
	"sync"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/repository"
)

type OrderService interface {
	GetOrderByID(id string) (*models.Order, error)
	SaveOrder(order *models.Order) error
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

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{
		repo:  repo,
		cache: make(map[string]*models.Order),
	}
}
