package controllers_test

import (
	"net/http"
	"strconv"
	"testing"

	"echo_practice/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Register_Login_Flow(t *testing.T) {
	app := testutil.NewApp(t)

	t.Run("register success", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/users", map[string]any{
			"user": map[string]string{"username": "alice", "email": "a@b.com", "password": "password123"},
		}, "")
		require.Equal(t, http.StatusCreated, rec.Code)
		var resp struct {
			User struct {
				Email    string `json:"email"`
				Username string `json:"username"`
				Token    string `json:"token"`
			} `json:"user"`
		}
		testutil.DecodeJSON(t, rec, &resp)
		assert.Equal(t, "alice", resp.User.Username)
		assert.Equal(t, "a@b.com", resp.User.Email)
		assert.NotEmpty(t, resp.User.Token)
	})

	t.Run("register duplicate email → 422", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/users", map[string]any{
			"user": map[string]string{"username": "alice2", "email": "a@b.com", "password": "password123"},
		}, "")
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Contains(t, rec.Body.String(), "email already taken")
	})

	t.Run("login wrong password → 401", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/users/login", map[string]any{
			"user": map[string]string{"email": "a@b.com", "password": "wrong"},
		}, "")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("login ok", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/users/login", map[string]any{
			"user": map[string]string{"email": "a@b.com", "password": "password123"},
		}, "")
		require.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			User struct{ Token string } `json:"user"`
		}
		testutil.DecodeJSON(t, rec, &resp)
		assert.NotEmpty(t, resp.User.Token)
	})
}

