package wraper

import "encoding/json"

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type Response[T any] struct {
	BaseResponse
	Data T `json:"data,omitempty"`
}

func (r Response[T]) ToJSON() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func NewSuccessResponse[T any](data T) Response[T] {
	return Response[T]{
		BaseResponse: BaseResponse{
			Code:    200,
			Message: "Success",
		},
		Data: data,
	}
}

func NewBadResponse[T any](message string, data T) Response[T] {
	return Response[T]{
		BaseResponse: BaseResponse{
			Code:    500,
			Message: message,
		},
		Data: data,
	}
}
