package web

// ErrorResponse how we respond to clients when something goes wrong. May add ErrorCodes later
type ErrorResponse struct {
	Error string `json:"error"`
}
