package viewmodel

import "github.com/hayrat/go-template2/backend/common/model"

type UserVM struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type CommentCreateVM struct {
	ID      int64  `bun:",pk,autoincrement" json:"id"`
	PostId  int64  `bun:",notnull" json:"post_id"`
	UserId  int64  `bun:",notnull" json:"user_id"`
	Content string `bun:",notnull" json:"content"`
}

func (vm CommentCreateVM) ToDBModel(m model.Comment) model.Comment {
	m.ID = vm.ID
	m.PostId = vm.PostId
	m.UserId = vm.UserId
	m.Content = vm.Content

	return m
}

type CommentUpdateVM struct {
	ID      int64  `bun:",pk,autoincrement" json:"id"`
	PostId  int64  `bun:",notnull" json:"post_id"`
	UserId  int64  `bun:",notnull" json:"user_id"`
	Content string `bun:",notnull" json:"content"`
}

func (vm CommentUpdateVM) ToDBModel(m model.Comment) model.Comment {
	m.ID = vm.ID
	m.PostId = vm.PostId
	m.UserId = vm.UserId
	m.Content = vm.Content

	return m

}

type CommentListVM struct {
	ID      int64  `bun:",pk,autoincrement" json:"id"`
	PostId  int64  `bun:",notnull" json:"post_id"`
	UserId  int64  `bun:",notnull" json:"user_id"`
	Content string `bun:",notnull" json:"content"`
}

func (vm CommentListVM) ToViewModel(m model.Comment) CommentListVM {
	vm.ID = m.ID
	vm.PostId = m.PostId
	vm.UserId = m.UserId
	vm.Content = m.Content

	return vm
}

type CommentDetailVM struct {
	ID      int64  `bun:",pk,autoincrement" json:"id"`
	PostId  int64  `bun:",notnull" json:"post_id"`
	UserId  int64  `bun:",notnull" json:"user_id"`
	Content string `bun:",notnull" json:"content"`
}

func (vm CommentDetailVM) ToViewModel(m model.Comment) CommentDetailVM {
	vm.ID = m.ID
	vm.PostId = m.PostId
	vm.UserId = m.UserId
	vm.Content = m.Content

	return vm
}

type CommentMeVM struct {
	ID      int64  `bun:",pk,autoincrement" json:"id"`
	PostId  int64  `bun:",notnull" json:"post_id"`
	UserId  int64  `bun:",notnull" json:"user_id"`
	Content string `bun:",notnull" json:"content"`
}

func (vm CommentMeVM) ToViewModel(m model.Comment) CommentMeVM {
	vm.ID = m.ID
	vm.PostId = m.PostId
	vm.UserId = m.UserId
	vm.Content = m.Content

	return vm
}

type CommentMeUpdateVM struct {
	ID      int64  `bun:",pk,autoincrement" json:"id"`
	PostId  int64  `bun:",notnull" json:"post_id"`
	UserId  int64  `bun:",notnull" json:"user_id"`
	Content string `bun:",notnull" json:"content"`
}

func (vm CommentMeUpdateVM) ToDBModel(m model.Comment) model.Comment {
	m.ID = vm.ID
	m.PostId = vm.PostId
	m.UserId = vm.UserId
	m.Content = vm.Content

	return m
}
