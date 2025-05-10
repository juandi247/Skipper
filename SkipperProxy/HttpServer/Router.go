package HttpServer

import (
	"fmt"
	"net/http"
	"strings"
)

// todo: this Router will not be used, the proxy will not handle the request and redirect to route methods
// the proxy will pass the http request through the tcp client connection, and then receive the reponse from lcoalhost

type Handler func(w http.ResponseWriter, r *http.Request)

type Router struct {
	routes          map[string]map[string]Handler
	server          *Server
	notfoundHandler Handler
}

func NewRouter(server *Server) *Router {
	return &Router{
		routes:          make(map[string]map[string]Handler),
		notfoundHandler: defaultNotFoundHandler,
	}
}

func (r *Router) NotFound(handler Handler) {
	r.notfoundHandler = handler
}

func (r *Router) addRoute(method, path string, handler Handler) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]Handler)
	}
	r.routes[method][path] = handler
}

func (r *Router) Any(path string, handler Handler) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodHead,
		http.MethodConnect,
		http.MethodTrace,
	}

	for _, method := range methods {
		r.addRoute(method, path, handler)
	}
}
func (r *Router) ServeFavicon() {
	r.GET("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		// Responde con un código de estado 200 OK para evitar que se procese más
		w.WriteHeader(http.StatusOK)
	})
}

func (r *Router) GET(path string, handler Handler) {
	r.addRoute(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler Handler) {
	r.addRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler Handler) {
	r.addRoute(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler Handler) {
	r.addRoute(http.MethodDelete, path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, err := r.findHandler(req.Method, req.URL.Path)
	if err != nil {
		r.notfoundHandler(w, req)
		return
	}
	// if r.server != nil {
	//  handler = r.server.ApplyMiddleware(handler)
	// }
	handler(w, req)
}

func (r *Router) findHandler(method, path string) (Handler, error) {
	if methodRoutes, ok := r.routes[method]; ok {
		if handler, ok := methodRoutes[path]; ok {
			return handler, nil
		}

		for routePath, handler := range methodRoutes {
			if isWildcardMatch(routePath, path) {
				return handler, nil
			}
		}
	}
	return nil, fmt.Errorf("no handler found for %s %s", method, path)
}

func isWildcardMatch(routePath, requestPath string) bool {
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(routeParts) != len(requestParts) {
		return false
	}

	for i, part := range routeParts {
		if part == "*" {
			continue
		}
		if part != requestParts[i] {
			return false
		}
	}

	return true
}

func defaultNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}
