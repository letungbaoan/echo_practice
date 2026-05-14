package repositories_test

import (
	"testing"

	"echo_practice/internal/repositories"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoriteRepository_FavoriteUnfavorite(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewFavoriteRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "post")

	got, err := repo.IsFavorited(alice.ID, article.ID)
	require.NoError(t, err)
	assert.False(t, got)

	require.NoError(t, repo.Favorite(alice.ID, article.ID))
	got, _ = repo.IsFavorited(alice.ID, article.ID)
	assert.True(t, got)

	require.NoError(t, repo.Favorite(alice.ID, article.ID), "favorite must be idempotent")

	cnt, err := repo.CountByArticle(article.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, cnt)

	require.NoError(t, repo.Unfavorite(alice.ID, article.ID))
	got, _ = repo.IsFavorited(alice.ID, article.ID)
	assert.False(t, got)

	require.NoError(t, repo.Unfavorite(alice.ID, article.ID), "unfavorite must be idempotent")
}

func TestFavoriteRepository_BatchQueries(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewFavoriteRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	a1 := testutil.MakeArticle(t, db, alice.ID, "a1")
	a2 := testutil.MakeArticle(t, db, alice.ID, "a2")
	a3 := testutil.MakeArticle(t, db, alice.ID, "a3")

	testutil.MakeFavorite(t, db, alice.ID, a1.ID)
	testutil.MakeFavorite(t, db, bob.ID, a1.ID)
	testutil.MakeFavorite(t, db, alice.ID, a2.ID)

	favMap, err := repo.IsFavoritedMany(alice.ID, []uint{a1.ID, a2.ID, a3.ID})
	require.NoError(t, err)
	assert.True(t, favMap[a1.ID])
	assert.True(t, favMap[a2.ID])
	assert.False(t, favMap[a3.ID])

	counts, err := repo.CountByArticles([]uint{a1.ID, a2.ID, a3.ID})
	require.NoError(t, err)
	assert.Equal(t, 2, counts[a1.ID])
	assert.Equal(t, 1, counts[a2.ID])
	assert.Equal(t, 0, counts[a3.ID])
}

func TestFavoriteRepository_EmptyInputs(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewFavoriteRepository(db)

	favMap, err := repo.IsFavoritedMany(0, []uint{1, 2})
	require.NoError(t, err)
	assert.Empty(t, favMap)

	favMap, err = repo.IsFavoritedMany(1, nil)
	require.NoError(t, err)
	assert.Empty(t, favMap)

	counts, err := repo.CountByArticles(nil)
	require.NoError(t, err)
	assert.Empty(t, counts)
}