func TestUser_GetCurrentUser(t *testing.T) {
	app := testutil.NewApp(t)
	token := testutil.RegisterUser(t, app, "alice", "a@b.com", "password123")

	t.Run("no token → 401", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/user", nil, "")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("with token → 200", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/user", nil, token)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"username":"alice"`)
	})
}

func TestProfile_FollowFlow(t *testing.T) {
	app := testutil.NewApp(t)
	_ = testutil.RegisterUser(t, app, "alice", "a@b.com", "password123")
	bobToken := testutil.RegisterUser(t, app, "bob", "b@b.com", "password123")

	t.Run("get profile anonymous", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/profiles/alice", nil, "")
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"following":false`)
	})

	t.Run("bob follows alice", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/profiles/alice/follow", nil, bobToken)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"following":true`)
	})

	t.Run("get profile with bob token", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/profiles/alice", nil, bobToken)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"following":true`)
	})

	t.Run("unfollow", func(t *testing.T) {
		rec := app.Request(t, http.MethodDelete, "/api/profiles/alice/follow", nil, bobToken)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"following":false`)
	})

	t.Run("profile not found → 404", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/profiles/nobody", nil, "")
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

type articleResp struct {
	Article struct {
		Slug           string   `json:"slug"`
		Title          string   `json:"title"`
		TagList        []string `json:"tagList"`
		Favorited      bool     `json:"favorited"`
		FavoritesCount int      `json:"favoritesCount"`
		Author         struct {
			Username  string `json:"username"`
			Following bool   `json:"following"`
		} `json:"author"`
	} `json:"article"`
}

func TestArticle_CRUDFlow(t *testing.T) {
	app := testutil.NewApp(t)
	aliceToken := testutil.RegisterUser(t, app, "alice", "a@b.com", "password123")
	bobToken := testutil.RegisterUser(t, app, "bob", "b@b.com", "password123")

	var created articleResp

	t.Run("create requires auth", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/articles", map[string]any{
			"article": map[string]any{"title": "x", "description": "y", "body": "z"},
		}, "")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("create article", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/articles", map[string]any{
			"article": map[string]any{
				"title":       "Hello Go",
				"description": "intro",
				"body":        "content",
				"tagList":     []string{"go", "tutorial"},
			},
		}, aliceToken)
		require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
		testutil.DecodeJSON(t, rec, &created)
		assert.ElementsMatch(t, []string{"go", "tutorial"}, created.Article.TagList)
		assert.Equal(t, "alice", created.Article.Author.Username)
	})

	t.Run("get article", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/articles/"+created.Article.Slug, nil, "")
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("update by non-owner → 403", func(t *testing.T) {
		rec := app.Request(t, http.MethodPut, "/api/articles/"+created.Article.Slug, map[string]any{
			"article": map[string]any{"title": "hacked"},
		}, bobToken)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("update by owner", func(t *testing.T) {
		rec := app.Request(t, http.MethodPut, "/api/articles/"+created.Article.Slug, map[string]any{
			"article": map[string]any{"title": "Hello Go v2"},
		}, aliceToken)
		require.Equal(t, http.StatusOK, rec.Code)
		var updated articleResp
		testutil.DecodeJSON(t, rec, &updated)
		assert.Equal(t, "Hello Go v2", updated.Article.Title)
		created.Article.Slug = updated.Article.Slug
	})

	t.Run("list articles filter by author", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/articles?author=alice", nil, "")
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"articlesCount":1`)
	})

	t.Run("delete by non-owner → 403", func(t *testing.T) {
		rec := app.Request(t, http.MethodDelete, "/api/articles/"+created.Article.Slug, nil, bobToken)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("delete by owner", func(t *testing.T) {
		rec := app.Request(t, http.MethodDelete, "/api/articles/"+created.Article.Slug, nil, aliceToken)
		assert.Equal(t, http.StatusOK, rec.Code)

		rec = app.Request(t, http.MethodGet, "/api/articles/"+created.Article.Slug, nil, "")
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestArticle_Feed(t *testing.T) {
	app := testutil.NewApp(t)
	aliceToken := testutil.RegisterUser(t, app, "alice", "a@b.com", "password123")
	bobToken := testutil.RegisterUser(t, app, "bob", "b@b.com", "password123")

	app.Request(t, http.MethodPost, "/api/articles", map[string]any{
		"article": map[string]any{"title": "Alice post", "description": "d", "body": "b"},
	}, aliceToken)

	t.Run("bob feed empty before follow", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/articles/feed", nil, bobToken)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"articlesCount":0`)
	})

	app.Request(t, http.MethodPost, "/api/profiles/alice/follow", nil, bobToken)

	t.Run("bob feed has alice post", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/articles/feed", nil, bobToken)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"articlesCount":1`)
		assert.Contains(t, rec.Body.String(), `"following":true`)
	})
}

func TestComment_Flow(t *testing.T) {
	app := testutil.NewApp(t)
	aliceToken := testutil.RegisterUser(t, app, "alice", "a@b.com", "password123")
	bobToken := testutil.RegisterUser(t, app, "bob", "b@b.com", "password123")

	rec := app.Request(t, http.MethodPost, "/api/articles", map[string]any{
		"article": map[string]any{"title": "Post", "description": "d", "body": "b"},
	}, aliceToken)
	var ar articleResp
	testutil.DecodeJSON(t, rec, &ar)
	slug := ar.Article.Slug

	var commentID float64

	t.Run("bob add comment", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/articles/"+slug+"/comments", map[string]any{
			"comment": map[string]string{"body": "first"},
		}, bobToken)
		require.Equal(t, http.StatusCreated, rec.Code)
		var resp struct {
			Comment struct {
				ID   float64 `json:"id"`
				Body string  `json:"body"`
			} `json:"comment"`
		}
		testutil.DecodeJSON(t, rec, &resp)
		commentID = resp.Comment.ID
		assert.Equal(t, "first", resp.Comment.Body)
	})

	t.Run("list comments", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/articles/"+slug+"/comments", nil, "")
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"body":"first"`)
	})

	t.Run("alice cannot delete bob's comment", func(t *testing.T) {
		path := "/api/articles/" + slug + "/comments/" + intToStr(commentID)
		rec := app.Request(t, http.MethodDelete, path, nil, aliceToken)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("bob deletes own comment", func(t *testing.T) {
		path := "/api/articles/" + slug + "/comments/" + intToStr(commentID)
		rec := app.Request(t, http.MethodDelete, path, nil, bobToken)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestFavorite_Flow(t *testing.T) {
	app := testutil.NewApp(t)
	aliceToken := testutil.RegisterUser(t, app, "alice", "a@b.com", "password123")
	bobToken := testutil.RegisterUser(t, app, "bob", "b@b.com", "password123")

	rec := app.Request(t, http.MethodPost, "/api/articles", map[string]any{
		"article": map[string]any{"title": "Post", "description": "d", "body": "b"},
	}, aliceToken)
	var ar articleResp
	testutil.DecodeJSON(t, rec, &ar)
	slug := ar.Article.Slug

	t.Run("favorite", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/articles/"+slug+"/favorite", nil, bobToken)
		require.Equal(t, http.StatusOK, rec.Code)
		var r articleResp
		testutil.DecodeJSON(t, rec, &r)
		assert.True(t, r.Article.Favorited)
		assert.Equal(t, 1, r.Article.FavoritesCount)
	})

	t.Run("favorite again is idempotent", func(t *testing.T) {
		rec := app.Request(t, http.MethodPost, "/api/articles/"+slug+"/favorite", nil, bobToken)
		require.Equal(t, http.StatusOK, rec.Code)
		var r articleResp
		testutil.DecodeJSON(t, rec, &r)
		assert.Equal(t, 1, r.Article.FavoritesCount)
	})

	t.Run("filter favorited=bob", func(t *testing.T) {
		rec := app.Request(t, http.MethodGet, "/api/articles?favorited=bob", nil, "")
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"articlesCount":1`)
	})

	t.Run("unfavorite", func(t *testing.T) {
		rec := app.Request(t, http.MethodDelete, "/api/articles/"+slug+"/favorite", nil, bobToken)
		require.Equal(t, http.StatusOK, rec.Code)
		var r articleResp
		testutil.DecodeJSON(t, rec, &r)
		assert.False(t, r.Article.Favorited)
		assert.Equal(t, 0, r.Article.FavoritesCount)
	})
}

func TestTags_Endpoint(t *testing.T) {
	app := testutil.NewApp(t)
	aliceToken := testutil.RegisterUser(t, app, "alice", "a@b.com", "password123")

	app.Request(t, http.MethodPost, "/api/articles", map[string]any{
		"article": map[string]any{"title": "A", "description": "d", "body": "b", "tagList": []string{"go", "echo"}},
	}, aliceToken)

	rec := app.Request(t, http.MethodGet, "/api/tags", nil, "")
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"tags"`)
	assert.Contains(t, rec.Body.String(), `"go"`)
	assert.Contains(t, rec.Body.String(), `"echo"`)
}

func TestErrorFormat_Standard(t *testing.T) {
	app := testutil.NewApp(t)

	rec := app.Request(t, http.MethodGet, "/api/articles/nope", nil, "")
	require.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), `"errors":{"body":["not found"]}`)
}

func intToStr(f float64) string {
	return strconv.Itoa(int(f))
}
