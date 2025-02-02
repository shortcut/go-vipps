package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shortcut/go-vipps/logging"
)

// HTTPError hold the HTTP-response and Status-code.
type HTTPError struct {
	Body   []byte
	Status int
}

// Error returns the error as a string.
func (e HTTPError) Error() string {
	return fmt.Sprintf("request failed with status: %d", e.Status)
}

// APIClient holds a HTTP-client and a logger.
type APIClient struct {
	L logging.Logger
	C *http.Client
}

// NewRequest creates a new HTTP request, with a JSON-body.
func (c *APIClient) NewRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

// Do runs a HTTP request, and unmarshalls the response into `v`.
func (c *APIClient) Do(req *http.Request, v interface{}) error {
	now := time.Now()
	resp, err := c.C.Do(req)
	logArguments := []logging.LogArgument{logging.NewArg("method", req.Method), logging.NewArg("url", req.URL), logging.NewArg("durationMS", time.Since(now).Milliseconds())}
	if err != nil {
		c.L.Error(req.Context(), "error executing Vipps HTTP request", logArguments...)
		return err
	}
	c.L.Info(req.Context(), "executed Vipps HTTP request", logArguments...)
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.L.Error(req.Context(), "error reading response-body", logArguments...)
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		return HTTPError{
			Body:   body,
			Status: resp.StatusCode,
		}
	}
	if v == nil {
		return nil
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		c.L.Error(req.Context(), "error unmarshalling body", logArguments...)
		return err
	}
	return nil
}
