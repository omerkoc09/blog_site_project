package viewmodel

import (
	"fmt"

	"github.com/hayrat/go-template2/backend/common/model"
)

type PostCreateVM struct {
	ID          int64  `bun:",pk" json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MainContent string `json:"main_content"`
	UserId      int64  `json:"user_id"`
}

func (vm PostCreateVM) ToDBModel(m model.Post) model.Post {

	fmt.Printf("Received PostCreateVM: %+v\n", vm)

	m.ID = vm.ID
	m.Title = vm.Title
	m.Content = vm.Content
	m.MainContent = vm.MainContent
	m.UserId = vm.UserId

	return m
}

type PostUpdateVM struct {
	ID          int64  `bun:",pk" json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MainContent string `json:"main_content"`
	UserId      int64  `json:"user_id"`
}

func (vm PostUpdateVM) ToDBModel(m model.Post) model.Post {
	m.ID = vm.ID
	m.Title = vm.Title
	m.Content = vm.Content
	m.MainContent = vm.MainContent
	m.UserId = vm.UserId

	return m
}

type PostListVM struct {
	ID           int64         `bun:",pk" json:"id"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	MainContent  string        `json:"main_content"`
	Image        string        `json:"image"`
	UserId       int64         `json:"user_id"`
	LikeCount    int64         `json:"like_count"`
	CommentCount int64         `json:"comment_count"`
	Likes        []LikeMeVM    `bun:"rel:has-many,join:id=post_id,on_delete:cascade" json:"likes"`
	Comments     []CommentMeVM `bun:"rel:has-many,join:id=post_id,on_delete:cascade" json:"comments"`
}

func (vm PostListVM) ToViewModel(m model.Post) PostListVM {
	vm.ID = m.ID
	vm.Title = m.Title
	vm.Content = m.Content
	vm.MainContent = m.MainContent
	vm.Image = m.Image
	vm.UserId = m.UserId
	vm.LikeCount = m.LikeCount
	vm.CommentCount = m.CommentCount

	vm.Likes = make([]LikeMeVM, len(m.Likes))
	for i, like := range m.Likes {
		vm.Likes[i] = LikeMeVM{
			ID:     like.ID,
			UserId: like.UserId,
			PostId: like.PostId,
		}
	}

	vm.Comments = make([]CommentMeVM, len(m.Comments))
	for i, comment := range m.Comments {
		vm.Comments[i] = CommentMeVM{
			ID:      comment.ID,
			PostId:  comment.PostId,
			UserId:  comment.UserId,
			Content: comment.Content,
		}
	}

	return vm
}

type PostDetailVM struct {
	ID          int64  `bun:",pk" json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MainContent string `json:"main_content"`
	Image       string `json:"image"`
	UserId      int64  `json:"user_id"`
}

func (vm PostDetailVM) ToViewModel(m model.Post) PostDetailVM {
	vm.ID = m.ID
	vm.Title = m.Title
	vm.Content = m.Content
	vm.MainContent = m.MainContent
	vm.UserId = m.UserId
	vm.Image = m.Image

	return vm
}

type PostMeVM struct {
	ID           int64         `bun:",pk" json:"id"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	MainContent  string        `json:"main_content"`
	Image        string        `json:"image"`
	UserId       int64         `json:"user_id"`
	LikeCount    int64         `json:"like_count"`
	CommentCount int64         `json:"comment_count"`
	Likes        []LikeMeVM    `bun:"rel:has-many,join:id=post_id,on_delete:cascade" json:"likes"`
	Comments     []CommentMeVM `bun:"rel:has-many,join:id=post_id,on_delete:cascade" json:"comments"`
}

func (vm PostMeVM) ToViewModel(m model.Post) PostMeVM {
	vm.ID = m.ID
	vm.Title = m.Title
	vm.Content = m.Content
	vm.MainContent = m.MainContent
	vm.Image = m.Image
	vm.UserId = m.UserId
	vm.LikeCount = m.LikeCount
	vm.CommentCount = m.CommentCount

	vm.Likes = make([]LikeMeVM, len(m.Likes))
	for i, like := range m.Likes {
		vm.Likes[i] = LikeMeVM{
			ID:     like.ID,
			UserId: like.UserId,
			PostId: like.PostId,
		}
	}

	vm.Comments = make([]CommentMeVM, len(m.Comments))
	for i, comment := range m.Comments {
		vm.Comments[i] = CommentMeVM{
			ID:      comment.ID,
			PostId:  comment.PostId,
			UserId:  comment.UserId,
			Content: comment.Content,
		}
	}

	return vm
}

type PostMeUpdateVM struct {
	ID          int64  `bun:",pk" json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MainContent string `json:"main_content"`
	Image       string `json:"image"`
	UserId      int64  `json:"user_id"`
}

func (vm PostMeUpdateVM) ToDBModel(m model.Post) model.Post {
	m.ID = vm.ID
	m.Title = vm.Title
	m.Content = vm.Content
	m.MainContent = vm.MainContent
	m.Image = vm.Image
	m.UserId = vm.UserId

	return m
}
