package database

import (
	"sync"

	"github.com/3-lines-studio/datafrost/internal/core/entity"
)

type AdapterCache struct {
	mu      sync.Mutex
	entries map[int64]entity.DatabaseAdapter
	factory *Factory
}

func NewAdapterCache() *AdapterCache {
	return &AdapterCache{
		entries: make(map[int64]entity.DatabaseAdapter),
		factory: NewFactory(),
	}
}

func (c *AdapterCache) Get(id int64, connType string, credentials map[string]any) (entity.DatabaseAdapter, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if adapter, ok := c.entries[id]; ok {
		return adapter, nil
	}

	adapter, err := c.factory.GetAdapter(connType)
	if err != nil {
		return nil, err
	}

	if err := adapter.Connect(credentials); err != nil {
		return nil, err
	}

	c.entries[id] = adapter
	return adapter, nil
}

func (c *AdapterCache) Invalidate(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if adapter, ok := c.entries[id]; ok {
		adapter.Close()
		delete(c.entries, id)
	}
}

func (c *AdapterCache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, adapter := range c.entries {
		adapter.Close()
		delete(c.entries, id)
	}
}
