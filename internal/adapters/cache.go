package adapters

import (
	"sync"

	"github.com/3-lines-studio/datafrost/internal/models"
)

// AdapterCache maintains a live connection per saved connection ID,
// avoiding the cost of reconnecting on every request.
type AdapterCache struct {
	mu      sync.Mutex
	entries map[int64]models.DatabaseAdapter
	factory *Factory
}

func NewAdapterCache() *AdapterCache {
	return &AdapterCache{
		entries: make(map[int64]models.DatabaseAdapter),
		factory: NewFactory(),
	}
}

// Get returns a cached adapter for the given connection ID, creating and
// connecting one if it doesn't exist yet.
func (c *AdapterCache) Get(id int64, connType string, credentials map[string]any) (models.DatabaseAdapter, error) {
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

// Invalidate closes and removes the cached adapter for the given ID.
// Call this when a connection is deleted or its credentials are updated.
func (c *AdapterCache) Invalidate(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if adapter, ok := c.entries[id]; ok {
		adapter.Close()
		delete(c.entries, id)
	}
}

// Close closes all cached adapters.
func (c *AdapterCache) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, adapter := range c.entries {
		adapter.Close()
		delete(c.entries, id)
	}
}
