package vader

import "net/http"

// Middleware is a type used to chain handlers together.
//
// A middleware takes a handler as input and returns another handler that can wrap functionality
// around the original handler
type Middleware func(http.Handler) http.Handler

// Chain is a convenient method to chain multiple handlers together
func Chain(outer Middleware, inner ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(inner) - 1; i >= 0; i-- {
			next = inner[i](next)

		}
		return outer(next)
	}
}

// Finalize is a convenient method to chain multiple middlewares with a handler
func Finalize(handler http.Handler, middlewares ...Middleware) http.Handler {
	mw := Chain(middlewares[len(middlewares)-1], middlewares[0:len(middlewares)-1]...)
	return mw(handler)
}
