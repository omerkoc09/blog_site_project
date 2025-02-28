package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
)

type IAuthService interface {
	GenerateTokenPair(userID int64, refreshTokenID uuid.UUID) (model.AuthTokenPair, error)
	ParseRefreshToken(refreshToken string) (refreshTokenID uuid.UUID, userID int64, err error)
	GetAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) (model.AuthRefreshToken, error)
	CreateAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID, userID int64) error
	UpdateAuthRefreshTokenExpires(ctx context.Context, refreshToken model.AuthRefreshToken) error
	DeleteAuthRefreshToken(ctx context.Context, userID int64) error
}

type AuthService struct {
	DB                     *bun.DB
	jwtSecret              string
	accessTokenExpireTime  time.Duration
	refreshTokenExpireTime time.Duration
}

func NewAuthService(db *bun.DB, jwtSecret string, accessTokenExpireTime, refreshTokenExpireTime time.Duration) IAuthService {
	return &AuthService{
		DB:                     db,
		jwtSecret:              jwtSecret,
		accessTokenExpireTime:  accessTokenExpireTime,
		refreshTokenExpireTime: refreshTokenExpireTime,
	}
}

func (s AuthService) ParseRefreshToken(refreshToken string) (refreshTokenID uuid.UUID, userID int64, err error) {
	// Parse the refresh token.
	refreshClaims := refreshTokenClaims{}
	claims, err := jwt.ParseWithClaims(refreshToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return
	}
	if !claims.Valid {
		err = errors.New("invalid token")
		return
	}

	// Verify that the refresh token is not expired
	now := time.Now()
	rtokenClaims := claims.Claims.(*refreshTokenClaims)
	if rtokenClaims.VerifyExpiresAt(now, false) == false {
		err = errors.New("token expired")
		return
	}

	return rtokenClaims.ID, rtokenClaims.UserID, nil
}

type accessTokenClaims struct {
	jwt.RegisteredClaims
	ID     uuid.UUID `json:"id"`
	UserID int64     `json:"uid"`
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	ID     uuid.UUID `json:"id"`
	UserID int64     `json:"uid"`
}

func (s AuthService) GenerateTokenPair(userID int64, refreshTokenID uuid.UUID) (model.AuthTokenPair, error) {
	var err error
	var m model.AuthTokenPair
	now := time.Now()

	accessClaims := accessTokenClaims{
		ID:     uuid.New(),
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenExpireTime)),
		},
	}

	m.AccessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return m, err
	}

	refreshClaims := refreshTokenClaims{
		ID:     refreshTokenID,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenExpireTime)),
		},
	}

	m.RefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return m, err
	}

	return m, nil
}

func (s AuthService) GetAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) (model.AuthRefreshToken, error) {
	refreshToken := model.AuthRefreshToken{}
	err := s.DB.NewSelect().Model(&refreshToken).Where("token_id = ?", refreshTokenID).Scan(ctx)
	if err != nil {
		return refreshToken, errorsx.Database(err)
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return refreshToken, errors.New("token expired")
	}
	return refreshToken, nil
}

func (s AuthService) CreateAuthRefreshToken(ctx context.Context, refreshTokenID uuid.UUID, userID int64) error {
	refreshToken := model.AuthRefreshToken{
		TokenID:   refreshTokenID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.refreshTokenExpireTime),
	}
	_, err := s.DB.NewInsert().Model(&refreshToken).Exec(ctx)
	return errorsx.Database(err)
}

func (s AuthService) UpdateAuthRefreshTokenExpires(ctx context.Context, refreshToken model.AuthRefreshToken) error {
	_, err := s.DB.NewUpdate().Model(&model.AuthRefreshToken{}).
		Where("token_id = ?", refreshToken.TokenID).
		Set("expires_at = ?", time.Now().Add(s.refreshTokenExpireTime)).
		Exec(ctx)
	return errorsx.Database(err)
}

func (s AuthService) DeleteAuthRefreshToken(ctx context.Context, userID int64) error {
	_, err := s.DB.NewDelete().Model(&model.AuthRefreshToken{}).Where("user_id = ?", userID).Exec(ctx)
	return errorsx.Database(err)
}
