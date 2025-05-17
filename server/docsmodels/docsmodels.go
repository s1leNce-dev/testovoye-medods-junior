package docsmodels

type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type UUIDResponse struct {
	Message string `json:"message" example:"success"`
	Data    string `json:"data" example:"f8dc39e0-9b2b-4c58-bfa0-dc6dc4d6479b"`
}

type SuccessMessage struct {
	Message string `json:"message" example:"success"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid user_id format"`
}
