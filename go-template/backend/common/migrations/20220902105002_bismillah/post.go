package migration

import "github.com/hayrat/go-template2/backend/common/model"

type Post struct {
	BaseModel

	Title        string          `json:"title"`
	Content      string          `json:"content"`
	MainContent  string          `json:"main_content"`
	Image        string          `json:"image"`
	UserID       int64           `json:"user_id" gorm:"foreignKey:PostId"`
	LikeCount    int64           `json:"like_count"`
	CommentCount int64           `json:"comment_count"`
	Likes        []model.Like    `bun:"rel:has-many,join:id=post_id" json:"likes"`
	Comments     []model.Comment `bun:"rel:has-many,join:id=post_id" json:"comments"`
}

func (Post) ModelName() string {
	return "post"
}
