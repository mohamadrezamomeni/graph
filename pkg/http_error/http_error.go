package httperror

import (
	"net/http"

	appError "github.com/mohamadrezamomeni/graph/pkg/error"
)

func Error(err error) (string, int) {
	return getMessage(err), getStatus(err)
}

func getStatus(err error) int {
	appErr, ok := err.(*appError.AppError)

	if !ok {
		return http.StatusInternalServerError
	}
	errType := appErr.GetErrorType()
	return mapAppErrorTypeToHttpStatus(errType)
}

func mapAppErrorTypeToHttpStatus(errType appError.ErrorType) int {
	switch errType {
	case appError.BadRequest:
		return http.StatusBadRequest
	case appError.Forbidden:
		return http.StatusForbidden
	case appError.UnExpected:
		return http.StatusInternalServerError
	case appError.NotFound:
		return http.StatusNotFound
	case appError.Duplicate:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func getMessage(err error) string {
	appErr, ok := err.(*appError.AppError)
	if !ok {
		return "something went wrong"
	}

	message := appErr.Message()

	if len(message) > 0 {
		return message
	}

	code := getStatus(err)
	switch code {
	case http.StatusBadRequest:
		return "input is wrong"
	case http.StatusNotFound:
		return "no record found"
	case http.StatusConflict:
		return "this record exist"
	default:
		return "something went wrong"
	}
}
