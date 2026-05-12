package repositories

import (
	"echo_practice/internal/models"

	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) FollowUser(followerID, followingID uint) error {
	follow := models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return r.db.Create(&follow).Error
}

func (r *FollowRepository) UnfollowUser(followerID, followingID uint) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follow{}).Error
}

func (r *FollowRepository) IsFollowing(followerID, followingID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error
	return count > 0, err
}

func (r *FollowRepository) IsFollowingMany(followerID uint, followingIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if followerID == 0 || len(followingIDs) == 0 {
		return result, nil
	}
	var follows []models.Follow
	if err := r.db.Where("follower_id = ? AND following_id IN ?", followerID, followingIDs).Find(&follows).Error; err != nil {
		return nil, err
	}
	for _, f := range follows {
		result[f.FollowingID] = true
	}
	return result, nil
}
