package repositories

import (
	"echo_practice/internal/models"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) FindOrCreate(name string) (*models.Tag, error) {
	var tag models.Tag
	result := r.db.Where("name = ?", name).First(&tag)
	if result.Error == nil {
		return &tag, nil
	}
	if !IsNotFound(result.Error) {
		return nil, result.Error
	}
	tag = models.Tag{Name: name}
	if err := r.db.Create(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) FindByNames(names []string) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Where("name IN ?", names).Find(&tags).Error
	return tags, err
}

func (r *TagRepository) ListUsed() ([]string, error) {
	var names []string
	err := r.db.Model(&models.Tag{}).
		Distinct("tags.name").
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Order("tags.name ASC").
		Pluck("name", &names).Error
	return names, err
}
