package port

import "github.com/3-lines-studio/datafrost/internal/core/entity"

type SavedQueryRepository interface {
	ListByConnection(connectionID int64) ([]entity.SavedQuery, error)
	GetByID(id int64) (*entity.SavedQuery, error)
	Create(connectionID int64, name, query string) (*entity.SavedQuery, error)
	Update(id int64, name, query string) (*entity.SavedQuery, error)
	Delete(id int64) error
}
