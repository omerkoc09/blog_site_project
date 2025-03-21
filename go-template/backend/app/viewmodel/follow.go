package viewmodel

import "github.com/hayrat/go-template2/backend/common/model"

type FollowCreateVM struct {
	FollowingID int64 `json:"following_id"`
}

type FollowUpdateVM struct {
	FollowingID int64 `json:"following_id"`
}

type FollowListVM struct {
	ID          int64    `json:"id"`
	FollowerID  int64    `json:"follower_id"`
	FollowingID int64    `json:"following_id"`
	Follower    UserInfo `json:"follower,omitempty"`
	Following   UserInfo `json:"following,omitempty"`
}

type FollowDetailVM struct {
	ID          int64    `json:"id"`
	FollowerID  int64    `json:"follower_id"`
	FollowingID int64    `json:"following_id"`
	Follower    UserInfo `json:"follower"`
	Following   UserInfo `json:"following"`
	CreatedAt   string   `json:"created_at"`
}

// FollowersListVM represents a list of users who follow a specific user
type FollowersListVM struct {
	Count     int        `json:"count"`
	Followers []UserInfo `json:"followers"`
}

// FollowingListVM represents a list of users who are followed by a specific user
type FollowingListVM struct {
	Count     int        `json:"count"`
	Following []UserInfo `json:"following"`
}

// UserInfo is a simplified user structure for follow relations
type UserInfo struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email,omitempty"`
	// Add other fields you might need
}

func (vm FollowListVM) ToViewModel(m model.Follow) FollowListVM {
	vm.ID = m.ID
	vm.FollowerID = m.FollowerID
	vm.FollowingID = m.FollowingID

	// Map follower user if loaded
	if m.Follower.ID != 0 {
		vm.Follower = UserInfo{
			ID:      m.Follower.ID,
			Name:    m.Follower.Name,
			Surname: m.Follower.Surname,
			Email:   m.Follower.Email,
		}
	}

	// Map following user if loaded
	if m.Following.ID != 0 {
		vm.Following = UserInfo{
			ID:      m.Following.ID,
			Name:    m.Following.Name,
			Surname: m.Following.Surname,
			Email:   m.Following.Email,
		}
	}

	return vm
}

func (vm FollowDetailVM) ToViewModel(m model.Follow) FollowDetailVM {
	vm.ID = m.ID
	vm.FollowerID = m.FollowerID
	vm.FollowingID = m.FollowingID
	vm.CreatedAt = m.CreatedAt.Format("2006-01-02 15:04:05")

	// Map follower user
	if m.Follower.ID != 0 {
		vm.Follower = UserInfo{
			ID:      m.Follower.ID,
			Name:    m.Follower.Name,
			Surname: m.Follower.Surname,
			Email:   m.Follower.Email,
		}
	}

	// Map following user
	if m.Following.ID != 0 {
		vm.Following = UserInfo{
			ID:      m.Following.ID,
			Name:    m.Following.Name,
			Surname: m.Following.Surname,
			Email:   m.Following.Email,
		}
	}

	return vm
}

func UserToUserInfo(user model.User) UserInfo {
	return UserInfo{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
	}
}

func (vm FollowCreateVM) ToDBModel(m model.Follow) model.Follow {
	m.FollowingID = vm.FollowingID
	return m
}

func (vm FollowUpdateVM) ToDBModel(m model.Follow) model.Follow {
	m.FollowingID = vm.FollowingID
	return m
}
