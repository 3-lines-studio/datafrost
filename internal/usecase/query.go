package usecase

import (
	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase/port"
)

type QueryUsecase struct {
	connRepo port.ConnectionRepository
	cache    port.AdapterCache
}

func NewQueryUsecase(
	connRepo port.ConnectionRepository,
	cache port.AdapterCache,
) *QueryUsecase {
	return &QueryUsecase{
		connRepo: connRepo,
		cache:    cache,
	}
}

func (u *QueryUsecase) Execute(connectionID int64, query string) (*entity.QueryResult, error) {
	if query == "" {
		return nil, ErrQueryRequired
	}
	conn, err := u.connRepo.GetByID(connectionID)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, ErrConnectionNotFound
	}
	adapter, err := u.cache.Get(conn.ID, conn.Type, conn.Credentials)
	if err != nil {
		return nil, err
	}
	return adapter.ExecuteQuery(query)
}
