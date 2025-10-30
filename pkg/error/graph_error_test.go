package error

import (
	"fmt"
	"os"
	"testing"

	appLogger "github.com/mohamadrezamomeni/graph/pkg/log"
)

func TestMain(m *testing.M) {
	appLogger.Discard()
	code := m.Run()

	os.Exit(code)
}

func TestWithoutMainError(t *testing.T) {
	scopeTest := "test.TestWithoutMainError"

	for _, testCase := range []struct {
		message string
		err     error
	}{
		{
			err:     Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was with arguments %d", 1),
			message: "the scope is \"test.TestWithoutMainError\" and the main error is \"nothing\" the additional information is \"patern was with arguments 1\"",
		},
		{
			err:     Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was without any arguments"),
			message: "the scope is \"test.TestWithoutMainError\" and the main error is \"nothing\" the additional information is \"patern was without any arguments\"",
		},
		{
			err:     Wrap(fmt.Errorf("database error")).DeactiveWrite().DebuggingErrorf("patern was without any arguments"),
			message: "the scope is \"empty\" and the main error is \"database error\" the additional information is \"patern was without any arguments\"",
		},
		{
			err:     Wrap(fmt.Errorf("database error")).Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was without any arguments"),
			message: "the scope is \"test.TestWithoutMainError\" and the main error is \"database error\" the additional information is \"patern was without any arguments\"",
		},
		{
			err:     Wrap(fmt.Errorf("database error")).Scope(scopeTest).Input(struct{ Domain string }{Domain: "google.com"}, "ssss", map[string]string{"name": "mic"}).ErrorWrite(),
			message: `the scope is "test.TestWithoutMainError" and the main error is "database error" also we got ("struct { Domain string }{Domain:"google.com"}", "ssss", "map[name:mic]")`,
		},
	} {
		if testCase.err.Error() != testCase.message {
			t.Error("error to compare we expected error and the error we were given")
		}
	}
}

func TestErrorType(t *testing.T) {
	scope := "test.TestErrorType"

	for i, testCase := range []struct {
		err       *AppError
		errorType int
	}{
		{
			err:       Scope(scope).Forbidden(),
			errorType: Forbidden,
		},
		{
			err:       Scope(scope).NotFound(),
			errorType: NotFound,
		},
		{
			err:       Scope(scope),
			errorType: UnExpected,
		},
		{
			err:       Scope(scope).BadRequest(),
			errorType: BadRequest,
		},
		{
			err:       Scope(scope).UnExpected(),
			errorType: UnExpected,
		},
	} {
		if testCase.errorType != testCase.err.GetErrorType() {
			t.Errorf("some thing went wrong to compare errorType at index %d we expected %d but we got %d",
				i,
				testCase.errorType,
				testCase.err.GetErrorType(),
			)
		}
	}
}

func TestMessage(t *testing.T) {
	scope := "test.TestMessage"

	e := Scope(scope).Errorf("hello world")
	for _, testCase := range []struct {
		err     error
		message string
	}{
		{
			err:     e,
			message: "hello world",
		},
		{
			err:     Wrap(e).Errorf("hello another world"),
			message: "hello another world",
		},
		{
			err:     Wrap(e),
			message: "hello world",
		},
	} {
		appErr, _ := testCase.err.(*AppError)

		if appErr.Message() != testCase.message {
			t.Errorf("message must be %s but we got %s", testCase.message, appErr.Message())
		}
	}
}
