package vader

import (
	"context"
	"net/http"
	"regexp"
)

var PathParamKey interface{} = "path_param_key"

type PathParam map[string]string

type Router struct {
	routes map[string]*route
}

func NewRouter() *Router {
	return &Router{}
}

func (rr *Router) Handle(identifier, path string, handler http.Handler, methods ...string) (*Router, error) {
	matcher, err := regexp.Compile(path)
	if err != nil {
		return rr, err
	}
	rr.routes[identifier] = &route{matcher: matcher, handler: handler, methods: methods, pathParamNames: matcher.SubexpNames()}
	return rr, nil
}

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
