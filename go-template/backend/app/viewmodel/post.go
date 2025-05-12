package viewmodel

import (
	"github.com/hayrat/go-template2/backend/common/model"
	pkgModel "github.com/hayrat/go-template2/backend/pkg/model"
)

type PostCreateVM struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	MainContent string  `json:"main_content"`
	Image       string  `json:"image"`
	UserId      int64   `json:"user_id"`
	TopicIds    []int64 `json:"topic_ids"`
}

func (vm PostCreateVM) ToDBModel(m model.Post) model.Post {
	m.ID = vm.ID
	m.Title = vm.Title
	m.Content = vm.Content
	m.MainContent = vm.MainContent
	m.Image = vm.Image
	m.UserId = vm.UserId

	if len(vm.TopicIds) > 0 {
		m.Topics = make([]model.Topic, 0, len(vm.TopicIds))
		for _, topicId := range vm.TopicIds {
			m.Topics = append(m.Topics, model.Topic{BaseModel: pkgModel.BaseModel{ID: topicId}})
		}
	}

	return m
}

type PostListVM struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	MainContent string        `json:"main_content"`
	Image       string        `json:"image"`
	UserId      int64         `json:"user_id"`
	Likes       []LikeMeVM    `json:"likes"`
	Comments    []CommentMeVM `json:"comments"`
}

func (vm PostListVM) ToViewModel(m model.Post) PostListVM {
	vm.ID = m.ID
	vm.Title = m.Title
	vm.Content = m.Content
	vm.MainContent = m.MainContent
	vm.Image = m.Image
	vm.UserId = m.UserId

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
	ID          int64  `json:"id"`
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
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	MainContent string        `json:"main_content"`
	Image       string        `json:"image"`
	UserId      int64         `json:"user_id"`
	Likes       []LikeMeVM    `json:"likes"`
	Comments    []CommentMeVM `json:"comments"`
}

func (vm PostMeVM) ToViewModel(m model.Post) PostMeVM {
	vm.ID = m.ID
	vm.Title = m.Title
	vm.Content = m.Content
	vm.MainContent = m.MainContent
	vm.Image = m.Image
	vm.UserId = m.UserId

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
	ID          int64  `json:"id"`
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

type DummyPostVM struct{}

func (d DummyPostVM) ToDBModel(m model.Post) model.Post {
	return m // Return an empty Contact or the appropriate type
}

type DummyPostVM2 struct{}

func (d DummyPostVM2) ToViewModel(m model.Post) DummyPostVM2 {
	return DummyPostVM2{} // Return an empty Contact or the appropriate type
}
