package domain

type RefreshTokenResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenTTLSeconds int    `json:"access_token_ttl_seconds"`
	RefreshToken          string `json:"refresh_token"`
}
