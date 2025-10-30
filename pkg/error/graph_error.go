package error

import (
	"fmt"
	"reflect"
	"strings"

	appLogger "github.com/mohamadrezamomeni/graph/pkg/log"
)

const empty = "empty"

type ErrorType = int

const (
	UnExpected ErrorType = iota + 1
	Forbidden
	BadRequest
	NotFound
	Duplicate
)

type AppError struct {
	args      []any
	pattern   string
	scope     string
	err       error
	isPrinted bool
	input     []any
	errorType ErrorType
}

func Scope(scope string) *AppError {
	return &AppError{
		isPrinted: true,
		args:      []any{},
		pattern:   "",
		err:       nil,
		scope:     fmt.Sprintf("\"%s\"", scope),
		input:     []any{},
	}
}

func Wrap(err error) *AppError {
	return &AppError{
		isPrinted: true,
		args:      []any{},
		pattern:   "",
		err:       err,
		scope:     fmt.Sprintf("\"%s\"", empty),
	}
}

func (m *AppError) GetErrorType() ErrorType {
	if m.errorType != 0 {
		return m.errorType
	}

	m, ok := m.err.(*AppError)

	if ok {
		return m.GetErrorType()
	}

	return UnExpected
}

func (m *AppError) Message() string {
	message := m.matchPatternAndArgs()
	if len(message) > 0 {
		return message
	}

	m, ok := m.err.(*AppError)

	if ok {
		return m.Message()
	}

	return ""
}

func (m *AppError) UnExpected() *AppError {
	m.errorType = UnExpected
	return m
}

func (m *AppError) NotFound() *AppError {
	m.errorType = NotFound
	return m
}

func (m *AppError) BadRequest() *AppError {
	m.errorType = BadRequest
	return m
}

func (m *AppError) Forbidden() *AppError {
	m.errorType = Forbidden
	return m
}

func (m *AppError) Duplicate() *AppError {
	m.errorType = Forbidden
	return m
}

func (m *AppError) DeactiveWrite() *AppError {
	m.isPrinted = false
	return m
}

func (m *AppError) ActiveWrite() *AppError {
	m.isPrinted = true
	return m
}

func (m *AppError) Scope(scope string) *AppError {
	m.scope = fmt.Sprintf("\"%s\"", scope)
	return m
}

func (m *AppError) Error() string {
	message := fmt.Sprintf("the scope is %s and the main error is \"%s\"", m.scope, m.mainError())

	messageInput := m.getInputMessage()

	if len(messageInput) > 0 {
		message += fmt.Sprintf(` also we got ("%s")`, messageInput)
	}

	additionlMessage := m.matchPatternAndArgs()

	if len(additionlMessage) > 0 {
		message += " the additional information is " + `"` + additionlMessage + `"`
	}
	return message
}

func (m *AppError) matchPatternAndArgs() string {
	additionlMessage := ""
	if len(m.pattern) > 0 && len(m.args) > 0 {
		additionlMessage = fmt.Sprintf(m.pattern, m.args...)
	} else if len(m.pattern) > 0 {
		additionlMessage = m.pattern
	}
	return additionlMessage
}

func (m *AppError) Input(data ...any) *AppError {
	m.input = data
	return m
}

func (m *AppError) mainError() string {
	err, ok := m.err.(*AppError)

	if ok {
		return err.mainError()
	}

	if m.err != nil {
		return m.err.Error()
	}
	return "nothing"
}

func (m *AppError) getInputMessage() string {
	messages := []string{}
	for _, item := range m.input {
		messages = append(messages, m.translateInput(item))
	}
	return strings.Join(messages, `", "`)
}

func (m *AppError) translateInput(inpt any) string {
	val := reflect.ValueOf(inpt)

	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", val.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", val.Bool())
	case reflect.Struct:
		return fmt.Sprintf("%#v", val.Interface())
	default:
		return fmt.Sprintf("%v", val.Interface())
	}
}

func (m *AppError) Errorf(pattern string, args ...any) error {
	m.args = args
	m.pattern = pattern
	if m.isPrinted {
		appLogger.Warning(m.Error())
	}
	return m
}

func (m *AppError) DebuggingErrorf(pattern string, args ...any) error {
	m.args = args
	m.pattern = pattern
	if m.isPrinted {
		appLogger.Debug(m.Error())
	}
	return m
}

func (m *AppError) DebuggingError() *AppError {
	if m.isPrinted {
		appLogger.Debug(m.Error())
	}
	return m
}

func (m *AppError) ErrorWrite() error {
	if m.isPrinted {
		appLogger.Warning(m.Error())
	}
	return m
}

func GetErrorError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	m, ok := err.(*AppError)
	if ok {
		return m, true
	}
	return nil, false
}
