package dto

// An AuthenticationRequest represents an authentication request
type AuthenticationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// An AuthenticationResponse contains a valid auth token to use in subsequent requests
type AuthenticationResponse struct {
	JWT string `json:"bearer"`
}
