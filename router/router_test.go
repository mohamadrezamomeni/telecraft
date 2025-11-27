package router

import (
	"fmt"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/telecraft/adapter/state"
	"github.com/mohamadrezamomeni/telecraft/handler"
)

func isErrorConfigsMatched(expected []*tgbotapi.MessageConfig, got []*tgbotapi.MessageConfig) bool {
	if len(expected) != len(got) {
		return false
	}

	cache := make(map[string]int)
	for _, em := range expected {
		cache[em.Text] += 1
	}

	for _, gm := range got {
		cache[gm.Text] -= 1
	}

	for _, v := range cache {
		if v != 0 {
			return false
		}
	}

	return true
}

var (
	router     *Router
	cacheState *state.Cache
)

func TestMain(m *testing.M) {
	cacheState = state.New()
	router = New("root", cacheState)

	var rootHandler handler.HandlerFunc = func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			RedirectRoot: false,
			ReleaseState: true,
			Path:         "somewhere",
			MessageConfigs: []*tgbotapi.MessageConfig{
				{Text: "this is root path"},
			},
		}, nil
	}

	router.Register("root", rootHandler)
	router.Register("users/:userID", func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			RedirectRoot: false,
			ReleaseState: true,
			Path:         "",
			MessageConfigs: []*tgbotapi.MessageConfig{
				{
					Text: fmt.Sprintf("the userID is %s", u.Params["userID"]),
				},
			},
		}, nil
	})
	router.Register("users", func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			RedirectRoot: false,
			ReleaseState: true,
			Path:         "",
			MessageConfigs: []*tgbotapi.MessageConfig{
				{
					Text: "this is users path",
				},
			},
		}, nil
	})
	router.Register("users/:userID/books", func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			RedirectRoot: false,
			ReleaseState: true,
			Path:         "",
			MessageConfigs: []*tgbotapi.MessageConfig{
				{
					Text: "this is users path",
				},
			},
		}, nil
	}, func(next handler.HandlerFunc) handler.HandlerFunc {
		return func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
			if u.Params["userID"] == "1" {
				return nil, fmt.Errorf("#1")
			}
			return next(u)
		}
	})

	router.Register("users/:userID/accounts", func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			RedirectRoot: true,
			ReleaseState: true,
			Path:         "",
			MessageConfigs: []*tgbotapi.MessageConfig{
				{
					Text: fmt.Sprintf("account is created the id is %s", u.Message.Text),
				},
			},
		}, nil
	}, func(next handler.HandlerFunc) handler.HandlerFunc {
		return func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
			if _, ok := u.Data["account"]; !ok {
				return &handler.ResponseHandlerFunc{
					RedirectRoot: false,
					ReleaseState: false,
					Path:         fmt.Sprintf("users/%s/accounts", u.Params["userID"]),
					Data: map[string]string{
						"account": "",
					},
					MessageConfigs: []*tgbotapi.MessageConfig{
						{
							Text: "input your account:",
						},
					},
				}, nil
			}
			return next(u)
		}
	})

	code := m.Run()

	os.Exit(code)
}

