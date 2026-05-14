package services_test

import (
	"testing"

	"echo_practice/internal/apperrors"
	"echo_practice/internal/dto"
	"echo_practice/internal/repositories"
	"echo_practice/internal/services"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newArticleService(t *testing.T) (*services.ArticleService, *gorm.DB) {
	t.Helper()
	db := testutil.NewDB(t)
	svc := services.NewArticleService(
		repositories.NewArticleRepository(db),
		repositories.NewTagRepository(db),
		repositories.NewUserRepository(db),
		repositories.NewFollowRepository(db),
		repositories.NewFavoriteRepository(db),
	)
	return svc, db
}

func createReq(title, body string, tags []string) dto.CreateArticleRequest {
	var r dto.CreateArticleRequest
	r.Article.Title = title
	r.Article.Description = "desc"
	r.Article.Body = body
	r.Article.TagList = tags
	return r
}

func TestArticleService_CreateArticle(t *testing.T) {
	svc, db := newArticleService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")

	resp, err := svc.CreateArticle(alice.ID, createReq("Hello", "body", []string{"go", "test"}))
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Article.Slug)
	assert.Equal(t, "alice", resp.Article.Author.Username)
	assert.ElementsMatch(t, []string{"go", "test"}, resp.Article.TagList)
	assert.Equal(t, 0, resp.Article.FavoritesCount)
	assert.False(t, resp.Article.Favorited)
}

func TestArticleService_GetArticle(t *testing.T) {
	svc, db := newArticleService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "hello", "go")
	testutil.MakeFavorite(t, db, bob.ID, article.ID)
	testutil.MakeFollow(t, db, bob.ID, alice.ID)

	t.Run("anonymous", func(t *testing.T) {
		resp, err := svc.GetArticle(article.Slug, nil)
		require.NoError(t, err)
		assert.Equal(t, 1, resp.Article.FavoritesCount)
		assert.False(t, resp.Article.Favorited)
		assert.False(t, resp.Article.Author.Following)
	})

	t.Run("with bob (favorited + following)", func(t *testing.T) {
		uid := bob.ID
		resp, err := svc.GetArticle(article.Slug, &uid)
		require.NoError(t, err)
		assert.True(t, resp.Article.Favorited)
		assert.True(t, resp.Article.Author.Following)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetArticle("missing", nil)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}

func TestArticleService_UpdateArticle_Ownership(t *testing.T) {
	svc, db := newArticleService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "original")

	var req dto.UpdateArticleRequest
	req.Article.Title = "updated title"

	_, err := svc.UpdateArticle(article.Slug, bob.ID, req)
	assert.ErrorIs(t, err, apperrors.ErrForbidden, "bob is not the author")

	resp, err := svc.UpdateArticle(article.Slug, alice.ID, req)
	require.NoError(t, err)
	assert.Equal(t, "updated title", resp.Article.Title)
	assert.NotEqual(t, article.Slug, resp.Article.Slug, "slug regenerated on title change")
}

func TestArticleService_DeleteArticle_Ownership(t *testing.T) {
	svc, db := newArticleService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "post")

	err := svc.DeleteArticle(article.Slug, bob.ID)
	assert.ErrorIs(t, err, apperrors.ErrForbidden)

	require.NoError(t, svc.DeleteArticle(article.Slug, alice.ID))

	_, err = svc.GetArticle(article.Slug, nil)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestArticleService_ListArticles_Filters(t *testing.T) {
	svc, db := newArticleService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	testutil.MakeArticle(t, db, alice.ID, "go1", "go")
	testutil.MakeArticle(t, db, alice.ID, "go2", "go")
	testutil.MakeArticle(t, db, bob.ID, "rust", "rust")

	resp, err := svc.ListArticles(services.ListArticlesFilter{Tag: "go"}, nil)
	require.NoError(t, err)
	assert.Equal(t, 2, resp.ArticlesCount)
	assert.Len(t, resp.Articles, 2)

	resp, err = svc.ListArticles(services.ListArticlesFilter{AuthorUsername: "bob"}, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.ArticlesCount)
}

func TestArticleService_FeedArticles(t *testing.T) {
	svc, db := newArticleService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	cara := testutil.MakeUser(t, db, "cara", "c@b.com", "pw")
	testutil.MakeArticle(t, db, bob.ID, "bobs post")
	testutil.MakeArticle(t, db, cara.ID, "caras post")
	testutil.MakeFollow(t, db, alice.ID, bob.ID)

	resp, err := svc.FeedArticles(alice.ID, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.ArticlesCount)
	assert.Equal(t, "bobs post", resp.Articles[0].Title)
	assert.True(t, resp.Articles[0].Author.Following)
}

func TestArticleService_FavoriteUnfavorite(t *testing.T) {
	svc, db := newArticleService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "post")

	resp, err := svc.FavoriteArticle(article.Slug, bob.ID)
	require.NoError(t, err)
	assert.True(t, resp.Article.Favorited)
	assert.Equal(t, 1, resp.Article.FavoritesCount)

	resp, err = svc.FavoriteArticle(article.Slug, bob.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Article.FavoritesCount, "idempotent")

	resp, err = svc.UnfavoriteArticle(article.Slug, bob.ID)
	require.NoError(t, err)
	assert.False(t, resp.Article.Favorited)
	assert.Equal(t, 0, resp.Article.FavoritesCount)
}
