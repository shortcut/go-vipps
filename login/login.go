package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// IssuerURL is a type-alias for a string.
type IssuerURL string

// URLs for access-management.
const (
	IssuerURLTesting    IssuerURL = "https://apitest.vipps.no/access-management-1.0/access/"
	IssuerURLProduction IssuerURL = "https://api.vipps.no/access-management-1.0/access/"
)

// List of possible scopes supported by Vipps Login
const (
	// List containing the users addresses. Will always contain home, but can
	// also include work and other.
	ScopeAdress = "address"
	// User birth date (BankID verified)
	ScopeBirthDate = "birthDate"
	// User email (verified), the flag "email_verified : true" in the response
	// can be used by merchant to confirm for each request that the email
	// actually is verified
	ScopeEmail = "email"
	// User first, middle and given name (verified with National Population
	// Register)
	ScopeName = "name"
	// 	Verified phone number (verfied - the number used with Vipps)
	ScopePhoneNumber = "phoneNumber"
	// Norwegian national identity number (verified with BankID)
	ScopeNNIN = "nnin"
	// User bank account numbers
	ScopeAccountNumbers = "accountNumbers"
	// Signals that version 2 of the API should be used
	ScopeAPIV2 = "api_version_2"
)

// Provider is a convenience wrapper around oidc.Provider tailored to the Vipps
// Login API
type Provider struct {
	provider    *oidc.Provider
	oauthConfig oauth2.Config
	verifier    *oidc.IDTokenVerifier
}

// ProviderConfig represents a configuration for a Provider
type ProviderConfig struct {
	ClientID     string
	ClientSecret string
	IssuerURL    IssuerURL
	RedirectURL  string
	Scopes       []string
}

// Claims represents the claims contained in Vipps ID tokens
type Claims struct {
	Address struct {
		Country   string `json:"country"`
		Street    string `json:"street_address"`
		Type      string `json:"address_type"`
		Formatted string `json:"formatted"`
		Zip       string `json:"postal_code"`
		Region    string `json:"region"`
	} `json:"address"`
	OtherAddress []struct {
		Country   string `json:"country"`
		Street    string `json:"street_address"`
		Type      string `json:"address_type"`
		Formatted string `json:"formatted"`
		Zip       string `json:"postal_code"`
		Region    string `json:"region"`
	} `json:"other_address"`
	NIN           string `json:"nin"`
	PhoneNumber   string `json:"phone_number"`
	Name          string `json:"name"`
	BirthDate     Date   `json:"birthdate"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	UserID        string `json:"sub"`
}

// Date represents a date, without time-information.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// String turns the date into a string.
func (d *Date) String() string {
	return fmt.Sprintf("%d-%d-%d", d.Day, d.Month, d.Year)
}

// UnmarshalJSON unmarshalls the bytes and stores the date in the receiver.
func (d *Date) UnmarshalJSON(bytes []byte) error {
	var s, layout string
	layout = "2006-01-02"
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	d.Year = t.Year()
	d.Day = t.Day()
	d.Month = t.Month()
	return nil
}

// NewProvider returns a configured Vipps Login Provider.
func NewProvider(ctx context.Context, config *ProviderConfig) (*Provider, error) {
	if config.IssuerURL == "" {
		config.IssuerURL = IssuerURLTesting
	}
	provider, err := oidc.NewProvider(ctx, string(config.IssuerURL))
	if err != nil {
		return nil, err
	}

	oauthConfig := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  config.RedirectURL,
		Scopes:       append([]string{oidc.ScopeOpenID, ScopeAPIV2}, config.Scopes...),
	}

	return &Provider{
		provider:    provider,
		oauthConfig: oauthConfig,
		verifier: provider.Verifier(&oidc.Config{
			ClientID: config.ClientID,
		}),
	}, err
}

// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page that asks for
// permissions for the configured scopes explicitly
func (p *Provider) AuthCodeURL(state string) string {
	return p.oauthConfig.AuthCodeURL(state)
}

// ExchangeCodeForClaims takes an oauth2 authorization code, exchanges it for a
// token, and returns the contained ID token's claims, if any
func (p *Provider) ExchangeCodeForClaims(ctx context.Context, code string) (*Claims, error) {
	token, err := p.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err

	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("oauth2: no id_token in response")
	}

	idToken, err := p.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	claims := Claims{
		UserID: idToken.Subject,
	}

	userInfo, err := p.provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		return nil, err
	}

	if err = userInfo.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}
