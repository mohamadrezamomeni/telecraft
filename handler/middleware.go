package handler

func ApplyMiddlewares(h HandlerFunc, ms ...Middleware) HandlerFunc {
	finalHandler := h
	for i := len(ms) - 1; i >= 0; i-- {
		finalHandler = ms[i](finalHandler)
	}
	return finalHandler
}
