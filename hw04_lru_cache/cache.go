package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	mu       sync.RWMutex
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	queue := c.queue
	items := c.items

	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := items[key]; ok {
		item.Value = value
		queue.MoveToFront(item)
		return true
	}
	item := queue.PushFront(value)
	item.Key = key
	items[key] = item
	if queue.Len() > c.capacity {
		lastItem := queue.Back()
		queue.Remove(lastItem)
		delete(items, lastItem.Key)
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	queue := c.queue
	items := c.items

	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := items[key]; ok {
		queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
