package httpx

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
}

func normalizeData(data any) any {
	if data == nil {
		return map[string]any{}
	}
	return data
}
