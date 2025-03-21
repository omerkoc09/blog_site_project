package service

import (
	"context"
	"errors"

	"github.com/hayrat/go-template2/backend/common/model"
	"github.com/hayrat/go-template2/backend/pkg/service"
	"github.com/uptrace/bun"
)

type IFollowService interface {
	Follow(ctx context.Context, followerID, followingID int64) (model.Follow, error)
	Unfollow(ctx context.Context, followerID, followingID int64) error
	GetFollowers(ctx context.Context, userID int64) ([]model.User, error)
	GetFollowing(ctx context.Context, userID int64) ([]model.User, error)
	IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error)
}

type FollowService struct {
	service.BaseService[model.Follow]
}

func NewFollowService(db *bun.DB) IFollowService {
	return &FollowService{
		service.BaseService[model.Follow]{DB: db},
	}
}

func (s *FollowService) Follow(ctx context.Context, followerID, followingID int64) (model.Follow, error) {
	if followerID == followingID {
		return model.Follow{}, errors.New("you cannot follow yourself")
	}

	isfollow, err := s.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return model.Follow{}, err
	}
	if isfollow {
		return model.Follow{}, errors.New("you are already following this user")
	}
	follow := model.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	err = s.DB.NewInsert().Model(&follow).Scan(ctx)
	if err != nil {
		return model.Follow{}, err
	}
	return follow, nil
}
func (s *FollowService) Unfollow(ctx context.Context, followerID, followingID int64) error {
	_, err := s.DB.NewDelete().
		Model(&model.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Exec(ctx)
	return err
}

func (s *FollowService) IsFollowing(ctx context.Context, followerID, followingID int64) (bool, error) {
	var follow model.Follow
	exists, err := s.DB.NewSelect().
		Model(&follow).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Exists(ctx)

	return exists, err
}

func (s *FollowService) GetFollowers(ctx context.Context, userID int64) ([]model.User, error) {
	var follows []model.Follow
	err := s.DB.NewSelect().
		Model(&follows).
		Column("follower_id").
		Where("following_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if len(follows) == 0 {
		return []model.User{}, nil
	}

	// Extract follower IDs
	followerIDs := make([]int64, len(follows))
	for i, follow := range follows {
		followerIDs[i] = follow.FollowerID
	}

	// Query users based on extracted IDs
	var users []model.User
	err = s.DB.NewSelect().
		Model(&users).
		Where("id IN (?)", bun.In(followerIDs)).
		Scan(ctx)

	return users, err
}

func (s *FollowService) GetFollowing(ctx context.Context, userID int64) ([]model.User, error) {
	var follows []model.Follow
	err := s.DB.NewSelect().
		Model(&follows).
		Column("following_id").
		Where("follower_id = ?", userID).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if len(follows) == 0 {
		return []model.User{}, nil
	}

	// Extract following IDs
	followingIDs := make([]int64, len(follows))
	for i, follow := range follows {
		followingIDs[i] = follow.FollowingID
	}

	// Query users based on extracted IDs
	var users []model.User
	err = s.DB.NewSelect().
		Model(&users).
		Where("id IN (?)", bun.In(followingIDs)).
		Scan(ctx)

	return users, err
}
