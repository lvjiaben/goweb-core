package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodAwareRouting(t *testing.T) {
	engine := NewEngine(slog.Default())
	engine.GET("/same", func(c *Context) {
		c.Success(map[string]any{"handler": "get"})
	})
	engine.POST("/same", func(c *Context) {
		c.Success(map[string]any{"handler": "post"})
	})

	getResp := httptest.NewRecorder()
	engine.ServeHTTP(getResp, httptest.NewRequest(http.MethodGet, "/same", nil))
	if getResp.Code != http.StatusOK {
		t.Fatalf("GET status = %d", getResp.Code)
	}
	if got := responseData(t, getResp.Body.Bytes())["handler"]; got != "get" {
		t.Fatalf("GET handler = %v", got)
	}

	postResp := httptest.NewRecorder()
	engine.ServeHTTP(postResp, httptest.NewRequest(http.MethodPost, "/same", nil))
	if postResp.Code != http.StatusOK {
		t.Fatalf("POST status = %d", postResp.Code)
	}
	if got := responseData(t, postResp.Body.Bytes())["handler"]; got != "post" {
		t.Fatalf("POST handler = %v", got)
	}
}

func TestMethodMismatchReturns405(t *testing.T) {
	engine := NewEngine(slog.Default())
	engine.POST("/captcha", func(c *Context) {
		c.Success(map[string]any{"handler": "post"})
	})

	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/captcha", nil))

	if resp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusMethodNotAllowed)
	}
	if allow := resp.Header().Get("Allow"); allow != http.MethodPost {
		t.Fatalf("Allow = %q, want %q", allow, http.MethodPost)
	}

	var body Response
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != http.StatusMethodNotAllowed {
		t.Fatalf("code = %d, want %d", body.Code, http.StatusMethodNotAllowed)
	}
	if body.Msg != "method not allowed" {
		t.Fatalf("msg = %q", body.Msg)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T", body.Data)
	}
	if data["method"] != http.MethodGet {
		t.Fatalf("method data = %v", data["method"])
	}
}

func TestUnknownPathStillReturns404(t *testing.T) {
	engine := NewEngine(slog.Default())
	engine.POST("/known", func(c *Context) {
		c.Success(nil)
	})

	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, httptest.NewRequest(http.MethodGet, "/missing", nil))
	if resp.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusNotFound)
	}
}

func responseData(t *testing.T, raw []byte) map[string]any {
	t.Helper()

	var body Response
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T", body.Data)
	}
	return data
}
