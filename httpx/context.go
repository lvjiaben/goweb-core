package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/lvjiaben/goweb-core/errorsx"
)

const requestIDKey = "request_id"

type Context struct {
	Writer       http.ResponseWriter
	Request      *http.Request
	engine       *Engine
	route        *Route
	values       map[string]any
	status       int
	wroteHeader  bool
	responseCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		values:  make(map[string]any),
		status:  http.StatusOK,
	}
}

func (c *Context) Route() *Route {
	return c.route
}

func (c *Context) Set(key string, value any) {
	c.values[key] = value
}

func (c *Context) Get(key string) (any, bool) {
	value, ok := c.values[key]
	return value, ok
}

func (c *Context) MustGetString(key string) string {
	value, ok := c.values[key]
	if !ok {
		return ""
	}
	str, _ := value.(string)
	return str
}

func (c *Context) RequestID() string {
	if value := c.ResponseHeader("X-Request-Id"); value != "" {
		return value
	}
	return c.MustGetString(requestIDKey)
}

func (c *Context) ResponseHeader(key string) string {
	return c.Writer.Header().Get(key)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status() int {
	return c.status
}

func (c *Context) Query(key string) string {
	return strings.TrimSpace(c.Request.URL.Query().Get(key))
}

func (c *Context) QueryInt64(key string) (int64, error) {
	value := c.Query(key)
	if value == "" {
		return 0, nil
	}
	return strconv.ParseInt(value, 10, 64)
}

func (c *Context) BindJSON(out any) error {
	if c.Request.Body == nil {
		return errorsx.Wrap(errors.New("empty request body"), http.StatusBadRequest, errorsx.CodeBadRequest, "empty request body")
	}
	defer c.Request.Body.Close()

	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(out); err != nil {
		if errors.Is(err, io.EOF) {
			return errorsx.Wrap(err, http.StatusBadRequest, errorsx.CodeBadRequest, "empty request body")
		}
		return errorsx.Wrap(err, http.StatusBadRequest, errorsx.CodeBadRequest, "invalid json body")
	}
	return nil
}

func (c *Context) MultipartFormFile(field string) (multipartFile multipart.File, header *multipart.FileHeader, err error) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		return nil, nil, errorsx.Wrap(err, http.StatusBadRequest, errorsx.CodeBadRequest, "invalid multipart form")
	}
	file, fileHeader, err := c.Request.FormFile(field)
	if err != nil {
		return nil, nil, errorsx.Wrap(err, http.StatusBadRequest, errorsx.CodeBadRequest, "missing file")
	}
	return file, fileHeader, nil
}

func (c *Context) JSON(status int, code int, msg string, data any) {
	c.status = status
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	c.WriteHeader(status)
	_ = json.NewEncoder(c.Writer).Encode(Response{
		Code:      code,
		Msg:       msg,
		Data:      normalizeData(data),
		RequestID: c.RequestID(),
	})
}

func (c *Context) Success(data any) {
	c.JSON(http.StatusOK, 0, "ok", data)
}

func (c *Context) Fail(status int, code int, msg string, data any) {
	c.JSON(status, code, msg, data)
}

func (c *Context) Error(err error) {
	appErr := errorsx.From(err)
	c.JSON(appErr.HTTPStatus, appErr.Code, appErr.Msg, map[string]any{})
}

func (c *Context) BadRequest(msg string) {
	c.Fail(http.StatusBadRequest, errorsx.CodeBadRequest, msg, map[string]any{})
}

func (c *Context) Unauthorized(msg string) {
	c.Fail(http.StatusUnauthorized, errorsx.CodeUnauthorized, msg, map[string]any{})
}

func (c *Context) Forbidden(msg string) {
	c.Fail(http.StatusForbidden, errorsx.CodeForbidden, msg, map[string]any{})
}

func (c *Context) NotFound(msg string) {
	c.Fail(http.StatusNotFound, errorsx.CodeNotFound, msg, map[string]any{})
}

func (c *Context) WriteHeader(status int) {
	if c.wroteHeader {
		return
	}
	c.Writer.WriteHeader(status)
	c.wroteHeader = true
	c.status = status
}

func (c *Context) ClientIP() string {
	headers := []string{"X-Forwarded-For", "X-Real-Ip"}
	for _, header := range headers {
		if raw := strings.TrimSpace(c.Request.Header.Get(header)); raw != "" {
			if header == "X-Forwarded-For" {
				parts := strings.Split(raw, ",")
				if len(parts) > 0 {
					return strings.TrimSpace(parts[0])
				}
			}
			return raw
		}
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(c.Request.RemoteAddr)
}