func TestRoutingWithoutState(t *testing.T) {
	for i, testCase := range []struct {
		input        handler.Context
		expectError  bool
		errorMessage string
		res          *handler.ResponseHandlerFunc
	}{
		{
			input: handler.Context{
				Update: &tgbotapi.Update{
					Message: &tgbotapi.Message{
						Text: "/root",
					},
				},
				Params: map[string]string{
					"userID": "1",
				},
			},
			expectError:  false,
			errorMessage: "",
			res: &handler.ResponseHandlerFunc{
				MessageConfigs: []*tgbotapi.MessageConfig{
					{Text: "this is root path"},
				},
				ReleaseState: true,
				RedirectRoot: false,
			},
		},
		{
			input: handler.Context{
				Update: &tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: "/users",
					},
				},
				Params: map[string]string{
					"userID": "1",
				},
			},
			expectError:  false,
			errorMessage: "",
			res: &handler.ResponseHandlerFunc{
				MessageConfigs: []*tgbotapi.MessageConfig{
					{Text: "this is users path"},
				},
				RedirectRoot: false,
				ReleaseState: true,
				Path:         "",
			},
		},
		{
			input: handler.Context{
				Update: &tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: "/users/1",
					},
				},
				Params: map[string]string{
					"userID": "1",
				},
			},
			expectError:  false,
			errorMessage: "",
			res: &handler.ResponseHandlerFunc{
				MessageConfigs: []*tgbotapi.MessageConfig{
					{Text: "the userID is 1"},
				},
				RedirectRoot: false,
				ReleaseState: true,
				Path:         "",
			},
		},
		{
			input: handler.Context{
				Update: &tgbotapi.Update{
					CallbackQuery: &tgbotapi.CallbackQuery{
						Data: "/users/1/books",
					},
				},
				Params: map[string]string{
					"userID": "1",
				},
			},
			expectError:  true,
			errorMessage: "#1",
			res: &handler.ResponseHandlerFunc{
				MessageConfigs: []*tgbotapi.MessageConfig{
					{Text: "the userID is 1"},
				},
				RedirectRoot: false,
				ReleaseState: true,
				Path:         "",
			},
		},
	} {
		res, err := router.Route(&testCase.input)
		if testCase.expectError && err == nil {
			t.Errorf("expected an error at %d but we got nothing error", i)
		}
		if testCase.expectError && testCase.errorMessage != err.Error() {
			t.Errorf("we expected error message that is %s but we got %s at %d", err.Error(), testCase.errorMessage, i)
		}

		if !testCase.expectError && res.RedirectRoot != testCase.res.RedirectRoot {
			t.Errorf("the redirect root is not matched %d", i)
		}
		if !testCase.expectError && res.ReleaseState != testCase.res.ReleaseState {
			t.Errorf("the release state is not matched at %d", i)
		}

		if !testCase.expectError && !isErrorConfigsMatched(testCase.res.MessageConfigs, res.MessageConfigs) {
			t.Errorf("messageConfigs expected aren't matched with messageConfigs we are given at %d", i)
		}
	}
}

func TestRoutingWithState(t *testing.T) {
	for i, testCase := range []struct {
		requests []struct {
			input        *handler.Context
			expectError  bool
			errorMessage string
			res          *handler.ResponseHandlerFunc
		}
	}{
		{
			requests: []struct {
				input        *handler.Context
				expectError  bool
				errorMessage string
				res          *handler.ResponseHandlerFunc
			}{
				{
					input: &handler.Context{
						UserID: "1",
						Update: &tgbotapi.Update{
							Message: &tgbotapi.Message{
								Text: "/users/1/accounts",
							},
						},
					},
					res: &handler.ResponseHandlerFunc{
						RedirectRoot: false,
						ReleaseState: false,
						MessageConfigs: []*tgbotapi.MessageConfig{
							{
								Text: "input your account:",
							},
						},
					},
					errorMessage: "",
					expectError:  false,
				},
				{
					input: &handler.Context{
						UserID: "1",
						Data: map[string]any{
							"account": "",
						},
						Update: &tgbotapi.Update{
							Message: &tgbotapi.Message{
								Text: "2",
							},
						},
					},
					res: &handler.ResponseHandlerFunc{
						RedirectRoot: true,
						ReleaseState: true,
						MessageConfigs: []*tgbotapi.MessageConfig{
							{
								Text: "account is created the id is 2",
							},
						},
					},
					errorMessage: "",
					expectError:  false,
				},
			},
		},
	} {
		testSequentialRequests(t, testCase.requests, i)
	}
}

func testSequentialRequests(t *testing.T, requests []struct {
	input        *handler.Context
	expectError  bool
	errorMessage string
	res          *handler.ResponseHandlerFunc
}, testCaseIndex int,
) {
	for i, request := range requests {
		res, err := router.Route(request.input)
		if request.expectError && err == nil {
			t.Errorf(
				"we expected an error for testCase %d at request %d but we got nothing ",
				testCaseIndex,
				i,
			)
		}
		if request.expectError && err.Error() != request.errorMessage {
			t.Errorf(
				"we expected an error for testCase %d at request %d error message be %s but we got %s ",
				testCaseIndex,
				i,
				request.errorMessage,
				err.Error(),
			)
		}

		if !request.expectError && res.RedirectRoot != request.res.RedirectRoot {
			t.Errorf("the redirect root is not matched %d", i)
		}
		if !request.expectError && res.ReleaseState != request.res.ReleaseState {
			t.Errorf("the release state is not matched at %d", i)
		}

		if !request.expectError && !isErrorConfigsMatched(request.res.MessageConfigs, res.MessageConfigs) {
			t.Errorf("messageConfigs expected aren't matched with messageConfigs we are given at %d", i)
		}
	}
}
