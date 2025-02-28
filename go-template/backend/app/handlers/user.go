package handlers

import (
	"github.com/hayrat/go-template2/backend/app/viewmodel"
	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/common/service"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	base "github.com/hayrat/go-template2/backend/pkg/handlers"
)

type UserHandler struct {
	base.BaseHandler[model.User, viewmodel.UserCreateVM, viewmodel.DummyVM, viewmodel.UserListVM, viewmodel.UserDetailVM]
	userService service.IUserService
}

func NewUserHandler(s service.IUserService) UserHandler {
	h := UserHandler{
		BaseHandler: base.NewBaseHandler[model.User, viewmodel.UserCreateVM, viewmodel.DummyVM, viewmodel.UserListVM, viewmodel.UserDetailVM](s),
		userService: s,
	}

	return h
}

func (h UserHandler) Me(ctx *app.Ctx) error {
	id := ctx.GetUserID()
	user, err := h.userService.GetByID(ctx.Context(), id, "Comments", "Likes", "Posts")
	if err != nil {
		return err
	}

	result := viewmodel.UserMeVM{}.ToViewModel(user)

	return ctx.SuccessResponse(result)
}

func (h UserHandler) MeUpdate(ctx *app.Ctx) error {
	id := ctx.GetUserID()
	m, err := h.userService.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	var vm viewmodel.UserMeUpdateVM
	if errs := ctx.BodyParseValidate(&vm); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	m = vm.ToDBModel(m)
	err = h.userService.Update(ctx.Context(), m)
	if err != nil {
		return err
	}

	return nil
}

func (h UserHandler) GetByID(ctx *app.Ctx) error {

	var user model.User
	// Get postID from the request (URL parameter in this case)
	userID := ctx.ParamsInt64("id") // Assuming the URL is `/posts/:id`
	if userID == 0 {
		return errorsx.BadRequestError("Invalid post ID")
	}
	user, err := h.userService.GetByID(ctx.Context(), userID, "Comments", "Likes", "Posts")

	if err != nil {
		return err
	}

	result := viewmodel.UserMeVM{}.ToViewModel(user)

	return ctx.SuccessResponse(result)
}
