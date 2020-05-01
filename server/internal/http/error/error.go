package error

import (
	"time"
)

type ErrorResponse struct {
	Code   int       `json:"code"`
	Status int       `json:"status"`
	Error  string    `json:"error"`
	Time   time.Time `json:"time"`
}

type HTTPError struct {
	Code   int
	Status int
	Err    error
}

func (e HTTPError) Error() string {
	return e.Err.Error()
}

func (e HTTPError) toResponse() ErrorResponse {
	return ErrorResponse{
		Code:   e.Code,
		Status: e.Status,
		Error:  e.Err.Error(),
		Time:   time.Now(),
	}
}
