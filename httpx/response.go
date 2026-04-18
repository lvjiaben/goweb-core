package httpx

type Response struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
}

func normalizeData(data any) any {
	if data == nil {
		return map[string]any{}
	}
	return data
}
