package model

import "github.com/hayrat/go-template2/backend/pkg/model"

type Post struct {
	model.BaseModel

	Title       string    `json:"title"`
	Content     string    `json:"content"`
	MainContent string    `json:"main_content"`
	Image       string    `json:"image"`
	UserId      int64     `json:"user_id"`
	Likes       []Like    `bun:"rel:has-many,join:id=post_id,on_delete:cascade" json:"likes"`
	Comments    []Comment `bun:"rel:has-many,join:id=post_id,on_delete:cascade" json:"comments"`
}

func (Post) ModelName() string {
	return "post"
}
