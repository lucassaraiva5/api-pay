package inbound

type Method struct {
	Type string `json:"type"`
	Card Card   `json:"card"`
}
