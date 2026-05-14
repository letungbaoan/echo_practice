package testutil

import (
	"testing"

	"echo_practice/internal/database"
	"echo_practice/internal/models"
	"echo_practice/internal/utils"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.Exec("PRAGMA foreign_keys = ON").Error)
	require.NoError(t, database.Migrate(db))
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})
	return db
}

func MakeUser(t *testing.T, db *gorm.DB, username, email, password string) *models.User {
	t.Helper()
	hash, err := utils.HashPassword(password)
	require.NoError(t, err)
	u := &models.User{Username: username, Email: email, PasswordHash: hash}
	require.NoError(t, db.Create(u).Error)
	return u
}

func MakeArticle(t *testing.T, db *gorm.DB, authorID uint, title string, tags ...string) *models.Article {
	t.Helper()
	tagModels := make([]models.Tag, 0, len(tags))
	for _, name := range tags {
		var tag models.Tag
		err := db.Where("name = ?", name).First(&tag).Error
		if err != nil {
			tag = models.Tag{Name: name}
			require.NoError(t, db.Create(&tag).Error)
		}
		tagModels = append(tagModels, tag)
	}
	a := &models.Article{
		Slug:        utils.GenerateSlug(title),
		Title:       title,
		Description: "desc",
		Body:        "body",
		AuthorID:    authorID,
		Tags:        tagModels,
	}
	require.NoError(t, db.Create(a).Error)
	return a
}

func MakeFollow(t *testing.T, db *gorm.DB, followerID, followingID uint) {
	t.Helper()
	require.NoError(t, db.Create(&models.Follow{FollowerID: followerID, FollowingID: followingID}).Error)
}

func MakeFavorite(t *testing.T, db *gorm.DB, userID, articleID uint) {
	t.Helper()
	require.NoError(t, db.Create(&models.Favorite{UserID: userID, ArticleID: articleID}).Error)
}
