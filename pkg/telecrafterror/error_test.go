package telecrafterror

import (
	"fmt"
	"testing"
)

func TestWithoutMainError(t *testing.T) {
	scopeTest := "test.TestWithoutMainError"
	err := Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was with arguments %d", 1)

	message := "the scope is \"test.TestWithoutMainError\" and the main error is \"nothing\" the additional information is \"patern was with arguments 1\""

	if err.Error() != message {
		t.Error("error to compare we expected error and the error we were given")
	}

	err = Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was without any arguments")

	message = "the scope is \"test.TestWithoutMainError\" and the main error is \"nothing\" the additional information is \"patern was without any arguments\""
	if err.Error() != message {
		t.Error("error to compare we expected error and the error we were given")
	}

	err = Wrap(fmt.Errorf("database error")).DeactiveWrite().DebuggingErrorf("patern was without any arguments")
	message = "the scope is \"empty\" and the main error is \"database error\" the additional information is \"patern was without any arguments\""
	if err.Error() != message {
		t.Error("error to compare we expected error and the error we were given")
	}

	err = Wrap(fmt.Errorf("database error")).Scope(scopeTest).DeactiveWrite().DebuggingErrorf("patern was without any arguments")
	message = "the scope is \"test.TestWithoutMainError\" and the main error is \"database error\" the additional information is \"patern was without any arguments\""
	if err.Error() != message {
		t.Error("error to compare we expected error and the error we were given")
	}

	err = Wrap(fmt.Errorf("database error")).Scope(scopeTest).Input(struct{ Domain string }{Domain: "google.com"}, "ssss", map[string]string{"name": "mic"}).ErrorWrite()
	message = `the scope is "test.TestWithoutMainError" and the main error is "database error" also we got ("struct { Domain string }{Domain:"google.com"}", "ssss", "map[name:mic]")`
	if message != err.Error() {
		t.Error("the input message isn't generated well")
	}
}

func TestErrorType(t *testing.T) {
	scope := "test.TestErrorType"
	e := Scope(scope).Forbidden()

	if e.GetErrorType() != Forbidden {
		t.Error("error type must be forbidden")
	}

	e = Scope(scope).NotFound()

	if e.GetErrorType() != NotFound {
		t.Error("error type must be notfound")
	}

	e = Scope(scope)
	if e.GetErrorType() != UnExpected {
		t.Error("error type must be unexpected")
	}

	e = Scope(scope).BadRequest()
	e = Scope(scope).UnExpected()
	if e.GetErrorType() != UnExpected {
		t.Error("error type must be unexpected")
	}

	e = Scope(scope).BadRequest()
	e = Wrap(e)
	if e.GetErrorType() != BadRequest {
		t.Error("error type must be BadRequest")
	}
}

func TestMessage(t *testing.T) {
	scope := "test.TestMessage"
	message := "hello world"
	e := Scope(scope).Errorf(message)

	v, _ := e.(*TeleCraftError)
	if v.Message() != message {
		t.Errorf("message must be %s but we got %s", message, v.Message())
	}
	message2 := "hello another world"

	e1 := Wrap(e).Errorf(message2)
	v, _ = e1.(*TeleCraftError)
	if v.Message() != message2 {
		t.Errorf("message must be %s but we got %s", message2, v.Message())
	}

	e1 = Wrap(e)
	v, _ = e1.(*TeleCraftError)
	if v.Message() != message {
		t.Errorf("message must be %s but we got %s", message, v.Message())
	}
}
