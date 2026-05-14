package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"echo_practice/internal/controllers"
	"echo_practice/internal/middlewares"
	"echo_practice/internal/repositories"
	"echo_practice/internal/routes"
	"echo_practice/internal/services"
	"echo_practice/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const TestJWTSecret = "integration-test-secret"

type App struct {
	Echo *echo.Echo
	DB   *gorm.DB
}

func NewApp(t *testing.T) *App {
	t.Helper()
	db := NewDB(t)

	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo, TestJWTSecret)
	userCtrl := controllers.NewUserController(userSvc)

	followRepo := repositories.NewFollowRepository(db)
	profileSvc := services.NewProfileService(userRepo, followRepo)
	profileCtrl := controllers.NewProfileController(profileSvc)

	articleRepo := repositories.NewArticleRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	favoriteRepo := repositories.NewFavoriteRepository(db)
	articleSvc := services.NewArticleService(articleRepo, tagRepo, userRepo, followRepo, favoriteRepo)
	articleCtrl := controllers.NewArticleController(articleSvc)

	commentRepo := repositories.NewCommentRepository(db)
	commentSvc := services.NewCommentService(commentRepo, articleRepo, userRepo, followRepo)
	commentCtrl := controllers.NewCommentController(commentSvc)

	tagSvc := services.NewTagService(tagRepo)
	tagCtrl := controllers.NewTagController(tagSvc)

	e := echo.New()
	e.Validator = utils.NewValidator()
	e.HTTPErrorHandler = middlewares.ErrorHandler

	routes.Register(e, routes.Deps{
		UserController:    userCtrl,
		ProfileController: profileCtrl,
		ArticleController: articleCtrl,
		CommentController: commentCtrl,
		TagController:     tagCtrl,
		JWTSecret:         TestJWTSecret,
	})

	return &App{Echo: e, DB: db}
}

func (a *App) Request(t *testing.T, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(buf)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Token "+token)
	}
	rec := httptest.NewRecorder()
	a.Echo.ServeHTTP(rec, req)
	return rec
}

func DecodeJSON(t *testing.T, rec *httptest.ResponseRecorder, dest any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), dest))
}

func RegisterUser(t *testing.T, app *App, username, email, password string) string {
	t.Helper()
	rec := app.Request(t, http.MethodPost, "/api/users", map[string]any{
		"user": map[string]string{
			"username": username,
			"email":    email,
			"password": password,
		},
	}, "")
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	var resp struct {
		User struct {
			Token string `json:"token"`
		} `json:"user"`
	}
	DecodeJSON(t, rec, &resp)
	return resp.User.Token
}
