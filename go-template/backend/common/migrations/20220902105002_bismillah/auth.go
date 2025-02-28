package migration

import (
	"time"

	"github.com/google/uuid"
)

type AuthRefreshToken struct {
	TokenID   uuid.UUID `json:"token_id" bun:"type:uuid"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (AuthRefreshToken) ModelName() string {
	return "user_refresh_token"
}
