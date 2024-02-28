package types

type RegisterPOST struct {
	Nick     string `json:"nick"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Nick               string `json:"nick"`
	AuthorizationToken string `json:"authorization_token"`
}
