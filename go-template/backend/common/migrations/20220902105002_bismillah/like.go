package migration

type Like struct {
	BaseModel

	UserId int64 `json:"user_id"`
	PostId int64 `json:"post_id"`
}

func (Like) ModelName() string {
	return "like"
}
