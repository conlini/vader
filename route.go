package vader

import (
	"net/http"
	"regexp"
)

type route struct {
	matcher        *regexp.Regexp
	pathParamNames []string
	handler        http.Handler
	methods        []string
}

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
