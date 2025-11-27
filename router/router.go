package router

import (
	"strings"
	"time"

	"github.com/mohamadrezamomeni/telecraft/handler"
	"github.com/mohamadrezamomeni/telecraft/pkg/telecrafterror"
	"github.com/mohamadrezamomeni/telecraft/state"
	"github.com/mohamadrezamomeni/telecraft/tree"
)

type StateRepository interface {
	Set(string, *state.State) error
	Get(string) (*state.State, bool)
	Delete(string) error
}

type Router struct {
	data              *tree.Tree
	defaultRoute      string
	globalMiddlewares []handler.Middleware
	stateRepo         StateRepository
}

func New(defaultRoute string, stateRepo StateRepository) *Router {
	return &Router{
		data:         tree.New("", nil),
		defaultRoute: defaultRoute,
		stateRepo:    stateRepo,
	}
}

func (r *Router) Register(
	path string,
	h handler.HandlerFunc,
	ms ...handler.Middleware,
) {
	finalHandler := handler.ApplyMiddlewares(h, ms...)

	if r.globalMiddlewares != nil && len(r.globalMiddlewares) > 0 {
		finalHandler = handler.ApplyMiddlewares(finalHandler, r.globalMiddlewares...)
	}

	r.data.Set(
		r.makeHierarchyPath(path),
		finalHandler,
	)
}

func (r *Router) makeHierarchyPath(path string) []string {
	return strings.Split(path, "/")
}

func (r *Router) SetGlobalMiddlewares(middlewwares ...handler.Middleware) {
	r.globalMiddlewares = middlewwares
}

func (r *Router) Route(context *handler.Context) (*handler.ResponseHandlerFunc, error) {
	var res *handler.ResponseHandlerFunc
	var err error

	switch {
	case context.CallbackQuery != nil:
		res, err = r.callbackQuery(context)
	case context.Message != nil:
		res, err = r.message(context)
	}

	if res == nil {
		res, _ = r.RootHandler(context)
	}

	if res != nil && res.ReleaseState {
		r.stateRepo.Delete(context.UserID)
	} else if res != nil && (len(res.Path) > 0 || len(res.Data) > 0) {
		r.stateRepo.Set(context.UserID, &state.State{
			Data:       res.Data,
			Path:       res.Path,
			Expiration: time.Now().Add(2 * 60 * time.Second),
		})
	}
	return res, err
}

func (r *Router) callbackQuery(context *handler.Context) (*handler.ResponseHandlerFunc, error) {
	text := context.CallbackQuery.Data
	return r.getResponse(text, context)
}

func (r *Router) message(context *handler.Context) (*handler.ResponseHandlerFunc, error) {
	text := context.Message.Text
	return r.getResponse(text, context)
}

func (r *Router) getResponse(text string, context *handler.Context) (*handler.ResponseHandlerFunc, error) {
	var res *handler.ResponseHandlerFunc
	var err error
	if err != nil {
		return nil, err
	}

	if r.isPath(text) {
		r.stateRepo.Delete(context.UserID)
		path := r.getPathFromText(text)
		res, err = r.routeFromText(path, context)
	}

	if res == nil && err == nil {
		res, err = r.getResponseFromState(context)
	}

	return res, err
}

func (r *Router) getResponseFromState(context *handler.Context) (*handler.ResponseHandlerFunc, error) {
	scope := "telegram.router.getResponseFromState"

	state, isExist := r.stateRepo.Get(context.UserID)
	if !isExist {
		return nil, telecrafterror.Scope(scope).ErrorWrite()
	}

	handler, params := r.getHandlerWithParam(state.Path)

	context.Params = params

	res, err := handler(context)
	if err != nil {
		res, _ := r.RootHandler(context)
		return res, err
	}

	return res, nil
}

func (r *Router) getPathFromText(path string) string {
	return path[1:]
}

func (r *Router) routeFromText(path string, context *handler.Context) (*handler.ResponseHandlerFunc, error) {
	handler, params := r.getHandlerWithParam(path)

	r.enrichContext(context, params)

	res, err := handler(context)
	if err != nil {
		res, _ = r.RootHandler(context)
		return res, err
	}
	return res, nil
}

func (r *Router) RootHandler(context *handler.Context) (*handler.ResponseHandlerFunc, error) {
	handler, params := r.getHandlerWithParam(r.defaultRoute)
	r.enrichContext(context, params)
	res, err := handler(context)
	return res, err
}

func (r *Router) getHandlerWithParam(path string) (handler.HandlerFunc, map[string]string) {
	if node, params := r.data.MatchPath(r.makeHierarchyPath(path)); node != nil {
		return node.Handler, params
	}

	return r.RootHandler, nil
}

func (r *Router) enrichContext(context *handler.Context, params map[string]string) {
	context.Params = params
}

func (r *Router) isPath(text string) bool {
	action := byte('/')

	if text[0] == action {
		return true
	}

	return false
}
