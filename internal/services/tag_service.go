package services

import (
	"echo_practice/internal/dto"
	"echo_practice/internal/repositories"
)

type TagService struct {
	tagRepo *repositories.TagRepository
}

func NewTagService(tagRepo *repositories.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

func (s *TagService) ListTags() (*dto.TagsResponse, error) {
	names, err := s.tagRepo.ListUsed()
	if err != nil {
		return nil, err
	}
	if names == nil {
		names = []string{}
	}
	return &dto.TagsResponse{Tags: names}, nil
}
