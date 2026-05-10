package dto

type CreateArticleRequest struct {
	Article struct {
		Title       string   `json:"title"       validate:"required,min=1"`
		Description string   `json:"description" validate:"required,min=1"`
		Body        string   `json:"body"        validate:"required,min=1"`
		TagList     []string `json:"tagList"     validate:"omitempty"`
	} `json:"article" validate:"required"`
}

type UpdateArticleRequest struct {
	Article struct {
		Title       string   `json:"title"       validate:"omitempty,min=1"`
		Description string   `json:"description" validate:"omitempty,min=1"`
		Body        string   `json:"body"        validate:"omitempty,min=1"`
		TagList     []string `json:"tagList"     validate:"omitempty"`
	} `json:"article" validate:"required"`
}

type ArticleResponse struct {
	Article ArticlePayload `json:"article"`
}

type ArticlesResponse struct {
	Articles      []ArticlePayload `json:"articles"`
	ArticlesCount int              `json:"articlesCount"`
}

type ArticlePayload struct {
	Slug             string         `json:"slug"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	Body             string         `json:"body"`
	TagList          []string       `json:"tagList"`
	CreatedAt        string         `json:"createdAt"`
	UpdatedAt        string         `json:"updatedAt"`
	Favorited        bool           `json:"favorited"`
	FavoritesCount   int            `json:"favoritesCount"`
	Author           ProfilePayload `json:"author"`
}
