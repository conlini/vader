package vader

import (
	"net/http"
	"regexp"
)

// route is an internal struct representing a single HTTP route
type route struct {
	matcher        *regexp.Regexp // regex matcher for the request
	pathParamNames []string       //list of path param names in the API
	handler        http.Handler   // handler to invoke on match of this route
	methods        []string       // list of allowed HTTP methods
}

// allowed verifies if the HTTP method is permitted on this route
func (r *route) allowed(method string) bool {
	if r.methods != nil {
		for _, m := range r.methods {
			if m == method {
				return true
			}
		}
		return false
	}
	return true
}

// match verifies if the given request path will match this route spec and return Pathparams along with the same
func (r *route) match(path string) (PathParam, bool) {
	matches := r.matcher.FindStringSubmatch(path)
	if matches != nil {
		params := PathParam{}
		for i, name := range r.pathParamNames {
			params[name] = matches[i]
		}
		return params, true
	}

	return nil, false
}
