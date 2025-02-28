package service

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/model"
	"github.com/hayrat/go-template2/backend/pkg/viewmodel"
)

type IBaseService[T model.IDBModel] interface {
	Create(ctx context.Context, t *T) error
	CreateArr(ctx context.Context, t []T) error
	GetByID(ctx context.Context, id int64, relations ...string) (T, error)
	GetByColumn(ctx context.Context, column string, value any, relations ...string) (T, error)
	GetQuery(ctx context.Context, q *viewmodel.QueryModel, p *viewmodel.PaginationModel, relations ...string) ([]T, int, error)
	GetQueryWithDeleted(ctx context.Context, q *viewmodel.QueryModel, p *viewmodel.PaginationModel, relations ...string) ([]T, int, error)
	Update(ctx context.Context, t T) error
	UpdateColumn(ctx context.Context, t T, column string, value any) error
	Delete(ctx context.Context, id int64) error
}

// BaseService methodları override edilemez. Go generics yüzünden override ettirmiyor.
type BaseService[T model.IDBModel] struct {
	DB bun.IDB
}

func (s BaseService[T]) Create(ctx context.Context, t *T) error {
	_, err := s.DB.NewInsert().Model(t).Exec(ctx)
	return errorsx.Database(err)
}

func (s BaseService[T]) CreateArr(ctx context.Context, t []T) error {
	_, err := s.DB.NewInsert().Model(&t).Exec(ctx)
	return errorsx.Database(err)
}

func (s BaseService[T]) GetByID(ctx context.Context, id int64, relations ...string) (T, error) {
	return s.GetByColumn(ctx, "id", id, relations...)
}

func (s BaseService[T]) GetByColumn(ctx context.Context, column string, value any, relations ...string) (T, error) {
	var t T
	sq := s.DB.NewSelect().Model(&t).Where("? = ?", bun.Ident(t.ModelName()+"."+column), value)

	for _, relation := range relations {
		sq = sq.Relation(relation)
	}

	err := sq.Scan(ctx)
	return t, errorsx.Database(err)
}

func (s BaseService[T]) getQuery(ctx context.Context, widthDeleted bool, q *viewmodel.QueryModel, p *viewmodel.PaginationModel, relations ...string) ([]T, int, error) {
	var t []T
	sq := s.DB.NewSelect().Model(&t)
	if widthDeleted {
		sq = sq.WhereAllWithDeleted()
	}
	for _, relation := range relations {
		sq = sq.Relation(relation)
	}

	if q != nil {
		sq = Filter(sq, q)
	}
	if p != nil {
		sq = Paginate(sq, p)
		dataCount, err := sq.ScanAndCount(ctx)
		return t, dataCount, errorsx.Database(err)
	}
	err := sq.Scan(ctx)
	return t, 0, errorsx.Database(err)
}

func (s BaseService[T]) GetQuery(ctx context.Context, q *viewmodel.QueryModel, p *viewmodel.PaginationModel, relations ...string) ([]T, int, error) {
	return s.getQuery(ctx, false, q, p, relations...)
}

func (s BaseService[T]) GetQueryWithDeleted(ctx context.Context, q *viewmodel.QueryModel, p *viewmodel.PaginationModel, relations ...string) ([]T, int, error) {
	return s.getQuery(ctx, true, q, p, relations...)
}

func (s BaseService[T]) Update(ctx context.Context, t T) error {
	_, err := s.DB.NewUpdate().
		Model(&t).
		//OmitZero().
		WherePK().
		Exec(ctx)
	return errorsx.Database(err)
}

func (s BaseService[T]) UpdateColumn(ctx context.Context, t T, column string, value any) error {
	_, err := s.DB.NewUpdate().
		Model(&t).
		WherePK().
		Set("? = ?", bun.Ident(t.ModelName()+"."+column), value).
		Exec(ctx)
	return errorsx.Database(err)
}

func (s BaseService[T]) Delete(ctx context.Context, id int64) error {
	var t T
	_, err := s.DB.NewDelete().Model(&t).Where("id = ?", id).Exec(ctx)
	return errorsx.Database(err)
}
