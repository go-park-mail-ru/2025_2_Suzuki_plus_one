package http

import (
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/pkg/logger"
)

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	PATCH   Method = "PATCH"
	OPTIONS Method = "OPTIONS"
	HEAD    Method = "HEAD"
)

var AllowedMethods = map[Method]struct{}{
	GET:     {},
	POST:    {},
	PUT:     {},
	DELETE:  {},
	PATCH:   {},
	OPTIONS: {},
	HEAD:    {},
}

type Route map[Method]http.Handler

// func NewRoute() Route {

// }

type Routes map[string]Route

type Router struct {
	middleware      []func(http.Handler) http.Handler // List of middleware functions
	logger          logger.Logger                     // Logger instance
	prefix          string                            // Router prefix
	routes          Routes                            // Keep routes like [path: [method: handler]]
	notFoundHandler http.Handler                      // Handler for not found routes
}

// Add middleware to the router
func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.middleware = append(r.middleware, mw)
}

// Apply all middleware to the handler in the order they were added
func (r *Router) ApplyMiddleware(handler http.Handler) http.Handler {
	h := handler
	for i := len(r.middleware) - 1; i >= 0; i-- {
		h = r.middleware[i](h)
	}
	return h
}

// Add a new route to the router
func (r *Router) Add(method Method, path string, handler http.Handler) {
	if _, ok := AllowedMethods[method]; !ok {
		panic("Cannot add route with unsupported method: " + string(method))
	}

	// Check for trailing spaces and slashes
	if path[len(path)-1] == ' ' || path[len(path)-1] == '/' {
		panic("Cannot add route with trailing space or slash in path: " + path)
	}

	path = r.prefix + path
	if strings.Count(path, "//") > 0 {
		panic("Cannot add route with double slashes in path: " + path)
	}

	if _, ok := r.routes[path]; !ok {
		r.routes[path] = make(Route)
	}

	if _, exists := r.routes[path][method]; exists {
		panic("Can't add another route for " + string(method) + " " + path)
	}
	r.routes[path][method] = handler
}

func (r *Router) Handle(method Method, path string, handlerFunc http.HandlerFunc) {
	r.Add(method, path, handlerFunc)
}

// Get the handler for a given method and path
func (r *Router) getHandler(method Method, path string) http.Handler {
	route, ok := r.routes[path]
	if !ok {
		return r.notFoundHandler
	}

	handler, ok := route[method]
	if !ok {
		return r.notFoundHandler
	}

	return handler
}

// Get the handler for a given method and path
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := Method(req.Method)
	r.logger.Debug("ServeHTTP", r.logger.ToString("path", path), r.logger.ToString("method", string(method)))

	handler := r.getHandler(method, path)
	handler = r.ApplyMiddleware(handler)

	handler.ServeHTTP(w, req)
}

func NewRouter(prefix string, logger logger.Logger) *Router {
	return &Router{
		logger:          logger,
		middleware:      []func(http.Handler) http.Handler{},
		prefix:          prefix,
		routes:          make(Routes),
		notFoundHandler: http.NotFoundHandler(),
	}

}
