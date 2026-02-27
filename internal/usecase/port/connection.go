package port

import "github.com/3-lines-studio/datafrost/internal/core/entity"

type ConnectionRepository interface {
	Create(req entity.CreateConnectionRequest) (*entity.Connection, error)
	GetByID(id int64) (*entity.Connection, error)
	List() ([]entity.Connection, error)
	Delete(id int64) error
	Update(id int64, req entity.UpdateConnectionRequest) (*entity.Connection, error)
	SetLastConnected(id int64) error
	GetLastConnected() (int64, error)
}
