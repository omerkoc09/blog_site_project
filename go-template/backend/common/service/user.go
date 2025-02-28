package service

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/service"
)

type IUserService interface {
	service.IBaseService[model.User]
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type UserService struct {
	service.BaseService[model.User]
}

func NewUserService(db *bun.DB) IUserService {
	return &UserService{
		service.BaseService[model.User]{DB: db},
	}
}

func (s UserService) Create(ctx context.Context, m *model.User) error {
	exists, err := s.DB.NewSelect().Model(&model.User{}).Where("email = ?", m.Email).Exists(ctx)
	if err != nil {
		return errorsx.Database(err)
	}
	if exists {
		return errorsx.BadRequestError("Email başka bir kullanıcı tarafından kullanılmaktadır")
	}

	_, err = s.DB.NewInsert().Model(m).Exec(ctx)
	return errorsx.Database(err)
}

func (s UserService) Update(ctx context.Context, m model.User) error {
	exists, err := s.DB.NewSelect().Model(&model.User{}).Where("id != ? AND email = ?", m.ID, m.Email).Exists(ctx)
	if err != nil {
		return errorsx.Database(err)
	}
	if exists {
		return errorsx.BadRequestError("Email başka bir kullanıcı tarafından kullanılmaktadır")
	}

	_, err = s.DB.NewUpdate().
		Model(&m).
		WherePK().
		Exec(ctx)
	return errorsx.Database(err)
}

func (s UserService) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var m model.User
	err := s.DB.NewSelect().Model(&m).Where("email = ?", email).Scan(ctx)
	return m, errorsx.Database(err, "Hatalı Email veya Parola")
}
