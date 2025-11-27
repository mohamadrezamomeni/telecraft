package tree

import (
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/telecraft/handler"
)

func TestResovingHandler(t *testing.T) {
	root := New("", nil)

	root.Set(strings.Split("books/:id", "/"), func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			MessageConfigs: []*tgbotapi.MessageConfig{
				{Text: "books/333"},
			},
		}, nil
	})

	root.Set(strings.Split("books", "/"), func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			MessageConfigs: []*tgbotapi.MessageConfig{
				{Text: "books"},
			},
		}, nil
	})

	root.Set(strings.Split("books/:id/authers", "/"), func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			MessageConfigs: []*tgbotapi.MessageConfig{
				{Text: "authers"},
			},
		}, nil
	})

	root.Set(strings.Split("books/:id/authers/:autherID", "/"), func(u *handler.Context) (*handler.ResponseHandlerFunc, error) {
		return &handler.ResponseHandlerFunc{
			MessageConfigs: []*tgbotapi.MessageConfig{
				{Text: "authers/111"},
			},
		}, nil
	})

	for i, testCase := range []struct {
		input  []string
		output struct {
			isExsted    bool
			params      map[string]string
			handlertext string
		}
	}{
		{
			input: []string{
				"books",
				"333",
			},
			output: struct {
				isExsted    bool
				params      map[string]string
				handlertext string
			}{
				isExsted: true,
				params: map[string]string{
					"id": "333",
				},
				handlertext: "books/333",
			},
		},
		{
			input: []string{
				"books",
			},
			output: struct {
				isExsted    bool
				params      map[string]string
				handlertext string
			}{
				isExsted:    true,
				params:      map[string]string{},
				handlertext: "books",
			},
		},
		{
			input: []string{
				"books",
				"444",
				"authers",
			},
			output: struct {
				isExsted    bool
				params      map[string]string
				handlertext string
			}{
				params: map[string]string{
					"id": "444",
				},
				isExsted:    true,
				handlertext: "authers",
			},
		},
		{
			input: []string{
				"books",
				"444",
				"authers",
				"111",
			},
			output: struct {
				isExsted    bool
				params      map[string]string
				handlertext string
			}{
				params: map[string]string{
					"id":       "444",
					"autherID": "111",
				},
				isExsted:    true,
				handlertext: "authers/111",
			},
		},
		{
			input: []string{
				"books",
				"444",
				"users",
				"111",
			},
			output: struct {
				isExsted    bool
				params      map[string]string
				handlertext string
			}{
				isExsted:    false,
				params:      nil,
				handlertext: "authers/111",
			},
		},
		{
			input: []string{
				"users",
				"444",
				"authers",
				"111",
			},
			output: struct {
				isExsted    bool
				params      map[string]string
				handlertext string
			}{
				isExsted:    false,
				params:      nil,
				handlertext: "authers/111",
			},
		},
	} {
		node, params := root.MatchPath(testCase.input)
		if !testCase.output.isExsted && node != nil {
			t.Errorf("we expected res would be nil  at %d", i)
		}

		if testCase.output.isExsted && node == nil {
			t.Errorf("we expected res wouldn't be nil but we got noting at %d", i)
		}
		if node == nil {
			continue
		}

		res, _ := node.Handler(&handler.Context{})

		messages := res.MessageConfigs
		if len(messages) != 1 {
			t.Errorf("we expected the lengh of message config be 1 but we got %d message config", len(messages))
		}

		msg := messages[0]

		if msg.Text != testCase.output.handlertext {
			t.Errorf(
				"we expected the message would be %s but we got %s at %d",
				testCase.output.handlertext,
				msg.Text,
				i,
			)
		}

		if len(params) != len(testCase.output.params) {
			t.Errorf("error to compare output params at %d", i)
		}

		for k, v := range params {
			if ov, ok := testCase.output.params[k]; !ok || ov != v {
				t.Errorf("error to compare output params as %d", i)
			}
		}
	}
}
