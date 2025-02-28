package migration

import (
	"github.com/google/uuid"
)

type UserRole int

type User struct {
	BaseModel
	UUID uuid.UUID `json:"uuid" bun:"type:uuid"`

	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	About    string    `json:"about"`
	Role     UserRole  `json:"role"`
	Comments []Comment `bun:"rel:has-many,join:id=user_id,on_delete:cascade" json:"comments"`
	Likes    []Like    `bun:"rel:has-many,join:id=user_id,on_delete:cascade" json:"likes"`
}

func (User) ModelName() string {
	return "user"
}
