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
	articleRepo  *repositories.ArticleRepository
	tagRepo      *repositories.TagRepository
	userRepo     *repositories.UserRepository
	followRepo   *repositories.FollowRepository
	favoriteRepo *repositories.FavoriteRepository
}

func NewArticleService(
	articleRepo *repositories.ArticleRepository,
	tagRepo *repositories.TagRepository,
	userRepo *repositories.UserRepository,
	followRepo *repositories.FollowRepository,
	favoriteRepo *repositories.FavoriteRepository,
) *ArticleService {
	return &ArticleService{
		articleRepo:  articleRepo,
		tagRepo:      tagRepo,
		userRepo:     userRepo,
		followRepo:   followRepo,
		favoriteRepo: favoriteRepo,
	}
}

type ListArticlesFilter struct {
	Tag               string
	AuthorUsername    string
	FavoritedUsername string
	Limit             int
	Offset            int
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

	return s.buildArticleResponse(article, &authorID)
}

func (s *ArticleService) GetArticle(slug string, currentUserID *uint) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return s.buildArticleResponse(article, currentUserID)
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

	return s.buildArticleResponse(article, &authorID)
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

func (s *ArticleService) FavoriteArticle(slug string, currentUserID uint) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	if err := s.favoriteRepo.Favorite(currentUserID, article.ID); err != nil {
		return nil, err
	}
	return s.buildArticleResponse(article, &currentUserID)
}

func (s *ArticleService) UnfavoriteArticle(slug string, currentUserID uint) (*dto.ArticleResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	if err := s.favoriteRepo.Unfavorite(currentUserID, article.ID); err != nil {
		return nil, err
	}
	return s.buildArticleResponse(article, &currentUserID)
}

func (s *ArticleService) buildArticleResponse(article *models.Article, currentUserID *uint) (*dto.ArticleResponse, error) {
	favoritesCount, err := s.favoriteRepo.CountByArticle(article.ID)
	if err != nil {
		return nil, err
	}

	var favorited, following bool
	if currentUserID != nil {
		favorited, err = s.favoriteRepo.IsFavorited(*currentUserID, article.ID)
		if err != nil {
			return nil, err
		}
		following, err = s.followRepo.IsFollowing(*currentUserID, article.AuthorID)
		if err != nil {
			return nil, err
		}
	}

	return &dto.ArticleResponse{Article: toArticlePayload(article, favorited, favoritesCount, following)}, nil
}

func (s *ArticleService) ListArticles(filter ListArticlesFilter, currentUserID *uint) (*dto.ArticlesResponse, error) {
	limit, offset := normalizePagination(filter.Limit, filter.Offset)
	articles, total, err := s.articleRepo.List(repositories.ListFilter{
		Tag:               filter.Tag,
		AuthorUsername:    filter.AuthorUsername,
		FavoritedUsername: filter.FavoritedUsername,
		Limit:             limit,
		Offset:            offset,
	})
	if err != nil {
		return nil, err
	}
	return s.buildArticlesResponse(articles, total, currentUserID)
}

func (s *ArticleService) FeedArticles(currentUserID uint, limit, offset int) (*dto.ArticlesResponse, error) {
	limit, offset = normalizePagination(limit, offset)
	articles, total, err := s.articleRepo.Feed(currentUserID, limit, offset)
	if err != nil {
		return nil, err
	}
	return s.buildArticlesResponse(articles, total, &currentUserID)
}

func (s *ArticleService) buildArticlesResponse(articles []models.Article, total int64, currentUserID *uint) (*dto.ArticlesResponse, error) {
	authorIDs := make([]uint, 0, len(articles))
	articleIDs := make([]uint, 0, len(articles))
	for _, a := range articles {
		authorIDs = append(authorIDs, a.AuthorID)
		articleIDs = append(articleIDs, a.ID)
	}

	favCounts, err := s.favoriteRepo.CountByArticles(articleIDs)
	if err != nil {
		return nil, err
	}

	var following, favorited map[uint]bool
	if currentUserID != nil {
		following, err = s.followRepo.IsFollowingMany(*currentUserID, authorIDs)
		if err != nil {
			return nil, err
		}
		favorited, err = s.favoriteRepo.IsFavoritedMany(*currentUserID, articleIDs)
		if err != nil {
			return nil, err
		}
	}

	payloads := make([]dto.ArticlePayload, 0, len(articles))
	for i := range articles {
		a := &articles[i]
		payloads = append(payloads, toArticlePayload(a, favorited[a.ID], favCounts[a.ID], following[a.AuthorID]))
	}
	return &dto.ArticlesResponse{
		Articles:      payloads,
		ArticlesCount: int(total),
	}, nil
}

func normalizePagination(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func toArticlePayload(article *models.Article, favorited bool, favoritesCount int, following bool) dto.ArticlePayload {
	tagList := make([]string, 0)
	if article.Tags != nil {
		for _, tag := range article.Tags {
			tagList = append(tagList, tag.Name)
		}
	}

	return dto.ArticlePayload{
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
			Following: following,
		},
	}
}
