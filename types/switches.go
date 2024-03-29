package types

type SwitchesPOST struct {
	Run_after  uint     `json:"run_after"`
	Recipients []string `json:"recipients"`
	Content    string   `json:"content"`
	Subject    string   `json:"subject"`
}

type Switch struct {
	AccountId  int64    `json:"account_id,omitempty"`
	Id         int64    `json:"id"`
	Subject    string   `json:"subject"`
	Content    string   `json:"content"`
	Run_after  uint     `json:"run_after"`
	Recipients []string `json:"recipients"`
}

type SwitchesResponse struct {
	Switches []Switch
}
