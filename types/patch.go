package types

type SwitchesPATCH struct {
	Id         int64    `json:"id,omitempty"`
	Run_after  uint     `json:"run_after,omitempty"`
	Recipients []string `json:"recipients,omitempty"`
	Content    string   `json:"content,omitempty"`
}
