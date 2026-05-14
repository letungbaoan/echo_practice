package services_test

import (
	"testing"

	"echo_practice/internal/repositories"
	"echo_practice/internal/services"
	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagService_ListTags_Empty(t *testing.T) {
	db := testutil.NewDB(t)
	svc := services.NewTagService(repositories.NewTagRepository(db))

	resp, err := svc.ListTags()
	require.NoError(t, err)
	assert.NotNil(t, resp.Tags)
	assert.Empty(t, resp.Tags)
}

func TestTagService_ListTags_UsedOnly(t *testing.T) {
	db := testutil.NewDB(t)
	tagRepo := repositories.NewTagRepository(db)
	svc := services.NewTagService(tagRepo)

	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	testutil.MakeArticle(t, db, alice.ID, "a", "go", "echo")
	_, _ = tagRepo.FindOrCreate("orphan-tag")

	resp, err := svc.ListTags()
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"echo", "go"}, resp.Tags)
	assert.NotContains(t, resp.Tags, "orphan-tag")
}
