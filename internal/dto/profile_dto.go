package dto

type ProfileResponse struct {
	Profile ProfilePayload `json:"profile"`
}

type ProfilePayload struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}
