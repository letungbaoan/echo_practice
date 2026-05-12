package repositories

import (
	"echo_practice/internal/models"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Favorite(userID, articleID uint) error {
	fav := models.Favorite{UserID: userID, ArticleID: articleID}
	return r.db.
		Where("user_id = ? AND article_id = ?", userID, articleID).
		FirstOrCreate(&fav).Error
}

func (r *FavoriteRepository) Unfavorite(userID, articleID uint) error {
	return r.db.
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Delete(&models.Favorite{}).Error
}

func (r *FavoriteRepository) IsFavorited(userID, articleID uint) (bool, error) {
	if userID == 0 {
		return false, nil
	}
	var count int64
	err := r.db.Model(&models.Favorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

func (r *FavoriteRepository) CountByArticle(articleID uint) (int, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return int(count), err
}

func (r *FavoriteRepository) IsFavoritedMany(userID uint, articleIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if userID == 0 || len(articleIDs) == 0 {
		return result, nil
	}
	var favs []models.Favorite
	if err := r.db.Where("user_id = ? AND article_id IN ?", userID, articleIDs).Find(&favs).Error; err != nil {
		return nil, err
	}
	for _, f := range favs {
		result[f.ArticleID] = true
	}
	return result, nil
}

func (r *FavoriteRepository) CountByArticles(articleIDs []uint) (map[uint]int, error) {
	result := make(map[uint]int)
	if len(articleIDs) == 0 {
		return result, nil
	}
	type row struct {
		ArticleID uint
		Count     int
	}
	var rows []row
	if err := r.db.Model(&models.Favorite{}).
		Select("article_id, COUNT(*) as count").
		Where("article_id IN ?", articleIDs).
		Group("article_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r.ArticleID] = r.Count
	}
	return result, nil
}
