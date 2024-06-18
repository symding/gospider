package gospider

import "net/http"

// Request 请求
type Request struct {
	Url         string
	Proxy       string
	Cookie      []*http.Cookie
	MaxRetryNum int
	Meta        map[string]any
	Headers     map[string]string
}
