package web

// ErrorResponse how we respond to clients when something goes wrong. May add ErrorCodes later
type ErrorResponse struct {
	Error string `json:"error"`
}

// Error is used to add web information to request error.
type Error struct {
	Err    error
	Status int
}

// NewRequestError is used when a known error condition is encountered.
func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
