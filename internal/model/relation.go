package model

type relation struct {
	Name    string `json:"name"`
	Model   string `json:"model"`
	Foreign string `json:"foreign"`
	Key     string `json:"key"`
}
