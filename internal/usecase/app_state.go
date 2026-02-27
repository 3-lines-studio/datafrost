package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase/port"
)

type AppStateUsecase struct {
	repo port.AppStateRepository
}

func NewAppStateUsecase(repo port.AppStateRepository) *AppStateUsecase {
	return &AppStateUsecase{repo: repo}
}

func (u *AppStateUsecase) GetTheme() (string, error) {
	value, err := u.repo.Get("theme")
	if err != nil {
		return "light", nil
	}
	if value == "" {
		return "light", nil
	}
	return value, nil
}

func (u *AppStateUsecase) SetTheme(theme string) error {
	return u.repo.Set("theme", theme)
}

func (u *AppStateUsecase) GetLayout(key string) (string, error) {
	value, err := u.repo.Get("layout_" + key)
	if err != nil {
		return "", nil
	}
	return value, nil
}

func (u *AppStateUsecase) SaveLayout(key, layout string) error {
	return u.repo.Set("layout_"+key, layout)
}

func (u *AppStateUsecase) GetTabs(connectionID int) ([]entity.Tab, error) {
	key := fmt.Sprintf("tabs_%d", connectionID)
	value, err := u.repo.Get(key)
	if err != nil {
		return []entity.Tab{}, nil
	}
	if value == "" {
		return []entity.Tab{}, nil
	}
	var tabs []entity.Tab
	if err := json.Unmarshal([]byte(value), &tabs); err != nil {
		return []entity.Tab{}, nil
	}
	return tabs, nil
}

func (u *AppStateUsecase) SaveTabs(connectionID int, tabs []entity.Tab) error {
	key := fmt.Sprintf("tabs_%d", connectionID)
	data, err := json.Marshal(tabs)
	if err != nil {
		return err
	}
	return u.repo.Set(key, string(data))
}
