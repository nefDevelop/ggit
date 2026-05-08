package oidc

// Config contains configuration for an OIDC authentication provider.
type Config struct {
	IssuerURL    string   `ini:"issuer_url"`
	ClientID     string   `ini:"client_id"`
	ClientSecret string   `ini:"client_secret"`
	Scopes       []string `ini:"scopes" delim:","`
	RedirectURL  string   `ini:"redirect_url"`

	// Mapping fields
	AttributeUsername string `ini:"attr_username"`
	AttributeName     string `ini:"attr_name"`
	AttributeEmail    string `ini:"attr_email"`

	SkipVerify bool `ini:"skip_verify"`
}
