package model

type LiveCheckResponse struct {
	Status string `json:"status"`
}

type ReadyCheckResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}
