package repositories_test

import (
	"testing"

	"echo_practice/internal/models"
	"echo_practice/internal/repositories"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateAndFind(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewUserRepository(db)

	u := &models.User{Email: "a@b.com", Username: "alice", PasswordHash: "hash"}
	require.NoError(t, repo.Create(u))
	assert.NotZero(t, u.ID)

	got, err := repo.FindByEmail("a@b.com")
	require.NoError(t, err)
	assert.Equal(t, "alice", got.Username)

	got, err = repo.FindByUsername("alice")
	require.NoError(t, err)
	assert.Equal(t, u.ID, got.ID)

	got, err = repo.FindByID(u.ID)
	require.NoError(t, err)
	assert.Equal(t, "a@b.com", got.Email)
}

func TestUserRepository_NotFound(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewUserRepository(db)

	_, err := repo.FindByEmail("missing@x.com")
	assert.True(t, repositories.IsNotFound(err))

	_, err = repo.FindByUsername("nobody")
	assert.True(t, repositories.IsNotFound(err))

	_, err = repo.FindByID(999)
	assert.True(t, repositories.IsNotFound(err))
}

func TestUserRepository_Update(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewUserRepository(db)

	u := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	u.Bio = "updated bio"
	require.NoError(t, repo.Update(u))

	got, _ := repo.FindByID(u.ID)
	assert.Equal(t, "updated bio", got.Bio)
}
