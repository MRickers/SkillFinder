package model

import "encoding/json"

type ErrorDto struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
}

func Serialize(dto ErrorDto) ([]byte, error) {
	return json.Marshal(dto)
}
