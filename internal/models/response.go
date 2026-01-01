package models

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ShareList struct {
	Missing    *ShareSection `json:"missing,omitempty"`
	Duplicates *ShareSection `json:"duplicates,omitempty"`
	Text       string        `json:"text"`
}

type ShareSection struct {
	Count   int                       `json:"count"`
	Sections map[string][]ShareItem   `json:"sections"`
}

type ShareItem struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity,omitempty"`
}
