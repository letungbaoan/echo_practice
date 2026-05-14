package repositories_test

import (
	"testing"

	"echo_practice/internal/repositories"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagRepository_FindOrCreate(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewTagRepository(db)

	tag1, err := repo.FindOrCreate("go")
	require.NoError(t, err)
	assert.NotZero(t, tag1.ID)

	tag2, err := repo.FindOrCreate("go")
	require.NoError(t, err)
	assert.Equal(t, tag1.ID, tag2.ID, "second call must reuse existing tag")
}

func TestTagRepository_ListUsed(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewTagRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	testutil.MakeArticle(t, db, alice.ID, "article one", "go", "tutorial")
	testutil.MakeArticle(t, db, alice.ID, "article two", "echo")
	_, _ = repo.FindOrCreate("unused-tag")

	names, err := repo.ListUsed()
	require.NoError(t, err)
	assert.Equal(t, []string{"echo", "go", "tutorial"}, names)
	assert.NotContains(t, names, "unused-tag")
}
