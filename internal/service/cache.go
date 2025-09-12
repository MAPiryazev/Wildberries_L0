package service

import (
	"container/list"
	"log"
	"sync"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
)

type LRUCache struct {
	mu    sync.Mutex
	max   int
	cache map[string]*list.Element
	list  *list.List
}

type entry struct {
	key   string
	order *models.Order
}

func NewLRUCache(max int) *LRUCache {
	if max < 100 {
		max = 100
		log.Println("значение емкости кэша очень мало, ставим ", max, " по умолчанию")
	}
	return &LRUCache{
		max:   max,
		cache: make(map[string]*list.Element),
		list:  list.New(),
	}
}

func (c *LRUCache) Get(key string) (*models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.cache[key]
	if ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry).order, true
	} else {
		return nil, false
	}
}

func (c *LRUCache) Put(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.cache[order.OrderUID]
	if ok {
		c.list.MoveToFront(elem)
		elem.Value.(*entry).order = order
		return
	}
	if c.list.Len() >= c.max {
		oldest := c.list.Back()
		if oldest != nil {
			c.list.Remove(oldest)
			delete(c.cache, oldest.Value.(*entry).key)
		}
	}
	elem = c.list.PushFront(&entry{key: order.OrderUID, order: order})
	c.cache[order.OrderUID] = elem
}
