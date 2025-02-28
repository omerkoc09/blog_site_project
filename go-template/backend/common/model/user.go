package model

import "github.com/hayrat/go-template2/backend/pkg/model"

type UserRole int

const (
	UserRoleNormal UserRole = 1
	UserRoleAdmin  UserRole = 10
)

type User struct {
	model.BaseModel

	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	About    string    `json:"about"`
	Role     UserRole  `json:"role"`
	Posts    []Post    `bun:"rel:has-many,join:id=user_id,on_delete:cascade" json:"posts"`
	Comments []Comment `bun:"rel:has-many,join:id=user_id,on_delete:cascade" json:"comments"`
	Likes    []Like    `bun:"rel:has-many,join:id=user_id,on_delete:cascade" json:"likes"`
}

func (User) ModelName() string {
	return "user"
}

func (k User) String() string {
	return k.Name + " " + k.Surname
}

func (r UserRole) String() string {
	switch r {
	case UserRoleNormal:
		return "normal"
	case UserRoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}
