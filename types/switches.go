package types

type SwitchesPOST struct {
	Run_after  uint     `json:"run_after"`
	Recipients []string `json:"recipients"`
	Content    string   `json:"content"`
}

type Switch struct {
	Id         int64    `json:"id"`
	Run_after  uint     `json:"run_after"`
	Recipients []string `json:"recipients"`
	Content    string   `json:"content"`
}

type SwitchesResponse struct {
	Switches []*Switch
}
