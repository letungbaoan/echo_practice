package models

import "time"

type Follow struct {
	FollowerID  uint      `gorm:"primaryKey;autoIncrement:false"`
	FollowingID uint      `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt   time.Time

	Follower  User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE"`
	Following User `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE"`
}
