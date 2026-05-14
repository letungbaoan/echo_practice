package repositories

import (
	"echo_practice/internal/models"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *ArticleRepository) FindBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Author").Preload("Tags").First(&article, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *ArticleRepository) Delete(article *models.Article) error {
	return r.db.Select("Tags").Delete(article).Error
}

func (r *ArticleRepository) FindByID(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Author").Preload("Tags").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

type ListFilter struct {
	Tag               string
	AuthorUsername    string
	FavoritedUsername string
	Limit             int
	Offset            int
}

func (r *ArticleRepository) List(filter ListFilter) ([]models.Article, int64, error) {
	q := r.db.Model(&models.Article{})

	if filter.Tag != "" {
		q = q.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.name = ?", filter.Tag)
	}
	if filter.AuthorUsername != "" {
		q = q.Joins("JOIN users AS authors ON authors.id = articles.author_id").
			Where("authors.username = ?", filter.AuthorUsername)
	}
	if filter.FavoritedUsername != "" {
		q = q.Joins("JOIN favorites ON favorites.article_id = articles.id").
			Joins("JOIN users AS fav_users ON fav_users.id = favorites.user_id").
			Where("fav_users.username = ?", filter.FavoritedUsername)
	}

	session := q.Session(&gorm.Session{})

	var total int64
	if err := session.Distinct("articles.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []models.Article
	if err := session.
		Preload("Author").Preload("Tags").
		Order("articles.created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

func (r *ArticleRepository) Feed(userID uint, limit, offset int) ([]models.Article, int64, error) {
	q := r.db.Model(&models.Article{}).
		Joins("JOIN follows ON follows.following_id = articles.author_id").
		Where("follows.follower_id = ?", userID)

	session := q.Session(&gorm.Session{})

	var total int64
	if err := session.Distinct("articles.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var articles []models.Article
	if err := session.
		Preload("Author").Preload("Tags").
		Order("articles.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}
