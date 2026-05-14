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

func newCommentService(t *testing.T) (*services.CommentService, *gorm.DB) {
	t.Helper()
	db := testutil.NewDB(t)
	svc := services.NewCommentService(
		repositories.NewCommentRepository(db),
		repositories.NewArticleRepository(db),
		repositories.NewUserRepository(db),
		repositories.NewFollowRepository(db),
	)
	return svc, db
}

func commentReq(body string) dto.CreateCommentRequest {
	var r dto.CreateCommentRequest
	r.Comment.Body = body
	return r
}

func TestCommentService_AddComment(t *testing.T) {
	svc, db := newCommentService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "title")

	resp, err := svc.AddComment(article.Slug, bob.ID, commentReq("first!"))
	require.NoError(t, err)
	assert.Equal(t, "first!", resp.Comment.Body)
	assert.Equal(t, "bob", resp.Comment.Author.Username)
}

func TestCommentService_AddComment_ArticleNotFound(t *testing.T) {
	svc, db := newCommentService(t)
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")

	_, err := svc.AddComment("missing", bob.ID, commentReq("hi"))
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestCommentService_ListComments_WithFollowing(t *testing.T) {
	svc, db := newCommentService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "post")

	_, _ = svc.AddComment(article.Slug, bob.ID, commentReq("b"))
	_, _ = svc.AddComment(article.Slug, alice.ID, commentReq("a"))

	testutil.MakeFollow(t, db, bob.ID, alice.ID)

	uid := bob.ID
	resp, err := svc.ListComments(article.Slug, &uid)
	require.NoError(t, err)
	assert.Len(t, resp.Comments, 2)

	for _, c := range resp.Comments {
		if c.Author.Username == "alice" {
			assert.True(t, c.Author.Following)
		}
		if c.Author.Username == "bob" {
			assert.False(t, c.Author.Following)
		}
	}
}

func TestCommentService_DeleteComment_Ownership(t *testing.T) {
	svc, db := newCommentService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	bob := testutil.MakeUser(t, db, "bob", "b@b.com", "pw")
	article := testutil.MakeArticle(t, db, alice.ID, "post")

	resp, _ := svc.AddComment(article.Slug, bob.ID, commentReq("bob's comment"))

	err := svc.DeleteComment(article.Slug, resp.Comment.ID, alice.ID)
	assert.ErrorIs(t, err, apperrors.ErrForbidden)

	require.NoError(t, svc.DeleteComment(article.Slug, resp.Comment.ID, bob.ID))
}

func TestCommentService_DeleteComment_WrongArticle(t *testing.T) {
	svc, db := newCommentService(t)
	alice := testutil.MakeUser(t, db, "alice", "a@b.com", "pw")
	a1 := testutil.MakeArticle(t, db, alice.ID, "one")
	a2 := testutil.MakeArticle(t, db, alice.ID, "two")

	resp, _ := svc.AddComment(a1.Slug, alice.ID, commentReq("hi"))

	err := svc.DeleteComment(a2.Slug, resp.Comment.ID, alice.ID)
	assert.ErrorIs(t, err, apperrors.ErrNotFound, "comment belongs to a different article")
}
