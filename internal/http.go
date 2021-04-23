package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/shortcut/go-vipps/logging"
)

type HTTPError struct {
	Body   []byte
	Status int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("request failed with status: %d", e.Status)
}

type APIClient struct {
	L logging.Logger
	C *http.Client
}

func (c *APIClient) NewRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Request, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (c *APIClient) Do(req *http.Request, v interface{}) error {
	now := time.Now()
	resp, err := c.C.Do(req)
	if err != nil {
		c.L.Error(req.Context(), "error executing Vipps HTTP request", logging.NewArg("method", req.Method), logging.NewArg("url", req.URL), logging.NewArg("durationMS", time.Since(now).Milliseconds()))
		return err
	}
	c.L.Info(req.Context(), "executed Vipps HTTP request", logging.NewArg("method", req.Method), logging.NewArg("url", req.URL), logging.NewArg("durationMS", time.Since(now).Milliseconds()))
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
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
		return err
	}
	return nil
}
