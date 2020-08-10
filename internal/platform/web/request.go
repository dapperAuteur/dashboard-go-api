package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Decode looks for a JSON document in the request body and unmarshals it into val.
func Decode(r *http.Request, val interface{}) error {

	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return errors.Wrapf(err, "decoding request body")
	}
	return nil
}
