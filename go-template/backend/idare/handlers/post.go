package handlers

import (
	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/common/service"
	"github.com/hayrat/go-template2/backend/idare/viewmodel"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	base "github.com/hayrat/go-template2/backend/pkg/handlers"
)

type PostHandler struct {
	base.BaseHandler[model.Post, viewmodel.PostCreateVM, viewmodel.PostUpdateVM, viewmodel.PostListVM, viewmodel.PostDetailVM]
	postService service.IPostService
}

func NewPostHandler(s service.IPostService) PostHandler {
	h := PostHandler{
		BaseHandler: base.NewBaseHandler[model.Post, viewmodel.PostCreateVM, viewmodel.PostUpdateVM, viewmodel.PostListVM, viewmodel.PostDetailVM](s),
		postService: s,
	}

	return h
}

func (h PostHandler) GetByID(ctx *app.Ctx) error {

	var post model.Post
	// Get postID from the request (URL parameter in this case)
	postID := ctx.ParamsInt64("id") // Assuming the URL is `/posts/:id`
	if postID == 0 {
		return errorsx.BadRequestError("Invalid post ID")
	}
	post, err := h.postService.GetByID(ctx.Context(), postID, "Comments", "Likes")

	if err != nil {
		return err
	}

	likeCount := len(post.Likes)

	commentCount := len(post.Comments)

	post.LikeCount = int64(likeCount)

	post.CommentCount = int64(commentCount)

	result := viewmodel.PostMeVM{}.ToViewModel(post)

	return ctx.SuccessResponse(result)
}

func (h PostHandler) MeUpdate(ctx *app.Ctx) error {
	postID := ctx.ParamsInt64("id") // Assuming the URL is `/posts/:id`
	if postID == 0 {
		return errorsx.BadRequestError("Invalid post ID")
	}
	m, err := h.postService.GetByID(ctx.Context(), postID)
	if err != nil {
		return err
	}

	var vm viewmodel.PostMeUpdateVM
	if errs := ctx.BodyParseValidate(&vm); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	m = vm.ToDBModel(m)
	err = h.postService.Update(ctx.Context(), m)
	if err != nil {
		return err
	}

	return nil
}

func (h PostHandler) Query(ctx *app.Ctx) error {
	q, p, err := ctx.GetQueryPaginationModel()
	if err != nil {
		return errorsx.BadRequestError(err.Error())
	}

	relations := []string{"Likes", "Comments"}

	data, dataCount, err := h.postService.GetQuery(ctx.Context(), q, p, relations...)
	if err != nil {
		return err
	}

	// Update counts before converting to view models
	for i := range data {
		// Update like_count based on Likes slice length
		data[i].LikeCount = int64(len(data[i].Likes))

		// Update comment_count based on Comments slice length
		data[i].CommentCount = int64(len(data[i].Comments))
	}

	var result []viewmodel.PostListVM
	for _, d := range data {
		var vm viewmodel.PostListVM
		result = append(result, vm.ToViewModel(d))
	}

	return ctx.SuccessResponse(result, dataCount)
}

func (h PostHandler) DeleteImageHandler(ctx *app.Ctx) error {

	return h.postService.DeleteImage(ctx)
}
