package models

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey"`
	Body      string    `gorm:"type:text;not null"`
	ArticleID uint      `gorm:"index;not null"`
	AuthorID  uint      `gorm:"index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Article Article `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE"`
	Author  User    `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
}
