package telecraft

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mohamadrezamomeni/telecraft/handler"
	"github.com/mohamadrezamomeni/telecraft/pkg/telecrafterror"
	"github.com/mohamadrezamomeni/telecraft/router"
	"github.com/mohamadrezamomeni/telecraft/state"
)

type TeleCraft struct {
	Router           *router.Router
	telecraftOptions *TeleCraftOptions
	bot              *tgbotapi.BotAPI
}

type TeleCraftOptions struct {
	RepoType      string
	DefaultRoute  string
	maxGoroutines int
	Timeout       int
	Token         string
	errorMessage  string
}

func New(telecraftOptions *TeleCraftOptions) *TeleCraft {
	scope := "new.Telecraft"

	bot, err := tgbotapi.NewBotAPI(telecraftOptions.Token)
	if err != nil {
		panic(telecrafterror.Wrap(err).Scope(scope).BadRequest().Errorf("error to initialize bot"))
	}

	stateRepo, err := state.NewRepository(telecraftOptions.RepoType)
	if err != nil {
		panic(err.Error())
	}

	return &TeleCraft{
		bot:              bot,
		Router:           router.New(telecraftOptions.RepoType, stateRepo),
		telecraftOptions: telecraftOptions,
	}
}

func (t *TeleCraft) Serve() {
	u := tgbotapi.NewUpdate(t.telecraftOptions.Timeout)
	updates := t.bot.GetUpdatesChan(u)

	limiter := make(chan struct{}, t.telecraftOptions.maxGoroutines)
	for update := range updates {
		go t.handleRequest(&update, limiter)
	}
}

func (t *TeleCraft) handleRequest(update *tgbotapi.Update, limiter chan struct{}) {
	limiter <- struct{}{}

	context := &handler.Context{
		Update: update,
	}

	// TODO: need to handle this error
	res, _ := t.Router.Route(context)

	if res != nil {
		t.send(res, context)
	}

	<-limiter
}

func (t *TeleCraft) send(res *handler.ResponseHandlerFunc, context *handler.Context) {
	for _, messageConfig := range res.MessageConfigs {
		t.bot.Send(messageConfig)
	}
}
