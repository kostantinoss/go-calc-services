package internal

import (
	"log"
	"net/http"
)

// Router used to determine where to route each incomming request
type Router struct {
	logger *log.Logger
	routes map[string]http.HandlerFunc
}

func (r *Router) Init() {
	r.routes = make(map[string]http.HandlerFunc)
}

func (r *Router) AddRoute(path string, handler http.HandlerFunc) {
	r.routes[path] = handler
}

func (r *Router) matchRoute(path string) (http.HandlerFunc, bool) {
	handler, exists := r.routes[path]
	if !exists {
		r.logger.Printf("Route not found: %s", path)
		return nil, false
	}
	return handler, true
}
