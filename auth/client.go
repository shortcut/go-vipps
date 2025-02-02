// Package auth provides a HTTP client suitable for consuming Vipps APIs.
//
// Vipps APIs require that clients authorize with a client id, a client secret,
// and an API subscription key.
package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/shortcut/go-vipps"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	tokenEndpoint = "/accessToken/get"
)

type customTransport struct {
	config             clientcredentials.Config
	apiSubscriptionKey string
	rt                 http.RoundTripper
}

// NewClient returns a http.Client with a custom Transport that adds required
// headers for authorizing Vipps API clients. JWT tokens are automatically
// fetched and renewed upon expiry.
func NewClient(environment vipps.Environment, credentials vipps.Credentials) *http.Client {
	var baseURL string
	if environment == vipps.EnvironmentTesting {
		baseURL = vipps.BaseURLTesting
	} else {
		baseURL = vipps.BaseURL
	}

	tr := &customTransport{
		config: clientcredentials.Config{
			ClientID:     credentials.ClientID,
			ClientSecret: credentials.ClientSecret,
			TokenURL:     baseURL + tokenEndpoint,
		},
		apiSubscriptionKey: credentials.APISubscriptionKey,
		rt:                 http.DefaultTransport,
	}

	tokenClient := &http.Client{
		Transport: tr,
		Timeout:   20 * time.Second,
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenClient)
	return tr.config.Client(ctx)
}

// RoundTrip satisfies interface http.RoundTripper
func (ct *customTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.URL.Path == tokenEndpoint {
		request.Header.Add("client_id", ct.config.ClientID)
		request.Header.Add("client_secret", ct.config.ClientSecret)
	}
	request.Header.Add("Ocp-Apim-Subscription-Key", ct.apiSubscriptionKey)

	return ct.rt.RoundTrip(request)
}
