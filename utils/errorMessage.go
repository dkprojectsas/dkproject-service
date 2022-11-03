package utils

import (
	"fmt"
	"runtime"
)

const (
	Error4xx               = "error 4xx"
	ErrorUnauthorizeUser   = "error 401 unauthorize user"
	ErrorInternalServer    = "error 500 internal server error"
	ErrorBadRequest        = "error 400 bad request"
	ErrorParameterNotFound = "error 404 parameter not found"
)

type errMessageModel struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Errors  string `json:"errors"`
}

func ErrorMessages(errMessage string, err error) *errMessageModel {
	return &errMessageModel{
		Status:  "error",
		Message: errMessage,
		Errors:  err.Error(),
	}
}

func DebugError(err error, msg string) {
	if err != nil {
		_, filename, line, _ := runtime.Caller(1)
		fmt.Printf("ERROR: '%s', %s:%d ", msg, filename, line)
	}
}
