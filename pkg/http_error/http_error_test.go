package httperror

import (
	"fmt"
	"net/http"
	"testing"

	appError "github.com/mohamadrezamomeni/graph/pkg/error"
)

func TestStatus(t *testing.T) {
	for _, testCase := range []struct {
		receivedStatus int
		expected       int
	}{
		{
			receivedStatus: getStatus(appError.Scope("test")),
			expected:       http.StatusInternalServerError,
		},
		{
			receivedStatus: getStatus(appError.Scope("test").Forbidden()),
			expected:       http.StatusForbidden,
		},
		{
			receivedStatus: getStatus(appError.Scope("test").BadRequest()),
			expected:       http.StatusBadRequest,
		},
		{
			receivedStatus: getStatus(fmt.Errorf("error")),
			expected:       http.StatusInternalServerError,
		},
	} {
		if testCase.receivedStatus != testCase.expected {
			t.Errorf("the error status we got is %d but we got %d", testCase.receivedStatus, testCase.expected)
		}
	}
}

func TestMessage(t *testing.T) {
	for _, testCase := range []struct {
		receivedMessage string
		expected        string
	}{
		{
			receivedMessage: getMessage(fmt.Errorf("")),
			expected:        "something went wrong",
		},
		{
			receivedMessage: getMessage(appError.Scope("test").Forbidden()),
			expected:        "something went wrong",
		},
		{
			receivedMessage: getMessage(appError.Scope("")),
			expected:        "something went wrong",
		},
		{
			receivedMessage: getMessage(appError.Scope("").Errorf("hello world")),
			expected:        "hello world",
		},
	} {
		if testCase.receivedMessage != testCase.expected {
			t.Errorf("the error status we got is '%s' but we got '%s'", testCase.receivedMessage, testCase.expected)
		}
	}
}
