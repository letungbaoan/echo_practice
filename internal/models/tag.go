package models

import "time"

type Tag struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time

	Articles []Article `gorm:"many2many:article_tags;"`
}
