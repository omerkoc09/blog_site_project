package model

import "github.com/hayrat/go-template2/backend/pkg/model"

type Like struct {
	model.BaseModel

	UserId int64 `json:"user_id"`
	PostId int64 `json:"post_id"`
}

func (Like) ModelName() string {
	return "like"
}
