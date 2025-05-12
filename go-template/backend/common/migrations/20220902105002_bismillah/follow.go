package migration

type Follow struct {
	BaseModel

	FollowerID  int64 `json:"follower_id"`
	FollowingID int64 `json:"following_id"`
	Follower    User  `rel:"belongs-to" json:"follower"`
	Following   User  `rel:"belongs-to" json:"following"`
}

func (Follow) ModelName() string {
	return "follows"
}

func (Follow) TableName() string {
	return "follows"
}
