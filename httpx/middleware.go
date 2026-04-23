package httpx

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lvjiaben/goweb-core/errorsx"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAgeSeconds    int
}

func Recover(logger *slog.Logger) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			defer func() {
				if recovered := recover(); recovered != nil {
					if logger != nil {
						logger.Error("panic recovered", "panic", recovered, "path", c.Request.URL.Path, "request_id", c.RequestID())
					}
					c.Fail(http.StatusInternalServerError, errorsx.CodeInternal, "internal server error", map[string]any{})
				}
			}()
			next(c)
		}
	}
}

func RequestID() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			requestID := strings.TrimSpace(c.Request.Header.Get("X-Request-Id"))
			if requestID == "" {
				requestID = uuid.NewString()
			}
			c.Set(requestIDKey, requestID)
			c.SetHeader("X-Request-Id", requestID)
			next(c)
		}
	}
}

func Logger(logger *slog.Logger) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			start := time.Now()
			next(c)
			if logger == nil {
				return
			}
			logger.Info(
				"http request",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"status", c.Status(),
				"duration", time.Since(start).String(),
				"request_id", c.RequestID(),
				"ip", c.ClientIP(),
				"permission_code", routePermission(c.Route()),
			)
		}
	}
}

func CORS(cfg CORSConfig) Middleware {
	allowMethods := joinCSVOrDefault(cfg.AllowMethods, []string{http.MethodGet, http.MethodPost, http.MethodOptions})
	allowHeaders := joinCSVOrDefault(cfg.AllowHeaders, []string{"Authorization", "Content-Type", "X-Request-Id"})
	exposeHeaders := strings.Join(cfg.ExposeHeaders, ", ")

	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			setCORSHeaders(c, cfg, allowMethods, allowHeaders, exposeHeaders)
			if c.Request.Method == http.MethodOptions {
				c.WriteHeader(http.StatusNoContent)
				return
			}
			next(c)
		}
	}
}

func setCORSHeaders(c *Context, cfg CORSConfig, allowMethods string, allowHeaders string, exposeHeaders string) {
	origin := strings.TrimSpace(c.Request.Header.Get("Origin"))
	allowOrigin := "*"
	if len(cfg.AllowOrigins) > 0 {
		allowOrigin = matchOrigin(origin, cfg.AllowOrigins)
	}
	if allowOrigin != "" {
		c.SetHeader("Access-Control-Allow-Origin", allowOrigin)
	}
	c.SetHeader("Access-Control-Allow-Methods", allowMethods)
	c.SetHeader("Access-Control-Allow-Headers", allowHeaders)
	if exposeHeaders != "" {
		c.SetHeader("Access-Control-Expose-Headers", exposeHeaders)
	}
	if cfg.AllowCredentials {
		c.SetHeader("Access-Control-Allow-Credentials", "true")
	}
	if cfg.MaxAgeSeconds > 0 {
		c.SetHeader("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAgeSeconds))
	}
}

func matchOrigin(origin string, allowOrigins []string) string {
	if origin == "" {
		if len(allowOrigins) > 0 {
			return allowOrigins[0]
		}
		return "*"
	}
	for _, item := range allowOrigins {
		if item == "*" || strings.EqualFold(strings.TrimSpace(item), origin) {
			return origin
		}
	}
	return ""
}

func joinCSVOrDefault(values []string, defaults []string) string {
	if len(values) == 0 {
		return strings.Join(defaults, ", ")
	}
	return strings.Join(values, ", ")
}

func routePermission(route *Route) string {
	if route == nil {
		return ""
	}
	return route.PermissionCode
}
