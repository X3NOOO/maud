package types

type SwitchesPOST struct {
	Run_after  uint     `json:"run_after"`
	Recipients []string `json:"recipients"`
	Content    string   `json:"content"`
}

type Switch struct {
	AccountId  int64	`json:"account_id,omitempty"`
	Id         int64    `json:"id"`
	Content    string   `json:"content"`
	Run_after  uint     `json:"run_after"`
	Recipients []string `json:"recipients"`
}

type SwitchesResponse struct {
	Switches []Switch
}
