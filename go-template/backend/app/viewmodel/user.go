package viewmodel

import (
	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/pkg/utils"
)

type UserLoginVM struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreateVM struct {
	Email    string         `json:"email"`
	Phone    string         `json:"phone"`
	Name     string         `json:"name"`
	Surname  string         `json:"surname"`
	Role     model.UserRole `json:"role"`
	Password string         `json:"password"`
	About    string         `json:"about"`
}

func (vm UserCreateVM) ToDBModel(m model.User) model.User {
	m.Email = utils.EmailTemizle(vm.Email)
	m.Phone = utils.TelefonTemizle(vm.Phone)
	m.Name = utils.ToTitle(vm.Name)
	m.Surname = utils.ToTitle(vm.Surname)
	m.Role = vm.Role
	m.Password, _ = utils.HashPassword(vm.Password)
	m.About = utils.ToTitle(vm.About)

	return m
}

type UserListVM struct {
	ID      int64          `json:"id"`
	Email   string         `json:"email"`
	Phone   string         `json:"phone"`
	Name    string         `json:"name"`
	Surname string         `json:"surname"`
	About   string         `json:"about"`
	Role    model.UserRole `json:"role"`
}

func (vm UserListVM) ToViewModel(m model.User) UserListVM {
	vm.ID = m.ID
	vm.Email = m.Email
	vm.Phone = m.Phone
	vm.Name = m.Name
	vm.Surname = m.Surname
	vm.Role = m.Role
	vm.About = m.About

	return vm
}

type UserDetailVM struct {
	ID      int64          `json:"id"`
	Email   string         `json:"email"`
	Phone   string         `json:"phone"`
	Name    string         `json:"name"`
	Surname string         `json:"surname"`
	About   string         `json:"about"`
	Role    model.UserRole `json:"role"`
}

func (vm UserDetailVM) ToViewModel(m model.User) UserDetailVM {
	vm.ID = m.ID
	vm.Email = m.Email
	vm.Phone = m.Phone
	vm.Name = m.Name
	vm.Surname = m.Surname
	vm.Role = m.Role
	vm.About = m.About

	return vm
}

type UserMeVM struct {
	ID       int64          `json:"id"`
	Email    string         `json:"email"`
	Phone    string         `json:"phone"`
	Name     string         `json:"name"`
	Surname  string         `json:"surname"`
	About    string         `json:"about"`
	Role     model.UserRole `json:"role"`
	Comments []CommentMeVM
	Likes    []LikeMeVM
	Posts    []PostMeVM
}

func (vm UserMeVM) ToViewModel(m model.User) UserMeVM {
	vm.ID = m.ID
	vm.Email = m.Email
	vm.Phone = m.Phone
	vm.Name = m.Name
	vm.Surname = m.Surname
	vm.Role = m.Role
	vm.About = m.About

	vm.Comments = make([]CommentMeVM, len(m.Comments))
	for i, comment := range m.Comments {
		vm.Comments[i] = CommentMeVM{
			ID:      comment.ID,
			PostId:  comment.PostId,
			UserId:  comment.UserId,
			Content: comment.Content,
		}
	}

	vm.Likes = make([]LikeMeVM, len(m.Likes))
	for i, like := range m.Likes {
		vm.Likes[i] = LikeMeVM{
			ID:     like.ID,
			UserId: like.UserId,
			PostId: like.PostId,
		}
	}

	vm.Posts = make([]PostMeVM, len(m.Posts))
	for i, post := range m.Posts {
		vm.Posts[i] = PostMeVM{
			ID:      post.ID,
			UserId:  post.UserId,
			Title:   post.Title,
			Content: post.Content,
		}
	}

	return vm
}

type UserMeUpdateVM struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
	About    string `json:"about"`
}

func (vm UserMeUpdateVM) ToDBModel(m model.User) model.User {
	m.Email = utils.EmailTemizle(vm.Email)
	m.Phone = utils.TelefonTemizle(vm.Phone)
	m.Name = utils.ToTitle(vm.Name)
	m.Surname = utils.ToTitle(vm.Surname)
	m.About = utils.ToTitle(vm.About)
	if vm.Password != "" {
		m.Password, _ = utils.HashPassword(vm.Password)
	}

	return m
}

type DummyVM struct{}

func (d DummyVM) ToDBModel(m model.User) model.User {
	return m // Return an empty Contact or the appropriate type
}

type DummyVM2 struct{}

func (d DummyVM2) ToViewModel(m model.User) DummyVM2 {
	return DummyVM2{} // Return an empty Contact or the appropriate type
}
