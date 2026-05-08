package services

import (
	"echo_practice/internal/apperrors"
	"echo_practice/internal/dto"
	"echo_practice/internal/models"
	"echo_practice/internal/repositories"
	"echo_practice/internal/utils"
)

type UserService struct {
	repo      *repositories.UserRepository
	jwtSecret string
}

func NewUserService(repo *repositories.UserRepository, jwtSecret string) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret}
}

func (s *UserService) Register(req dto.RegisterRequest) (*models.User, string, error) {
	if existing, err := s.repo.FindByEmail(req.User.Email); err == nil && existing != nil {
		return nil, "", apperrors.ErrEmailTaken
	} else if err != nil && !repositories.IsNotFound(err) {
		return nil, "", err
	}

	if existing, err := s.repo.FindByUsername(req.User.Username); err == nil && existing != nil {
		return nil, "", apperrors.ErrUsernameTaken
	} else if err != nil && !repositories.IsNotFound(err) {
		return nil, "", err
	}

	hash, err := utils.HashPassword(req.User.Password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Email:        req.User.Email,
		Username:     req.User.Username,
		PasswordHash: hash,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (s *UserService) Login(req dto.LoginRequest) (*models.User, string, error) {
	user, err := s.repo.FindByEmail(req.User.Email)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, "", apperrors.ErrInvalidLogin
		}
		return nil, "", err
	}

	if !utils.CheckPassword(req.User.Password, user.PasswordHash) {
		return nil, "", apperrors.ErrInvalidLogin
	}

	token, err := utils.GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func ToUserResponse(u *models.User, token string) dto.UserResponse {
	return dto.UserResponse{
		User: dto.UserPayload{
			Email:    u.Email,
			Username: u.Username,
			Token:    token,
			Bio:      u.Bio,
			Image:    u.Image,
		},
	}
}
