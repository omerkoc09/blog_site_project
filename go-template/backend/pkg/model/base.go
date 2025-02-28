package model

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type IDBModel interface {
	ModelName() string
}

type IToDBModel[T IDBModel] interface {
	ToDBModel(t T) T
}

type IToViewModel[T IDBModel, VM any] interface {
	ToViewModel(t T) VM
}

type BaseModel struct {
	ID        int64     `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	DeletedAt time.Time `json:"deleted_at" bun:",soft_delete,nullzero"`
}

func (u *BaseModel) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.UpdateQuery:
		u.UpdatedAt = time.Now()
	}
	return nil
}
