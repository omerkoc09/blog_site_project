package model

import "github.com/hayrat/go-template2/backend/pkg/model"

type Comment struct {
	model.BaseModel

	PostId  int64  `bun:",notnull" json:"post_id"`
	UserId  int64  `bun:",notnull" json:"user_id"`
	Content string `bun:",notnull" json:"content"`
}

func (Comment) ModelName() string {
	return "comment"
}
