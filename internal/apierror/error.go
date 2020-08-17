package apierror

import "errors"

// Predefined Errors indentify expected failure conditions.
var (
	// ErrNotFound is used when a specific document is requested but does not exist.
	ErrNotFound = errors.New("document NOT found")

	// ErrInvalID is used when an invalid ID is provided.
	ErrInvalidID = errors.New("_id is NOT in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is NOT allowed")
)
