package ecom

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shortcut/go-vipps"
	"github.com/shortcut/go-vipps/internal"
)

// ErrEcom represents errors returned from the Vipps Ecom API.
type ErrEcom []APIError

// APIError represents a single error returned from the Vipps Ecom API.
type APIError struct {
	Group   string `json:"errorGroup"`
	Message string `json:"errorMessage"`
	Code    string `json:"errorCode"`
}

func (e ErrEcom) Error() string {
	s := []string{"vipps:"}
	if len(e) > 1 {
		s = append(s, "multiple errors:")
	}
	for _, e := range e {
		s = append(s, fmt.Sprintf("[%s] %s (code %s)", e.Group, e.Message, e.Code))
	}
	return strings.Join(s, " ")
}

func wrapErr(err error) error {
	if err, ok := err.(internal.HTTPError); ok {
		var wrappedErr ErrEcom
		unmarshalErr := json.Unmarshal(err.Body, &wrappedErr)
		if unmarshalErr != nil {
			return vipps.ErrUnexpectedResponse{
				Body:   err.Body,
				Status: err.Status,
			}
		}
		return wrappedErr
	}
	return err
}
