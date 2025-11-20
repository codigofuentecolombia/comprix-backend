package dto

type ServerRequest struct {
	URL     string   `json:"url"`
	Method  string   `json:"method"`
	Params  string   `json:"params"`
	Headers []string `json:"headers"`
}
