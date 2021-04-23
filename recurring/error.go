package recurring

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shortcut/go-vipps"
	"github.com/shortcut/go-vipps/internal"
)

// ErrRecurring represents errors returned from the Vipps Recurring Payments
// API.
type ErrRecurring []APIError

// APIError represents a single error returned from the Vipps
// Recurring Payments API.
type APIError struct {
	Field     string `json:"field"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	ContextID string `json:"contextId"`
}

func (e ErrRecurring) Error() string {
	s := []string{"vipps:"}
	if len(e) > 1 {
		s = append(s, "multiple errors:")
	}
	for _, e := range e {
		s = append(s, fmt.Sprintf("field %s: %s (code %s)", e.Field, e.Message, e.Code))
	}
	return strings.Join(s, " ")
}

func wrapErr(err error) error {
	if err, ok := err.(internal.HTTPError); ok {
		var wrappedErr ErrRecurring
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
