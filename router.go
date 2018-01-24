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
	return &Router{
		routes: make(map[string]*route),
	}
}

// Handle allows to register a new handler for a given request path
//
// the identifer is a short "api_name" used to ID the request
//
// NOTE; the identifier must be unique per Router
//
// Optional; methods - HTTP methods allowed for this handler. If nil, the Router allows all methods to pass thru
func (rr *Router) Handle(identifier, path string, handler http.Handler, methods ...string) error {
	// if path does not begin with '/' add it as an optional entry
	// add a leading '^' and a trailing '/?$' for regex compliation
	//
	// this means registering a path of the form "/a/b" will be compiled as "^/a/b/?$" to ensure right match
	if path[0] != '/' {
		path = "/" + path
	}

	n := len(path)
	if path[n-1] != '/' {
		path = path + "/"
	}
	path = "^" + path + "?$"
	matcher, err := regexp.Compile(path)
	if err != nil {
		return err
	}
	rr.routes[identifier] = &route{matcher: matcher, handler: handler, methods: methods, pathParamNames: matcher.SubexpNames()}
	return nil
}

// ServeHTTP is the http.Handler implementation on Router
func (rr *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request_path := req.URL.EscapedPath()
	for _, r := range rr.routes {

		if params, OK := r.match(request_path); OK {
			if r.allowed(req.Method) {
				// stick the params in a "well known key in the context"
				ctx := context.Background()
				ctx = context.WithValue(ctx, PathParamKey, params)
				rCon := req.WithContext(ctx)
				r.handler.ServeHTTP(w, rCon)
				return
			} else {
				// write a 405
				handle405(w, req)
				return
			}

		}
	}
	// write a NotFound Hander
	handleNotFound(w, req)
}
