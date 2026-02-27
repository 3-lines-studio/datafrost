package usecase

import (
	"github.com/3-lines-studio/datafrost/internal/core/entity"
	"github.com/3-lines-studio/datafrost/internal/usecase/port"
)

type TableUsecase struct {
	connRepo port.ConnectionRepository
	cache    port.AdapterCache
}

func NewTableUsecase(
	connRepo port.ConnectionRepository,
	cache port.AdapterCache,
) *TableUsecase {
	return &TableUsecase{
		connRepo: connRepo,
		cache:    cache,
	}
}

func (u *TableUsecase) getAdapter(connectionID int64) (entity.DatabaseAdapter, *entity.Connection, error) {
	conn, err := u.connRepo.GetByID(connectionID)
	if err != nil {
		return nil, nil, err
	}
	if conn == nil {
		return nil, nil, ErrConnectionNotFound
	}
	adapter, err := u.cache.Get(conn.ID, conn.Type, conn.Credentials)
	if err != nil {
		return nil, nil, err
	}
	return adapter, conn, nil
}

func (u *TableUsecase) ListTables(connectionID int64) ([]entity.TableInfo, error) {
	adapter, _, err := u.getAdapter(connectionID)
	if err != nil {
		return nil, err
	}
	return adapter.ListTables()
}

func (u *TableUsecase) GetTableData(connectionID int64, tableName string, limit, offset int, filters []entity.Filter) (*entity.QueryResult, error) {
	adapter, _, err := u.getAdapter(connectionID)
	if err != nil {
		return nil, err
	}
	return adapter.GetTableData(tableName, limit, offset, filters)
}

func (u *TableUsecase) GetTableSchema(connectionID int64, tableName string) (*entity.TableSchema, error) {
	adapter, _, err := u.getAdapter(connectionID)
	if err != nil {
		return nil, err
	}
	return adapter.GetTableSchema(tableName)
}
