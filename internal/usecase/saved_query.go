package usecase

import (
	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase/port"
)

type SavedQueryUsecase struct {
	repo port.SavedQueryRepository
}

func NewSavedQueryUsecase(repo port.SavedQueryRepository) *SavedQueryUsecase {
	return &SavedQueryUsecase{repo: repo}
}

func (u *SavedQueryUsecase) ListByConnection(connectionID int64) ([]entity.SavedQuery, error) {
	return u.repo.ListByConnection(connectionID)
}

func (u *SavedQueryUsecase) Create(connectionID int64, name, query string) (*entity.SavedQuery, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if query == "" {
		return nil, ErrQueryRequired
	}
	return u.repo.Create(connectionID, name, query)
}

func (u *SavedQueryUsecase) Update(id int64, name, query string) (*entity.SavedQuery, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if query == "" {
		return nil, ErrQueryRequired
	}
	q, err := u.repo.Update(id, name, query)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, ErrQueryNotFound
	}
	return q, nil
}

func (u *SavedQueryUsecase) Delete(id int64) error {
	return u.repo.Delete(id)
}
