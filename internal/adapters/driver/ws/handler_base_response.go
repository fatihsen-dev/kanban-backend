package ws

import (
	"encoding/json"
)

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(message string, data interface{}) []byte {
	response := BaseResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	jsonBytes, _ := json.Marshal(response)
	return jsonBytes
}

func NewErrorResponse(err string) []byte {
	response := BaseResponse{
		Status:  false,
		Message: err,
	}
	jsonBytes, _ := json.Marshal(response)
	return jsonBytes
}

func NewAbortResponse(message string) []byte {
	response := BaseResponse{
		Status:  false,
		Message: message,
	}
	jsonBytes, _ := json.Marshal(response)
	return jsonBytes
}
