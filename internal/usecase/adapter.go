package usecase

import (
	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase/port"
)

type AdapterUsecase struct {
	factory port.AdapterFactory
}

func NewAdapterUsecase(factory port.AdapterFactory) *AdapterUsecase {
	return &AdapterUsecase{factory: factory}
}

func (u *AdapterUsecase) List() []entity.AdapterInfo {
	return u.factory.ListAdapters()
}
