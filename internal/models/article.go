package models

import "time"

type Article struct {
	ID          uint      `gorm:"primaryKey"`
	Slug        string    `gorm:"uniqueIndex;not null"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"type:text;default:''"`
	Body        string    `gorm:"type:text;default:''"`
	AuthorID    uint      `gorm:"index;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Author    User       `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Tags      []Tag      `gorm:"many2many:article_tags;"`
	Comments  []Comment  `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE"`
	Favorites []Favorite `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE"`
}
