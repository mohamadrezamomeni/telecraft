package telecrafterror

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	telecraftLogger "github.com/mohamadrezamomeni/telecraft/pkg/log"
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

type TeleCraftError struct {
	args      []any
	pattern   string
	scope     string
	err       error
	isPrinted bool
	input     []any
	errorType ErrorType
}

func Scope(scope string) *TeleCraftError {
	return &TeleCraftError{
		isPrinted: true,
		args:      []any{},
		pattern:   "",
		err:       nil,
		scope:     fmt.Sprintf("\"%s\"", scope),
		input:     []any{},
	}
}

func Wrap(err error) *TeleCraftError {
	return &TeleCraftError{
		isPrinted: true,
		args:      []any{},
		pattern:   "",
		err:       err,
		scope:     fmt.Sprintf("\"%s\"", empty),
	}
}

func (m *TeleCraftError) GetErrorType() ErrorType {
	errorType := m.errorType

	if errorType != 0 {
		return m.errorType
	}

	m, ok := m.err.(*TeleCraftError)

	if ok {
		return m.GetErrorType()
	}

	return UnExpected
}

func (m *TeleCraftError) Message() string {
	message := m.matchPatternAndArgs()
	if len(message) > 0 {
		return message
	}

	m, ok := m.err.(*TeleCraftError)

	if ok {
		return m.Message()
	}

	return ""
}

func (m *TeleCraftError) UnExpected() *TeleCraftError {
	m.errorType = UnExpected
	return m
}

func (m *TeleCraftError) NotFound() *TeleCraftError {
	m.errorType = NotFound
	return m
}

func (m *TeleCraftError) BadRequest() *TeleCraftError {
	m.errorType = BadRequest
	return m
}

func (m *TeleCraftError) Forbidden() *TeleCraftError {
	m.errorType = Forbidden
	return m
}

func (m *TeleCraftError) Duplicate() *TeleCraftError {
	m.errorType = Forbidden
	return m
}

func (m *TeleCraftError) DeactiveWrite() *TeleCraftError {
	m.isPrinted = false
	return m
}

func (m *TeleCraftError) ActiveWrite() *TeleCraftError {
	m.isPrinted = true
	return m
}

func (m *TeleCraftError) Scope(scope string) *TeleCraftError {
	m.scope = fmt.Sprintf("\"%s\"", scope)
	return m
}

func (m *TeleCraftError) Error() string {
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

func (m *TeleCraftError) matchPatternAndArgs() string {
	additionlMessage := ""
	if len(m.pattern) > 0 && len(m.args) > 0 {
		additionlMessage = fmt.Sprintf(m.pattern, m.args...)
	} else if len(m.pattern) > 0 {
		additionlMessage = m.pattern
	}
	return additionlMessage
}

func (m *TeleCraftError) Input(data ...any) *TeleCraftError {
	m.input = data
	return m
}

func (m *TeleCraftError) mainError() string {
	err, ok := m.err.(*TeleCraftError)

	if ok {
		return err.mainError()
	}

	if m.err != nil {
		return m.err.Error()
	}
	return "nothing"
}

func (m *TeleCraftError) getInputMessage() string {
	messages := []string{}
	for _, item := range m.input {
		messages = append(messages, m.translateInput(item))
	}
	return strings.Join(messages, `", "`)
}

func (m *TeleCraftError) translateInput(inpt any) string {
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

func (m *TeleCraftError) Fatalf(pattern string, args ...any) {
	m.args = args
	m.pattern = pattern
	if m.isPrinted {
		telecraftLogger.Debugging(m.Error())
	}
	os.Exit(1)
}

func (m *TeleCraftError) Errorf(pattern string, args ...any) error {
	m.args = args
	m.pattern = pattern
	if m.isPrinted {
		telecraftLogger.Warrning(m.Error())
	}
	return m
}

func (m *TeleCraftError) DebuggingErrorf(pattern string, args ...any) error {
	m.args = args
	m.pattern = pattern
	if m.isPrinted {
		telecraftLogger.Debugging(m.Error())
	}
	return m
}

func (m *TeleCraftError) DebuggingError() *TeleCraftError {
	if m.isPrinted {
		telecraftLogger.Debugging(m.Error())
	}
	return m
}

func (m *TeleCraftError) Fatal() {
	if m.isPrinted {
		telecraftLogger.Debugging(m.Error())
	}
	os.Exit(1)
}

func (m *TeleCraftError) ErrorWrite() error {
	if m.isPrinted {
		telecraftLogger.Warrning(m.Error())
	}
	return m
}

func GetMomoError(err error) (*TeleCraftError, bool) {
	if err == nil {
		return nil, false
	}
	m, ok := err.(*TeleCraftError)
	if ok {
		return m, true
	}
	return nil, false
}
