package model

import "github.com/hayrat/go-template2/backend/pkg/model"

type Post struct {
	model.BaseModel

	Title        string    `json:"title"`
	Content      string    `json:"content"`
	MainContent  string    `json:"main_content"`
	Image        string    `json:"image"`
	UserId       int64     `json:"user_id"`
	LikeCount    int64     `json:"like_count"`
	CommentCount int64     `json:"comment_count"`
	Likes        []Like    `bun:"rel:has-many,join:id=post_id,on_delete:cascade" json:"likes"`
	Comments     []Comment `bun:"rel:has-many,join:id=post_id,on_delete:cascade" json:"comments"`
}

func (Post) ModelName() string {
	return "post"
}

//TODO: created at değişkeni ekle database de bir sorun çıkarıyor. şimdilik sildim
