package types

type LoginPOST struct {
	Nick     string `json:"nick"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthorizationToken string `json:"authorization_token"`
}
