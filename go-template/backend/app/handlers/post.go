package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	userIDStr := ctx.FormValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	topicIdsStr := ctx.FormValue("topic_ids")
	var topicIds []int64
	if topicIdsStr != "" {
		for _, idStr := range strings.Split(topicIdsStr, ",") {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			topicIds = append(topicIds, id)
		}
	}

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

	// Topic ilişkilerini kaydet
	if len(topicIds) > 0 {
		if ps, ok := h.postService.(*service.PostService); ok {
			for _, topicId := range topicIds {
				postTopic := model.PostTopic{PostID: post.ID, TopicID: topicId}
				if _, err := ps.DB.NewInsert().Model(&postTopic).Exec(ctx.Context()); err != nil {
					return errorsx.InternalError(err, "Post-topic ilişkisi eklenemedi.")
				}
			}
		}
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

// Belirli bir post_id için ilişkili topic'leri döndüren endpoint
func (h PostHandler) GetTopicsByPostID(ctx *app.Ctx) error {
	postID := ctx.ParamsInt64("id")
	if postID == 0 {
		return errorsx.BadRequestError("Geçersiz post ID")
	}

	db := h.postService.(*service.PostService).DB
	type TopicResult struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	var topics []TopicResult
	err := db.NewSelect().
		Model((*model.PostTopic)(nil)).
		ColumnExpr("topic.id, topic.name").
		Join("JOIN topic ON topic.id = post_topic.topic_id").
		Where("post_topic.post_id = ?", postID).
		Scan(ctx.Context(), &topics)
	if err != nil {
		return errorsx.InternalError(err, "Topic'ler çekilemedi")
	}
	return ctx.JSON(topics)
}
