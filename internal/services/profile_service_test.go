package services_test

import (
	"testing"

	"echo_practice/internal/apperrors"
	"echo_practice/internal/repositories"
	"echo_practice/internal/services"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileService_GetProfile(t *testing.T) {
	db := testutil.NewDB(t)
	svc := services.NewProfileService(
		repositories.NewUserRepository(db),
		repositories.NewFollowRepository(db),
	)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")

	t.Run("anonymous", func(t *testing.T) {
		resp, err := svc.GetProfile("alice", nil)
		require.NoError(t, err)
		assert.Equal(t, "alice", resp.Profile.Username)
		assert.False(t, resp.Profile.Following)
	})

	t.Run("bob follows alice", func(t *testing.T) {
		testutil.MakeFollow(t, db, bob.ID, alice.ID)
		resp, err := svc.GetProfile("alice", &bob.ID)
		require.NoError(t, err)
		assert.True(t, resp.Profile.Following)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := svc.GetProfile("nobody", nil)
		assert.ErrorIs(t, err, apperrors.ErrNotFound)
	})
}

func TestProfileService_FollowUnfollow(t *testing.T) {
	db := testutil.NewDB(t)
	svc := services.NewProfileService(
		repositories.NewUserRepository(db),
		repositories.NewFollowRepository(db),
	)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")

	resp, err := svc.FollowUser(bob.ID, "alice")
	require.NoError(t, err)
	assert.True(t, resp.Profile.Following)

	resp, err = svc.UnfollowUser(bob.ID, "alice")
	require.NoError(t, err)
	assert.False(t, resp.Profile.Following)

	_ = alice
}
