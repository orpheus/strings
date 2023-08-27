package controller

type ErrorDescriptor struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type ApiError struct {
	ErrorDescriptor `json:"errorDescriptor"`
}

func NewApiError(code int, msg string) ApiError {
	return ApiError{ErrorDescriptor{
		ErrorCode:    code,
		ErrorMessage: msg,
	}}
}
