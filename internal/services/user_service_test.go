package services_test

import (
	"testing"

	"echo_practice/internal/apperrors"
	"echo_practice/internal/dto"
	"echo_practice/internal/repositories"
	"echo_practice/internal/services"
	"echo_practice/internal/testutil"
	"echo_practice/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const jwtSecret = "test-secret"

func newUserService(t *testing.T) (*services.UserService, *repositories.UserRepository) {
	t.Helper()
	db := testutil.NewDB(t)
	repo := repositories.NewUserRepository(db)
	return services.NewUserService(repo, jwtSecret), repo
}

func registerReq(username, email, password string) dto.RegisterRequest {
	var r dto.RegisterRequest
	r.User.Username = username
	r.User.Email = email
	r.User.Password = password
	return r
}

func TestUserService_Register_Success(t *testing.T) {
	svc, _ := newUserService(t)

	user, token, err := svc.Register(registerReq("alice", "a@b.com", "password123"))
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.NotEmpty(t, token)

	claims, err := utils.ParseToken(token, jwtSecret)
	require.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	svc, _ := newUserService(t)
	_, _, err := svc.Register(registerReq("alice", "a@b.com", "password123"))
	require.NoError(t, err)

	_, _, err = svc.Register(registerReq("bob", "a@b.com", "password456"))
	assert.ErrorIs(t, err, apperrors.ErrEmailTaken)
}

func TestUserService_Register_DuplicateUsername(t *testing.T) {
	svc, _ := newUserService(t)
	_, _, err := svc.Register(registerReq("alice", "a@b.com", "password123"))
	require.NoError(t, err)

	_, _, err = svc.Register(registerReq("alice", "different@x.com", "password456"))
	assert.ErrorIs(t, err, apperrors.ErrUsernameTaken)
}

func TestUserService_Login(t *testing.T) {
	svc, _ := newUserService(t)
	_, _, err := svc.Register(registerReq("alice", "a@b.com", "password123"))
	require.NoError(t, err)

	var loginReq dto.LoginRequest
	loginReq.User.Email = "a@b.com"
	loginReq.User.Password = "password123"

	user, token, err := svc.Login(loginReq)
	require.NoError(t, err)
	assert.Equal(t, "alice", user.Username)
	assert.NotEmpty(t, token)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	svc, _ := newUserService(t)
	_, _, err := svc.Register(registerReq("alice", "a@b.com", "password123"))
	require.NoError(t, err)

	var loginReq dto.LoginRequest
	loginReq.User.Email = "a@b.com"
	loginReq.User.Password = "wrong"

	_, _, err = svc.Login(loginReq)
	assert.ErrorIs(t, err, apperrors.ErrInvalidLogin)
}

func TestUserService_Login_UnknownEmail(t *testing.T) {
	svc, _ := newUserService(t)
	var loginReq dto.LoginRequest
	loginReq.User.Email = "nobody@x.com"
	loginReq.User.Password = "anything"

	_, _, err := svc.Login(loginReq)
	assert.ErrorIs(t, err, apperrors.ErrInvalidLogin)
}

func TestUserService_GetCurrentUser(t *testing.T) {
	svc, _ := newUserService(t)
	u, _, err := svc.Register(registerReq("alice", "a@b.com", "password123"))
	require.NoError(t, err)

	got, err := svc.GetCurrentUser(u.ID)
	require.NoError(t, err)
	assert.Equal(t, "alice", got.Username)

	_, err = svc.GetCurrentUser(999)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}

func TestUserService_UpdateUser(t *testing.T) {
	svc, _ := newUserService(t)
	u, _, _ := svc.Register(registerReq("alice", "a@b.com", "password123"))

	var req dto.UpdateRequest
	req.User.Email = "new@b.com"
	req.User.Bio = "hello"

	got, err := svc.UpdateUser(u.ID, req)
	require.NoError(t, err)
	assert.Equal(t, "new@b.com", got.Email)
	assert.Equal(t, "hello", got.Bio)
}

func TestUserService_UpdateUser_EmailTaken(t *testing.T) {
	svc, _ := newUserService(t)
	_, _, _ = svc.Register(registerReq("alice", "alice@b.com", "password123"))
	u2, _, _ := svc.Register(registerReq("bob", "bob@b.com", "password123"))

	var req dto.UpdateRequest
	req.User.Email = "alice@b.com"

	_, err := svc.UpdateUser(u2.ID, req)
	assert.ErrorIs(t, err, apperrors.ErrEmailTaken)
}
