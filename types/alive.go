package types

type AlivePOST struct {
	Nick string `json:"nick"`
}

type AliveResponse struct {
	Date string `json:"date"`
}