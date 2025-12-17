package provider

import (
	"net/http"
	"time"
)

const (
	StatusOK      = "ok"
	StatusSuccess = "success"
	CodeSuccess   = 200
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 2 * time.Second,
	}
}
