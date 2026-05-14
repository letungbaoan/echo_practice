package repositories_test

import (
	"testing"

	"echo_practice/internal/repositories"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFollowRepository_FollowUnfollow(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewFollowRepository(db)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")

	is, err := repo.IsFollowing(alice.ID, bob.ID)
	require.NoError(t, err)
	assert.False(t, is)

	require.NoError(t, repo.FollowUser(alice.ID, bob.ID))
	is, _ = repo.IsFollowing(alice.ID, bob.ID)
	assert.True(t, is)

	require.NoError(t, repo.UnfollowUser(alice.ID, bob.ID))
	is, _ = repo.IsFollowing(alice.ID, bob.ID)
	assert.False(t, is)
}

func TestFollowRepository_IsFollowingMany(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewFollowRepository(db)

	a := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	b := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	c := testutil.MakeUser(t, db, "cara", "c@b.com", "pw")
	d := testutil.MakeUser(t, db, "dan", "d@b.com", "pw")

	testutil.MakeFollow(t, db, a.ID, b.ID)
	testutil.MakeFollow(t, db, a.ID, d.ID)

	got, err := repo.IsFollowingMany(a.ID, []uint{b.ID, c.ID, d.ID})
	require.NoError(t, err)
	assert.True(t, got[b.ID])
	assert.False(t, got[c.ID])
	assert.True(t, got[d.ID])
}

func TestFollowRepository_IsFollowingMany_EmptyInputs(t *testing.T) {
	db := testutil.NewDB(t)
	repo := repositories.NewFollowRepository(db)

	got, err := repo.IsFollowingMany(0, []uint{1})
	require.NoError(t, err)
	assert.Empty(t, got)

	got, err = repo.IsFollowingMany(1, nil)
	require.NoError(t, err)
	assert.Empty(t, got)
}
