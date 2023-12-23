package model

type ErrorDto struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}
