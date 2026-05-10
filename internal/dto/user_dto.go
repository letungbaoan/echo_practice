package dto

type RegisterRequest struct {
	User struct {
		Username string `json:"username" validate:"required,min=3,max=32"`
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=6,max=72"`
	} `json:"user" validate:"required"`
}

type LoginRequest struct {
	User struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"user" validate:"required"`
}

type UserResponse struct {
	User UserPayload `json:"user"`
}

type UserPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type UpdateRequest struct {
	User struct {
		Email    string `json:"email"    validate:"omitempty,email"`
		Username string `json:"username" validate:"omitempty,min=3,max=32"`
		Password string `json:"password" validate:"omitempty,min=6,max=72"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"user" validate:"required"`
}
