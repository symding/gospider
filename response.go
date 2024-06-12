package gospider

// Response 响应
type Response struct {
	Request Request
	Error   error
	Content string
	Meta    map[string]any
	StatusCode int
}
