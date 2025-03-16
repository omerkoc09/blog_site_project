package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hayrat/go-template2/backend/app/viewmodel"
	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/common/service"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	base "github.com/hayrat/go-template2/backend/pkg/handlers"
)

type PostHandler struct {
	base.BaseHandler[model.Post, viewmodel.PostCreateVM, viewmodel.DummyPostVM, viewmodel.PostListVM, viewmodel.PostDetailVM]
	postService service.IPostService
}

func NewPostHandler(s service.IPostService) PostHandler {
	h := PostHandler{
		BaseHandler: base.NewBaseHandler[model.Post, viewmodel.PostCreateVM, viewmodel.DummyPostVM, viewmodel.PostListVM, viewmodel.PostDetailVM](s),
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

func (h PostHandler) MeUpdateWithImage(ctx *app.Ctx) error {
	postID := ctx.ParamsInt64("id")
	if postID == 0 {
		return errorsx.BadRequestError("Invalid post ID")
	}

	// Get existing post
	post, err := h.postService.GetByID(ctx.Context(), postID)
	if err != nil {
		return err
	}

	// Update text fields from form data
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	mainContent := ctx.FormValue("main_content")

	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}
	if mainContent != "" {
		post.MainContent = mainContent
	}

	// Handle image upload if present
	file, err := ctx.FormFile("image")
	if err == nil {
		// Create timestamp-based filename
		timestamp := time.Now().UnixNano()
		filename := fmt.Sprintf("post_%d_%d%s", postID, timestamp, filepath.Ext(file.Filename))

		// Ensure upload directory exists
		uploadDir := "uploads/images"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return errorsx.InternalError(err, "Failed to create upload directory")
		}

		// Set up file path
		uploadPath := filepath.Join(uploadDir, filename)

		// Open source file
		src, err := file.Open()
		if err != nil {
			return errorsx.InternalError(err, "Failed to open uploaded file")
		}
		defer src.Close()

		// Create destination file
		dst, err := os.Create(uploadPath)
		if err != nil {
			return errorsx.InternalError(err, "Failed to create destination file")
		}
		defer dst.Close()

		// Copy file contents
		if _, err = io.Copy(dst, src); err != nil {
			return errorsx.InternalError(err, "Failed to save file")
		}

		// Update post image path
		post.Image = uploadPath
	}

	// Save all changes
	err = h.postService.Update(ctx.Context(), post)
	if err != nil {
		return errorsx.InternalError(err, "Failed to update post")
	}

	return ctx.SuccessResponse(post)
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

func (h PostHandler) CreatePostWithImage(ctx *app.Ctx) error {
	// Form verilerini al
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	mainContent := ctx.FormValue("main_content")
	userIDStr := "1"
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return errorsx.BadRequestError("Geçerli bir kullanıcı ID'si giriniz: " + err.Error())
	}

	// Eksik alan kontrolü
	if title == "" || content == "" || mainContent == "" {
		return errorsx.BadRequestError("Tüm alanlar doldurulmalıdır.")
	}

	// Yeni post nesnesi oluştur
	post := model.Post{
		Title:       title,
		Content:     content,
		MainContent: mainContent,
		UserId:      int64(userID),
	}

	// Post'u kaydet
	err = h.postService.Create(ctx.Context(), &post)
	if err != nil {
		return errorsx.InternalError(err, "Post oluşturulamadı.")
	}

	// **Resim ekleme işlemi**
	file, err := ctx.FormFile("image")
	if err == nil {
		// Dosya adını oluştur
		timestamp := time.Now().UnixNano()
		filename := fmt.Sprintf("post_%d_%d%s", post.ID, timestamp, filepath.Ext(file.Filename))

		// uploads/images klasörünün varlığını kontrol et
		uploadDir := "uploads/images"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return errorsx.InternalError(err, "Klasör oluşturulamadı.")
		}

		// Dosya kaydetme işlemi
		uploadPath := filepath.Join(uploadDir, filename)
		src, err := file.Open()
		if err != nil {
			return errorsx.InternalError(err, "Dosya açılamadı.")
		}
		defer src.Close()

		dst, err := os.Create(uploadPath)
		if err != nil {
			return errorsx.InternalError(err, "Dosya oluşturulamadı.")
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return errorsx.InternalError(err, "Dosya kaydedilemedi.")
		}

		// Post'a resim ekleyip güncelle
		post.Image = uploadPath
		err = h.postService.Update(ctx.Context(), post)
		if err != nil {
			return errorsx.InternalError(err, "Post güncellenemedi.")
		}
	}

	// Başarı cevabı
	return ctx.JSON(post)
}
