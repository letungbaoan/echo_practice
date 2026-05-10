package services

import (
	"echo_practice/internal/apperrors"
	"echo_practice/internal/dto"
	"echo_practice/internal/models"
	"echo_practice/internal/repositories"
	"echo_practice/internal/utils"
	"time"
)

type ArticleService struct {
	articleRepo *repositories.ArticleRepository
	tagRepo     *repositories.TagRepository
	userRepo    *repositories.UserRepository
}

func NewArticleService(
	articleRepo *repositories.ArticleRepository,
	tagRepo *repositories.TagRepository,
	userRepo *repositories.UserRepository,
) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		tagRepo:     tagRepo,
		userRepo:    userRepo,
	}
}

func (s *ArticleService) CreateArticle(authorID uint, req dto.CreateArticleRequest) (*dto.ArticleResponse, error) {
	slug := utils.GenerateSlug(req.Article.Title)

	article := &models.Article{
		Slug:        slug,
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		AuthorID:    authorID,
	}

	if len(req.Article.TagList) > 0 {
		tags := make([]models.Tag, 0)
		for _, tagName := range req.Article.TagList {
			tag, err := s.tagRepo.FindOrCreate(tagName)
			if err != nil {
				return nil, err
			}
			tags = append(tags, *tag)
		}
		article.Tags = tags
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	author, err := s.userRepo.FindByID(authorID)
	if err != nil {
		return nil, err
	}
	article.Author = *author

	return toArticleResponse(article, false, 0), nil
}

func (s *ArticleService) GetArticle(slug string, currentUserID *uint) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return toArticleResponse(article, false, 0), nil
}

func (s *ArticleService) UpdateArticle(slug string, authorID uint, req dto.UpdateArticleRequest) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	if article.AuthorID != authorID {
		return nil, apperrors.ErrForbidden
	}

	payload := req.Article
	if payload.Title != "" {
		article.Title = payload.Title
		article.Slug = utils.GenerateSlug(payload.Title)
	}
	if payload.Description != "" {
		article.Description = payload.Description
	}
	if payload.Body != "" {
		article.Body = payload.Body
	}

	if len(payload.TagList) > 0 {
		tags := make([]models.Tag, 0)
		for _, tagName := range payload.TagList {
			tag, err := s.tagRepo.FindOrCreate(tagName)
			if err != nil {
				return nil, err
			}
			tags = append(tags, *tag)
		}
		article.Tags = tags
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	return toArticleResponse(article, false, 0), nil
}

func (s *ArticleService) DeleteArticle(slug string, authorID uint) error {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return apperrors.ErrNotFound
		}
		return err
	}

	if article.AuthorID != authorID {
		return apperrors.ErrForbidden
	}

	return s.articleRepo.Delete(slug)
}

func toArticleResponse(article *models.Article, favorited bool, favoritesCount int) *dto.ArticleResponse {
	tagList := make([]string, 0)
	if article.Tags != nil {
		for _, tag := range article.Tags {
			tagList = append(tagList, tag.Name)
		}
	}

	return &dto.ArticleResponse{
		Article: dto.ArticlePayload{
			Slug:           article.Slug,
			Title:          article.Title,
			Description:    article.Description,
			Body:           article.Body,
			TagList:        tagList,
			CreatedAt:      article.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      article.UpdatedAt.Format(time.RFC3339),
			Favorited:      favorited,
			FavoritesCount: favoritesCount,
			Author: dto.ProfilePayload{
				Username:  article.Author.Username,
				Bio:       article.Author.Bio,
				Image:     article.Author.Image,
				Following: false,
			},
		},
	}
}
