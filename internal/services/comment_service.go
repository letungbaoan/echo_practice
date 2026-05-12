package services

import (
	"echo_practice/internal/apperrors"
	"echo_practice/internal/dto"
	"echo_practice/internal/models"
	"echo_practice/internal/repositories"
	"time"
)

type CommentService struct {
	commentRepo *repositories.CommentRepository
	articleRepo *repositories.ArticleRepository
	userRepo    *repositories.UserRepository
	followRepo  *repositories.FollowRepository
}

func NewCommentService(
	commentRepo *repositories.CommentRepository,
	articleRepo *repositories.ArticleRepository,
	userRepo *repositories.UserRepository,
	followRepo *repositories.FollowRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
		userRepo:    userRepo,
		followRepo:  followRepo,
	}
}

func (s *CommentService) AddComment(slug string, authorID uint, req dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	comment := &models.Comment{
		Body:      req.Comment.Body,
		ArticleID: article.ID,
		AuthorID:  authorID,
	}
	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	author, err := s.userRepo.FindByID(authorID)
	if err != nil {
		return nil, err
	}
	comment.Author = *author

	return &dto.CommentResponse{Comment: toCommentPayload(comment, false)}, nil
}

func (s *CommentService) ListComments(slug string, currentUserID *uint) (*dto.CommentsResponse, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	comments, err := s.commentRepo.ListByArticle(article.ID)
	if err != nil {
		return nil, err
	}

	authorIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		authorIDs = append(authorIDs, c.AuthorID)
	}

	var following map[uint]bool
	if currentUserID != nil {
		following, err = s.followRepo.IsFollowingMany(*currentUserID, authorIDs)
		if err != nil {
			return nil, err
		}
	}

	payloads := make([]dto.CommentPayload, 0, len(comments))
	for i := range comments {
		c := &comments[i]
		payloads = append(payloads, toCommentPayload(c, following[c.AuthorID]))
	}
	return &dto.CommentsResponse{Comments: payloads}, nil
}

func (s *CommentService) DeleteComment(slug string, commentID, currentUserID uint) error {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		if repositories.IsNotFound(err) {
			return apperrors.ErrNotFound
		}
		return err
	}

	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		if repositories.IsNotFound(err) {
			return apperrors.ErrNotFound
		}
		return err
	}

	if comment.ArticleID != article.ID {
		return apperrors.ErrNotFound
	}
	if comment.AuthorID != currentUserID {
		return apperrors.ErrForbidden
	}

	return s.commentRepo.Delete(commentID)
}

func toCommentPayload(comment *models.Comment, following bool) dto.CommentPayload {
	return dto.CommentPayload{
		ID:        comment.ID,
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
		Author: dto.ProfilePayload{
			Username:  comment.Author.Username,
			Bio:       comment.Author.Bio,
			Image:     comment.Author.Image,
			Following: following,
		},
	}
}
