package types

import ()

type Logger interface {
	Info(message string)
	Error(message string)
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}