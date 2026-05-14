package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"echo_practice/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret"

func newContext(t *testing.T, authHeader string) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestJWTAuth_NoHeader(t *testing.T) {
	c, _ := newContext(t, "")
	called := false
	err := JWTAuth(testSecret)(func(c echo.Context) error { called = true; return nil })(c)
	assert.Error(t, err)
	assert.False(t, called)
}

func TestJWTAuth_WrongScheme(t *testing.T) {
	c, _ := newContext(t, "Bearer abc")
	err := JWTAuth(testSecret)(func(c echo.Context) error { return nil })(c)
	assert.Error(t, err)
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	c, _ := newContext(t, "Token invalid.jwt.string")
	err := JWTAuth(testSecret)(func(c echo.Context) error { return nil })(c)
	assert.Error(t, err)
}

func TestJWTAuth_ValidToken(t *testing.T) {
	token, err := utils.GenerateToken(42, testSecret)
	require.NoError(t, err)

	c, _ := newContext(t, "Token "+token)
	called := false
	err = JWTAuth(testSecret)(func(c echo.Context) error {
		called = true
		uid, gErr := GetUserID(c)
		assert.NoError(t, gErr)
		assert.Equal(t, uint(42), uid)
		return nil
	})(c)
	require.NoError(t, err)
	assert.True(t, called, "next handler should run")
}

func TestOptionalJWTAuth_NoHeader(t *testing.T) {
	c, _ := newContext(t, "")
	err := OptionalJWTAuth(testSecret)(func(c echo.Context) error {
		assert.Nil(t, CurrentUserID(c))
		return nil
	})(c)
	assert.NoError(t, err)
}

func TestOptionalJWTAuth_ValidToken(t *testing.T) {
	token, err := utils.GenerateToken(7, testSecret)
	require.NoError(t, err)

	c, _ := newContext(t, "Token "+token)
	err = OptionalJWTAuth(testSecret)(func(c echo.Context) error {
		uid := CurrentUserID(c)
		require.NotNil(t, uid)
		assert.Equal(t, uint(7), *uid)
		return nil
	})(c)
	assert.NoError(t, err)
}

func TestOptionalJWTAuth_InvalidToken_StillPasses(t *testing.T) {
	c, _ := newContext(t, "Token bad")
	called := false
	err := OptionalJWTAuth(testSecret)(func(c echo.Context) error {
		called = true
		assert.Nil(t, CurrentUserID(c))
		return nil
	})(c)
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestGetUserID_Missing(t *testing.T) {
	c, _ := newContext(t, "")
	_, err := GetUserID(c)
	assert.Error(t, err)
}
