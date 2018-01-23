package vader

import (
	"context"
	"net/http"
	"regexp"
)

var PathParamKey interface{} = "path_param_key"

// PathParam type provides kv access to path params in the request
type PathParam map[string]string

// Router is simple mux for handling and registering http REST routes
type Router struct {
	routes map[string]*route
}

// NewRouter instantiates a new Router
func NewRouter() *Router {
	return &Router{}
}

// Handle allows to register a new handler for a given request path
//
// the identifer is a short "api_name" used to ID the request
//
// NOTE; the identifier must be unique per Router
//
// Optional; methods - HTTP methods allowed for this handler. If nil, the Router allows all methods to pass thru
func (rr *Router) Handle(identifier, path string, handler http.Handler, methods ...string) (*Router, error) {
	matcher, err := regexp.Compile(path)
	if err != nil {
		return rr, err
	}
	rr.routes[identifier] = &route{matcher: matcher, handler: handler, methods: methods, pathParamNames: matcher.SubexpNames()}
	return rr, nil
}

// ServeHTTP is the http.Handler implementation on Router
func (rr *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request_path := req.URL.EscapedPath()
	for _, r := range rr.routes {
		if r.allowed(req.Method) {
			if params, OK := r.match(request_path); OK {
				// stick the params in a "well known key in the context"
				ctx := context.Background()
				ctx = context.WithValue(ctx, PathParamKey, params)
				rCon := req.WithContext(ctx)
				r.handler.ServeHTTP(w, rCon)
				return
			}
		} else {
			// write a 405
		}
	}
	// write a NotFound Hander
}
