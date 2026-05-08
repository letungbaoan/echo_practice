package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	Username     string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Bio          string    `gorm:"type:text;default:''"`
	Image        string    `gorm:"default:''"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Articles  []Article  `gorm:"foreignKey:AuthorID"`
	Comments  []Comment  `gorm:"foreignKey:AuthorID"`
	Favorites []Favorite `gorm:"foreignKey:UserID"`

	Following []Follow `gorm:"foreignKey:FollowerID"`
	Followers []Follow `gorm:"foreignKey:FollowingID"`
}
