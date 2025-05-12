package migration

type Post struct {
	BaseModel

	Title       string    `json:"title"`
	Content     string    `json:"content"`
	MainContent string    `json:"main_content"`
	Image       string    `json:"image"`
	UserID      int64     `json:"user_id" gorm:"foreignKey:PostId"`
	Likes       []Like    `bun:"rel:has-many,join:id=post_id" json:"likes"`
	Comments    []Comment `bun:"rel:has-many,join:id=post_id" json:"comments"`
	Topics      []Topic   `bun:"m2m:post_topic" json:"topics"`
}

func (Post) ModelName() string {
	return "post"
}
