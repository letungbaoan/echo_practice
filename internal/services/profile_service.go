package services

import (
	"echo_practice/internal/apperrors"
	"echo_practice/internal/dto"
	"echo_practice/internal/repositories"
)

type ProfileService struct {
	userRepo   *repositories.UserRepository
	followRepo *repositories.FollowRepository
}

func NewProfileService(userRepo *repositories.UserRepository, followRepo *repositories.FollowRepository) *ProfileService {
	return &ProfileService{
		userRepo:   userRepo,
		followRepo: followRepo,
	}
}

func (s *ProfileService) GetProfile(username string, currentUserID *uint) (*dto.ProfileResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	following := false
	if currentUserID != nil {
		following, _ = s.followRepo.IsFollowing(*currentUserID, user.ID)
	}

	return &dto.ProfileResponse{
		Profile: dto.ProfilePayload{
			Username:  user.Username,
			Bio:       user.Bio,
			Image:     user.Image,
			Following: following,
		},
	}, nil
}

func (s *ProfileService) FollowUser(followerID uint, username string) (*dto.ProfileResponse, error) {
	target, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	_ = s.followRepo.FollowUser(followerID, target.ID)

	return &dto.ProfileResponse{
		Profile: dto.ProfilePayload{
			Username:  target.Username,
			Bio:       target.Bio,
			Image:     target.Image,
			Following: true,
		},
	}, nil
}

func (s *ProfileService) UnfollowUser(followerID uint, username string) (*dto.ProfileResponse, error) {
	target, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	_ = s.followRepo.UnfollowUser(followerID, target.ID)

	return &dto.ProfileResponse{
		Profile: dto.ProfilePayload{
			Username:  target.Username,
			Bio:       target.Bio,
			Image:     target.Image,
			Following: false,
		},
	}, nil
}
