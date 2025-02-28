package viewmodel

import "github.com/hayrat/go-template2/backend/common/model"

type LikeCreateVM struct {
	ID     int64 `bun:",pk" json:"id"`
	UserId int64 `json:"user_id"`
	PostId int64 `json:"post_id"`
}

func (vm LikeCreateVM) ToDBModel(m model.Like) model.Like {
	m.ID = vm.ID
	m.UserId = vm.UserId
	m.PostId = vm.PostId

	return m
}

type LikeUpdateVM struct {
	ID     int64 `bun:",pk" json:"id"`
	UserId int64 `json:"user_id"`
	PostId int64 `json:"post_id"`
}

func (vm LikeUpdateVM) ToDBModel(m model.Like) model.Like {
	m.ID = vm.ID
	m.UserId = vm.UserId
	m.PostId = vm.PostId

	return m

}

type LikeListVM struct {
	ID     int64 `bun:",pk" json:"id"`
	UserId int64 `json:"user_id"`
	PostId int64 `json:"post_id"`
}

func (vm LikeListVM) ToViewModel(m model.Like) LikeListVM {
	vm.ID = m.ID
	vm.UserId = m.UserId
	vm.PostId = m.PostId

	return vm
}

type LikeDetailVM struct {
	ID     int64 `bun:",pk" json:"id"`
	UserId int64 `json:"user_id"`
	PostId int64 `json:"post_id"`
}

func (vm LikeDetailVM) ToViewModel(m model.Like) LikeDetailVM {
	vm.ID = m.ID
	vm.UserId = m.UserId
	vm.PostId = m.PostId

	return vm
}

type LikeMeVM struct {
	ID     int64 `bun:",pk" json:"id"`
	UserId int64 `json:"user_id"`
	PostId int64 `json:"post_id"`
}

func (vm LikeMeVM) ToViewModel(m model.Like) LikeMeVM {
	vm.ID = m.ID
	vm.UserId = m.UserId
	vm.PostId = m.PostId

	return vm
}

type LikeMeUpdateVM struct {
	ID     int64 `bun:",pk" json:"id"`
	UserId int64 `json:"user_id"`
	PostId int64 `json:"post_id"`
}

func (vm LikeMeUpdateVM) ToDBModel(m model.Like) model.Like {
	m.ID = vm.ID
	m.UserId = vm.UserId
	m.PostId = vm.PostId

	return m
}
