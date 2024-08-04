package exception

import (
	"fmt"
	"net/http"
	"strconv"
)

type ErrorException struct {
	ErrorCode      string
	ErrorMessage   string
	HttpStatusCode int
}

func (e *ErrorException) Error() string {
	return fmt.Sprintf("Error : %s - %s", e.ErrorCode, e.ErrorMessage)
}

func NotFoundException(errorCode string, errorMessage string) *ErrorException {
	if errorCode == "" {
		errorCode = strconv.Itoa(http.StatusNotFound)
	}

	return &ErrorException{
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
		HttpStatusCode: http.StatusNotFound,
	}
}

func ValidationException(errorCode string, errorMessage string) *ErrorException {
	if errorCode == "" {
		errorCode = strconv.Itoa(http.StatusBadRequest)
	}

	return &ErrorException{
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
		HttpStatusCode: http.StatusBadRequest,
	}
}

func UnhandledException(errorCode string, errorMessage string) *ErrorException {
	if errorCode == "" {
		errorCode = strconv.Itoa(http.StatusInternalServerError)
	}

	return &ErrorException{
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
		HttpStatusCode: http.StatusInternalServerError,
	}
}

func UnauthorizedException(errorCode string, errorMessage string) *ErrorException {
	if errorCode == "" {
		errorCode = strconv.Itoa(http.StatusUnauthorized)
	}

	return &ErrorException{
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
		HttpStatusCode: http.StatusUnauthorized,
	}
}
