package migration

type Comment struct {
	BaseModel

	PostId  int64  `bun:",notnull" json:"post_id"`
	UserId  int64  `bun:",notnull" json:"user_id"`
	Content string `bun:",notnull" json:"content"`
}

func (Comment) ModelName() string {
	return "comment"
}
