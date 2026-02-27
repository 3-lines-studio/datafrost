package port

import "github.com/3-lines-studio/datafrost/internal/core/entity"

type AdapterCache interface {
	Get(id int64, connType string, credentials map[string]any) (entity.DatabaseAdapter, error)
	Invalidate(id int64)
	Close()
}
