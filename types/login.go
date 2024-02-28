package types

type LoginPOST struct {
	Nick     string `json:"nick"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Nick               string `json:"nick"`
	AuthorizationToken string `json:"authorization_token"`
}
