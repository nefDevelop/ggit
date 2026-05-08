package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider_OAuth2Config(t *testing.T) {
	// Mock OIDC server
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":                 "http://" + r.Host,
			"authorization_endpoint": "http://" + r.Host + "/auth",
			"token_endpoint":         "http://" + r.Host + "/token",
			"jwks_uri":               "http://" + r.Host + "/keys",
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := &Config{
		IssuerURL:    server.URL,
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "http://localhost/callback",
		Scopes:       []string{"openid", "email"},
	}
	p := &Provider{config: cfg}

	oauth2Cfg, err := p.OAuth2Config(context.Background())
	require.NoError(t, err)
	assert.Equal(t, cfg.ClientID, oauth2Cfg.ClientID)
	assert.Equal(t, cfg.ClientSecret, oauth2Cfg.ClientSecret)
	assert.Equal(t, cfg.RedirectURL, oauth2Cfg.RedirectURL)
	assert.Contains(t, oauth2Cfg.Scopes, "openid")
	assert.Contains(t, oauth2Cfg.Scopes, "email")
	assert.Equal(t, server.URL+"/auth", oauth2Cfg.Endpoint.AuthURL)
	assert.Equal(t, server.URL+"/token", oauth2Cfg.Endpoint.TokenURL)
}
