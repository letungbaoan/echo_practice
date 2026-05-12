package dto

type CreateCommentRequest struct {
	Comment struct {
		Body string `json:"body" validate:"required,min=1"`
	} `json:"comment" validate:"required"`
}

type CommentResponse struct {
	Comment CommentPayload `json:"comment"`
}

type CommentsResponse struct {
	Comments []CommentPayload `json:"comments"`
}

type CommentPayload struct {
	ID        uint           `json:"id"`
	Body      string         `json:"body"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	Author    ProfilePayload `json:"author"`
}
