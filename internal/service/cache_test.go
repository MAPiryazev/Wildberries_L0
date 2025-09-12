package service

import (
	"testing"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
)

func TestCachePut(t *testing.T) {
	cache := NewLRUCache(100)

	cache.Put(&models.Order{OrderUID: "orderUID"})

	_, exists := cache.cache["orderUID"]
	if !exists {
		t.Fatalf("заказ не добавился в мапу кэша")
	}

	if cache.list.Front() == nil {
		t.Fatalf("заказ не добавился в список кэша")
	}
}

func TestCacheGet(t *testing.T) {
	cache := NewLRUCache(100)
	cache.Put(&models.Order{OrderUID: "orderUID"})

	got, ok := cache.Get("orderUID")
	if ok == false {
		t.Fatalf("Кэш не нашел существующее значение")
	}
	if got == nil {
		t.Fatalf("Кэш вернул битое значение")
	}
}
