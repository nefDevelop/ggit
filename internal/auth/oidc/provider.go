package oidc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nef734/ggit/internal/auth"
	"golang.org/x/oauth2"
)

func (p *Provider) withContext(ctx context.Context) context.Context {
	if p.config.SkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		return context.WithValue(ctx, oauth2.HTTPClient, client)
	}
	return ctx
}

// Provider contains configuration for an OIDC authentication provider.
type Provider struct {
	config *Config
}

// NewProvider creates a new OIDC authentication provider.
func NewProvider(cfg *Config) auth.Provider {
	return &Provider{
		config: cfg,
	}
}

// Authenticate is not supported for OIDC standard flow.
func (p *Provider) Authenticate(login, password string) (*auth.ExternalAccount, error) {
	return nil, fmt.Errorf("OIDC does not support direct password authentication")
}

func (p *Provider) Config() any {
	return p.config
}

func (p *Provider) HasTLS() bool {
	return true
}

func (p *Provider) UseTLS() bool {
	return true
}

func (p *Provider) SkipTLSVerify() bool {
	return p.config.SkipVerify
}

// OAuth2Config returns the oauth2.Config for this provider.
func (p *Provider) OAuth2Config(ctx context.Context) (*oauth2.Config, error) {
	ctx = p.withContext(ctx)
	provider, err := oidc.NewProvider(ctx, p.config.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	return &oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  p.config.RedirectURL,
		Scopes:       p.config.Scopes,
	}, nil
}

// VerifyCallback verifies the OIDC callback and returns the external account info.
func (p *Provider) VerifyCallback(ctx context.Context, code string) (*auth.ExternalAccount, error) {
	ctx = p.withContext(ctx)
	provider, err := oidc.NewProvider(ctx, p.config.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  p.config.RedirectURL,
		Scopes:       p.config.Scopes,
	}

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: p.config.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID Token: %v", err)
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %v", err)
	}

	// Map claims to ExternalAccount
	login := p.getClaim(claims, p.config.AttributeUsername, "preferred_username", "sub")
	name := p.getClaim(claims, p.config.AttributeName, "name", "nickname")
	email := p.getClaim(claims, p.config.AttributeEmail, "email")

	if login == "" {
		return nil, fmt.Errorf("could not find username claim")
	}

	return &auth.ExternalAccount{
		Login:    login,
		Name:     name,
		FullName: name,
		Email:    email,
	}, nil
}

func (p *Provider) getClaim(claims map[string]interface{}, preferred string, fallbacks ...string) string {
	if preferred != "" {
		if v, ok := claims[preferred].(string); ok {
			return v
		}
	}
	for _, f := range fallbacks {
		if v, ok := claims[f].(string); ok {
			return v
		}
	}
	return ""
}
