package database

import (
	"fmt"

	"github.com/3-lines-studio/datafrost/internal/core/entity"
)

type Factory struct {
	adapters map[string]entity.AdapterRegistration
}

func NewFactory() *Factory {
	factory := &Factory{
		adapters: make(map[string]entity.AdapterRegistration),
	}

	factory.Register(newTursoAdapterRegistration())
	factory.Register(newPostgresAdapterRegistration())
	factory.Register(newBigQueryAdapterRegistration())

	return factory
}

func (f *Factory) Register(reg entity.AdapterRegistration) {
	f.adapters[reg.Info.Type] = reg
}

func (f *Factory) GetAdapter(adapterType string) (entity.DatabaseAdapter, error) {
	reg, exists := f.adapters[adapterType]
	if !exists {
		return nil, fmt.Errorf("unknown adapter type: %s", adapterType)
	}
	return reg.Factory(), nil
}

func (f *Factory) GetAdapterInfo(adapterType string) (entity.AdapterInfo, error) {
	reg, exists := f.adapters[adapterType]
	if !exists {
		return entity.AdapterInfo{}, fmt.Errorf("unknown adapter type: %s", adapterType)
	}
	return reg.Info, nil
}

func (f *Factory) ListAdapters() []entity.AdapterInfo {
	infos := make([]entity.AdapterInfo, 0, len(f.adapters))
	for _, reg := range f.adapters {
		infos = append(infos, reg.Info)
	}
	return infos
}

func (f *Factory) TestConnection(adapterType string, credentials map[string]any) error {
	adapter, err := f.GetAdapter(adapterType)
	if err != nil {
		return err
	}
	defer adapter.Close()

	if err := adapter.Connect(credentials); err != nil {
		return err
	}

	return adapter.Ping()
}
