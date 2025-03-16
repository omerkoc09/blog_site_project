package model

import "github.com/hayrat/go-template2/backend/pkg/model"

type Comment struct {
	model.BaseModel

	PostId  int64  `json:"post_id"`
	UserId  int64  `json:"user_id"`
	Content string `json:"content"`
}

func (Comment) ModelName() string {
	return "comment"
}
