package repositories_test

import (
	"testing"

	"echo_practice/internal/models"
	"echo_practice/internal/repositories"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentRepository_CRUD(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewCommentRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "title")

	c1 := &models.Comment{Body: "first", ArticleID: article.ID, AuthorID: bob.ID}
	require.NoError(t, repo.Create(c1))
	assert.NotZero(t, c1.ID)

	c2 := &models.Comment{Body: "second", ArticleID: article.ID, AuthorID: alice.ID}
	require.NoError(t, repo.Create(c2))

	got, err := repo.FindByID(c1.ID)
	require.NoError(t, err)
	assert.Equal(t, "first", got.Body)
	assert.Equal(t, "bob", got.Author.Username)

	list, err := repo.ListByArticle(article.ID)
	require.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "second", list[0].Body, "DESC order: newest first")

	require.NoError(t, repo.Delete(c1.ID))
	list, err = repo.ListByArticle(article.ID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestCommentRepository_FindByID_NotFound(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewCommentRepository(db)

	_, err := repo.FindByID(999)
	assert.True(t, repositories.IsNotFound(err))
}
