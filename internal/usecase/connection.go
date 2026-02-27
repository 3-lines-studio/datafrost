package usecase

import (
	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase/port"
)

type ConnectionUsecase struct {
	repo    port.ConnectionRepository
	factory port.AdapterFactory
	cache   port.AdapterCache
}

func NewConnectionUsecase(
	repo port.ConnectionRepository,
	factory port.AdapterFactory,
	cache port.AdapterCache,
) *ConnectionUsecase {
	return &ConnectionUsecase{
		repo:    repo,
		factory: factory,
		cache:   cache,
	}
}

func (u *ConnectionUsecase) List() ([]entity.Connection, int64, error) {
	connections, err := u.repo.List()
	if err != nil {
		return nil, 0, err
	}
	lastID, _ := u.repo.GetLastConnected()
	return connections, lastID, nil
}

func (u *ConnectionUsecase) Create(req entity.CreateConnectionRequest) (*entity.Connection, error) {
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if req.Type == "" {
		return nil, ErrTypeRequired
	}
	return u.repo.Create(req)
}

func (u *ConnectionUsecase) Delete(id int64) error {
	u.cache.Invalidate(id)
	return u.repo.Delete(id)
}

func (u *ConnectionUsecase) Update(id int64, req entity.UpdateConnectionRequest) (*entity.Connection, error) {
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if req.Type == "" {
		return nil, ErrTypeRequired
	}
	u.cache.Invalidate(id)
	return u.repo.Update(id, req)
}

func (u *ConnectionUsecase) SetLastConnected(id int64) error {
	return u.repo.SetLastConnected(id)
}

func (u *ConnectionUsecase) Test(req entity.TestConnectionRequest) error {
	if req.Type == "" {
		return ErrTypeRequired
	}
	return u.factory.TestConnection(req.Type, req.Credentials)
}

func (u *ConnectionUsecase) TestExisting(id int64) error {
	conn, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	if conn == nil {
		return ErrConnectionNotFound
	}
	return u.factory.TestConnection(conn.Type, conn.Credentials)
}

func (u *ConnectionUsecase) GetConnection(id int64) (*entity.Connection, error) {
	conn, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, ErrConnectionNotFound
	}
	return conn, nil
}
