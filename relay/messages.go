package relay

type MessageIn struct {
	ID      string            `json:"id"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type MessageOut struct {
	ID         string            `json:"id"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}
