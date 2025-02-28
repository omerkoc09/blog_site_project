package handlers

import (
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/model"
	"github.com/hayrat/go-template2/backend/pkg/service"
)

type BaseHandler[
	T model.IDBModel,
	CreateVM model.IToDBModel[T],
	UpdateVM model.IToDBModel[T],
	ListVM model.IToViewModel[T, ListVM],
	DetailVM model.IToViewModel[T, DetailVM],
] struct {
	BaseQueryHandler[T, ListVM]
	BaseGetByIDHandler[T, DetailVM]
	BaseCreateHandler[T, CreateVM]
	BaseUpdateHandler[T, UpdateVM]
	BaseDeleteHandler[T]
}

func NewBaseHandler[
	T model.IDBModel,
	CreateVM model.IToDBModel[T],
	UpdateVM model.IToDBModel[T],
	ListVM model.IToViewModel[T, ListVM],
	DetailVM model.IToViewModel[T, DetailVM],
](s service.IBaseService[T], relations ...string) BaseHandler[T, CreateVM, UpdateVM, ListVM, DetailVM] {
	return BaseHandler[T, CreateVM, UpdateVM, ListVM, DetailVM]{
		BaseQueryHandler: BaseQueryHandler[T, ListVM]{
			BaseService: s,
			relations:   relations,
		},
		BaseGetByIDHandler: BaseGetByIDHandler[T, DetailVM]{
			BaseService: s,
			relations:   relations,
		},
		BaseCreateHandler: BaseCreateHandler[T, CreateVM]{
			BaseService: s,
		},
		BaseUpdateHandler: BaseUpdateHandler[T, UpdateVM]{
			BaseService: s,
		},
		BaseDeleteHandler: BaseDeleteHandler[T]{
			BaseService: s,
		},
	}
}

type BaseQueryHandler[T model.IDBModel, ListVM model.IToViewModel[T, ListVM]] struct {
	BaseService service.IBaseService[T]
	relations   []string
}

func NewBaseQueryHandler[T model.IDBModel, ListVM model.IToViewModel[T, ListVM]](s service.IBaseService[T], relations ...string) BaseQueryHandler[T, ListVM] {
	return BaseQueryHandler[T, ListVM]{BaseService: s, relations: relations}
}

type BaseGetByIDHandler[T model.IDBModel, DetailVM model.IToViewModel[T, DetailVM]] struct {
	BaseService service.IBaseService[T]
	relations   []string
}

func NewBaseGetByIDHandler[T model.IDBModel, DetailVM model.IToViewModel[T, DetailVM]](s service.IBaseService[T], relations ...string) BaseGetByIDHandler[T, DetailVM] {
	return BaseGetByIDHandler[T, DetailVM]{BaseService: s, relations: relations}
}

type BaseCreateHandler[T model.IDBModel, CreateVM model.IToDBModel[T]] struct {
	BaseService service.IBaseService[T]
}

func NewBaseCreateHandler[T model.IDBModel, CreateVM model.IToDBModel[T]](s service.IBaseService[T]) BaseCreateHandler[T, CreateVM] {
	return BaseCreateHandler[T, CreateVM]{BaseService: s}
}

type BaseUpdateHandler[T model.IDBModel, UpdateVM model.IToDBModel[T]] struct {
	BaseService service.IBaseService[T]
}

func NewBaseUpdateHandler[T model.IDBModel, UpdateVM model.IToDBModel[T]](s service.IBaseService[T]) BaseUpdateHandler[T, UpdateVM] {
	return BaseUpdateHandler[T, UpdateVM]{BaseService: s}
}

type BaseDeleteHandler[T model.IDBModel] struct {
	BaseService service.IBaseService[T]
}

func NewBaseDeleteHandler[T model.IDBModel](s service.IBaseService[T]) BaseDeleteHandler[T] {
	return BaseDeleteHandler[T]{BaseService: s}
}
func (h BaseQueryHandler[T, ListVM]) Query(ctx *app.Ctx) error {
	q, p, err := ctx.GetQueryPaginationModel()
	if err != nil {
		return errorsx.BadRequestError(err.Error())
	}
	data, dataCount, err := h.BaseService.GetQuery(ctx.Context(), q, p, h.relations...)
	if err != nil {
		return err
	}
	var result []any
	var lvm ListVM
	for _, d := range data {
		result = append(result, lvm.ToViewModel(d))
	}
	return ctx.SuccessResponse(result, dataCount)
}
func (h BaseQueryHandler[T, ListVM]) QueryWithDeleted(ctx *app.Ctx) error {
	q, p, err := ctx.GetQueryPaginationModel()
	if err != nil {
		return errorsx.BadRequestError(err.Error())
	}
	data, dataCount, err := h.BaseService.GetQueryWithDeleted(ctx.Context(), q, p, h.relations...)
	if err != nil {
		return err
	}
	var result []any
	var lvm ListVM
	for _, d := range data {
		result = append(result, lvm.ToViewModel(d))
	}
	return ctx.SuccessResponse(result, dataCount)
}

func (h BaseGetByIDHandler[T, DetailVM]) GetByID(ctx *app.Ctx) error {
	id := ctx.ParamsInt64("id")
	m, err := h.BaseService.GetByID(ctx.Context(), id, h.relations...)

	if err != nil {
		return err
	}
	var vm DetailVM
	result := vm.ToViewModel(m)
	return ctx.SuccessResponse(result)
}

func (h BaseCreateHandler[T, CreateVM]) Create(ctx *app.Ctx) error {
	var m T
	var vm CreateVM
	if errs := ctx.BodyParseValidate(&vm); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	m = vm.ToDBModel(m)
	err := h.BaseService.Create(ctx.Context(), &m)
	if err != nil {
		return err
	}

	return nil
}

func (h BaseUpdateHandler[T, UpdateVM]) Update(ctx *app.Ctx) error {
	id := ctx.ParamsInt64("id")
	m, err := h.BaseService.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	var vm UpdateVM
	if errs := ctx.BodyParseValidate(&vm); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}
	m = vm.ToDBModel(m)
	err = h.BaseService.Update(ctx.Context(), m)
	if err != nil {
		return err
	}

	return nil
}

func (h BaseDeleteHandler[T]) Delete(ctx *app.Ctx) error {
	id := ctx.ParamsInt64("id")
	err := h.BaseService.Delete(ctx.Context(), id)
	if err != nil {
		return err
	}

	return nil
}
