package repositories_test

import (
	"testing"

	"echo_practice/internal/repositories"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArticleRepository_CreateAndFindBySlug(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewArticleRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	a := testutil.MakeArticle(t, db, alice.ID, "Hello World", "go", "test")

	got, err := repo.FindBySlug(a.Slug)
	require.NoError(t, err)
	assert.Equal(t, a.Title, got.Title)
	assert.Equal(t, "alice", got.Author.Username)
	assert.Len(t, got.Tags, 2)
}

func TestArticleRepository_FindBySlug_NotFound(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewArticleRepository(db)

	_, err := repo.FindBySlug("does-not-exist")
	assert.True(t, repositories.IsNotFound(err))
}

func TestArticleRepository_Delete(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewArticleRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	a := testutil.MakeArticle(t, db, alice.ID, "doomed", "tag1", "tag2")

	require.NoError(t, repo.Delete(a))
	_, err := repo.FindBySlug(a.Slug)
	assert.True(t, repositories.IsNotFound(err))
}

func TestArticleRepository_List_Filters(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewArticleRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")

	testutil.MakeArticle(t, db, alice.ID, "Go intro", "go", "tutorial")
	testutil.MakeArticle(t, db, alice.ID, "Echo intro", "go", "echo")
	rust := testutil.MakeArticle(t, db, bob.ID, "Rust intro", "rust")

	testutil.MakeFavorite(t, db, alice.ID, rust.ID)

	t.Run("no filter", func(t *testing.T) {
		arts, total, err := repo.List(repositories.ListFilter{Limit: 10})
		require.NoError(t, err)
		assert.Len(t, arts, 3)
		assert.Equal(t, int64(3), total)
	})

	t.Run("by tag", func(t *testing.T) {
		arts, total, err := repo.List(repositories.ListFilter{Tag: "go", Limit: 10})
		require.NoError(t, err)
		assert.Len(t, arts, 2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("by author", func(t *testing.T) {
		arts, total, err := repo.List(repositories.ListFilter{AuthorUsername: "bob", Limit: 10})
		require.NoError(t, err)
		assert.Len(t, arts, 1)
		assert.Equal(t, "Rust intro", arts[0].Title)
		assert.Equal(t, int64(1), total)
	})

	t.Run("by favorited", func(t *testing.T) {
		arts, total, err := repo.List(repositories.ListFilter{FavoritedUsername: "alice", Limit: 10})
		require.NoError(t, err)
		assert.Len(t, arts, 1)
		assert.Equal(t, "Rust intro", arts[0].Title)
		assert.Equal(t, int64(1), total)
	})

	t.Run("pagination", func(t *testing.T) {
		arts, total, err := repo.List(repositories.ListFilter{Limit: 1, Offset: 1})
		require.NoError(t, err)
		assert.Len(t, arts, 1)
		assert.Equal(t, int64(3), total)
	})
}

func TestArticleRepository_Feed(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewArticleRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	cara := testutil.MakeUser(t, db, "cara", "c@b.com", "pw")

	testutil.MakeArticle(t, db, alice.ID, "A1")
	testutil.MakeArticle(t, db, bob.ID, "B1")
	testutil.MakeArticle(t, db, cara.ID, "C1")

	testutil.MakeFollow(t, db, alice.ID, bob.ID)

	arts, total, err := repo.Feed(alice.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, arts, 1)
	assert.Equal(t, "B1", arts[0].Title)
	assert.Equal(t, int64(1), total)
}
