package models

import "time"

type Favorite struct {
	UserID    uint      `gorm:"primaryKey;autoIncrement:false"`
	ArticleID uint      `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time

	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Article Article `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE"`
}
