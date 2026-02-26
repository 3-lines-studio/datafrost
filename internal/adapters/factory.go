package adapters

import (
	"datafrost/internal/models"
	"encoding/json"
	"fmt"
)

type Factory struct {
	adapters map[string]models.AdapterRegistration
}

func NewFactory() *Factory {
	factory := &Factory{
		adapters: make(map[string]models.AdapterRegistration),
	}

	factory.Register(NewTursoAdapterRegistration())
	factory.Register(NewPostgresAdapterRegistration())
	factory.Register(NewBigQueryAdapterRegistration())

	return factory
}

func (f *Factory) Register(reg models.AdapterRegistration) {
	f.adapters[reg.Info.Type] = reg
}

func (f *Factory) GetAdapter(adapterType string) (models.DatabaseAdapter, error) {
	reg, exists := f.adapters[adapterType]
	if !exists {
		return nil, fmt.Errorf("unknown adapter type: %s", adapterType)
	}
	return reg.Factory(), nil
}

func (f *Factory) GetAdapterInfo(adapterType string) (models.AdapterInfo, error) {
	reg, exists := f.adapters[adapterType]
	if !exists {
		return models.AdapterInfo{}, fmt.Errorf("unknown adapter type: %s", adapterType)
	}
	return reg.Info, nil
}

func (f *Factory) ListAdapters() []models.AdapterInfo {
	infos := make([]models.AdapterInfo, 0, len(f.adapters))
	for _, reg := range f.adapters {
		infos = append(infos, reg.Info)
	}
	return infos
}

func (f *Factory) TestConnection(adapterType string, credentials map[string]interface{}) error {
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

func SerializeCredentials(credentials map[string]interface{}) (string, error) {
	data, err := json.Marshal(credentials)
	if err != nil {
		return "", fmt.Errorf("failed to serialize credentials: %w", err)
	}
	return string(data), nil
}

func DeserializeCredentials(data string) (map[string]interface{}, error) {
	var credentials map[string]interface{}
	if err := json.Unmarshal([]byte(data), &credentials); err != nil {
		return nil, fmt.Errorf("failed to deserialize credentials: %w", err)
	}
	return credentials, nil
}
