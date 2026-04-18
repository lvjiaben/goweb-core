package httpx

import (
	"log/slog"
	"net/http"
	"path"
	"strings"
	"sync"
)

type HandlerFunc func(*Context)

type Middleware func(HandlerFunc) HandlerFunc

type Route struct {
	Method         string
	Path           string
	Name           string
	PermissionCode string
	Handler        HandlerFunc
}

type RouteOption func(*Route)

type Engine struct {
	logger      *slog.Logger
	mu          sync.RWMutex
	routes      map[string]*Route
	middlewares []Middleware
	notFound    HandlerFunc
}

type Group struct {
	engine      *Engine
	prefix      string
	middlewares []Middleware
}

func NewEngine(logger *slog.Logger) *Engine {
	return &Engine{
		logger:   logger,
		routes:   make(map[string]*Route),
		notFound: defaultNotFound,
	}
}

func (e *Engine) Use(middlewares ...Middleware) {
	e.middlewares = append(e.middlewares, middlewares...)
}

func (e *Engine) Group(prefix string, middlewares ...Middleware) *Group {
	return &Group{
		engine:      e,
		prefix:      cleanPath(prefix),
		middlewares: append([]Middleware{}, middlewares...),
	}
}

func (e *Engine) GET(routePath string, handler HandlerFunc, opts ...RouteOption) {
	e.Handle(http.MethodGet, routePath, handler, opts...)
}

func (e *Engine) POST(routePath string, handler HandlerFunc, opts ...RouteOption) {
	e.Handle(http.MethodPost, routePath, handler, opts...)
}

func (e *Engine) Handle(method string, routePath string, handler HandlerFunc, opts ...RouteOption) {
	route := &Route{
		Method:  strings.ToUpper(method),
		Path:    cleanPath(routePath),
		Handler: handler,
	}
	for _, opt := range opts {
		opt(route)
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.routes[routeKey(route.Method, route.Path)] = route
}

func (e *Engine) Routes() []*Route {
	e.mu.RLock()
	defer e.mu.RUnlock()

	out := make([]*Route, 0, len(e.routes))
	for _, route := range e.routes {
		out = append(out, route)
	}
	return out
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.mu.RLock()
	route := e.routes[routeKey(r.Method, cleanPath(r.URL.Path))]
	e.mu.RUnlock()

	ctx := newContext(w, r)
	ctx.engine = e
	ctx.route = route

	finalHandler := e.notFound
	if route != nil {
		finalHandler = route.Handler
	}

	handler := chain(finalHandler, e.middlewares...)
	handler(ctx)
}

func (g *Group) Use(middlewares ...Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) Group(prefix string, middlewares ...Middleware) *Group {
	return &Group{
		engine:      g.engine,
		prefix:      joinPath(g.prefix, prefix),
		middlewares: append(append([]Middleware{}, g.middlewares...), middlewares...),
	}
}

func (g *Group) GET(routePath string, handler HandlerFunc, opts ...RouteOption) {
	g.Handle(http.MethodGet, routePath, handler, opts...)
}

func (g *Group) POST(routePath string, handler HandlerFunc, opts ...RouteOption) {
	g.Handle(http.MethodPost, routePath, handler, opts...)
}

func (g *Group) Handle(method string, routePath string, handler HandlerFunc, opts ...RouteOption) {
	handler = chain(handler, g.middlewares...)
	g.engine.Handle(method, joinPath(g.prefix, routePath), handler, opts...)
}

func WithName(name string) RouteOption {
	return func(route *Route) {
		route.Name = name
	}
}

func WithPermission(code string) RouteOption {
	return func(route *Route) {
		route.PermissionCode = strings.TrimSpace(code)
	}
}

func chain(final HandlerFunc, middlewares ...Middleware) HandlerFunc {
	handler := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func routeKey(method string, routePath string) string {
	return strings.ToUpper(method) + ":" + cleanPath(routePath)
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	cleaned := path.Clean("/" + strings.TrimSpace(p))
	if cleaned == "." {
		return "/"
	}
	return cleaned
}

func joinPath(parts ...string) string {
	result := "/"
	for _, part := range parts {
		if strings.TrimSpace(part) == "" || part == "/" {
			continue
		}
		result = path.Join(result, part)
	}
	return cleanPath(result)
}

func defaultNotFound(c *Context) {
	c.NotFound("route not found")
}
