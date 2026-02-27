package port

import "github.com/3-lines-studio/datafrost/internal/core/entity"

type AdapterFactory interface {
	GetAdapter(adapterType string) (entity.DatabaseAdapter, error)
	GetAdapterInfo(adapterType string) (entity.AdapterInfo, error)
	ListAdapters() []entity.AdapterInfo
	TestConnection(adapterType string, credentials map[string]any) error
}
