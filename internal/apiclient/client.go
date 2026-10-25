// Package apiclient es el cliente HTTP de la API v2 de Ggit.
//
// Comparte los DTOs con el servidor (internal/route/api/v2/types), de modo que un
// cambio de contrato rompe la compilación del CLI en lugar de fallar en runtime.
package apiclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const apiPrefix = "/api/v2"

// Client es un cliente mínimo de la API v2. Es seguro para uso concurrente.
type Client struct {
	baseURL   string
	token     string
	basicAuth string
	http      *http.Client
}

// New crea un cliente apuntando a baseURL (con o sin esquema) usando token.
func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(NormalizeHost(baseURL), "/"),
		token:   token,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// WithBasicAuth devuelve una copia del cliente que autentica con usuario y
// contraseña. Se usa para crear el primer token, cuando aún no hay ninguno.
func (c *Client) WithBasicAuth(username, password string) *Client {
	clone := *c
	clone.basicAuth = base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return &clone
}

// Error representa una respuesta HTTP fuera del rango 2xx.
type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("la API respondió %d %s", e.Status, http.StatusText(e.Status))
	}
	return fmt.Sprintf("la API respondió %d: %s", e.Status, e.Message)
}

// Get ejecuta un GET sobre path (relativo a /api/v2) y decodifica la respuesta en out.
func (c *Client) Get(ctx context.Context, path string, query url.Values, out any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, out)
}

// Post ejecuta un POST sobre path (relativo a /api/v2) con body en JSON.
func (c *Client) Post(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPost, path, nil, body, out)
}

// Put ejecuta un PUT sobre path (relativo a /api/v2) con body en JSON.
func (c *Client) Put(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPut, path, nil, body, out)
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	endpoint := c.baseURL + apiPrefix + path
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var payload io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("serializar petición: %w", err)
		}
		payload = bytes.NewReader(buf)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, payload)
	if err != nil {
		return fmt.Errorf("construir petición: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// Ggit espera el esquema "token", no "Bearer" (ver internal/context/auth.go).
	switch {
	case c.token != "":
		req.Header.Set("Authorization", "token "+c.token)
	case c.basicAuth != "":
		req.Header.Set("Authorization", "Basic "+c.basicAuth)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("no se pudo contactar %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return newAPIError(resp)
	}
	if out == nil || resp.StatusCode == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("respuesta inválida de %s: %w", endpoint, err)
	}
	return nil
}

func newAPIError(resp *http.Response) error {
	// Los errores de la API llegan como {"message": "..."}; el resto se copia como texto.
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))

	var payload struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil || payload.Message == "" {
		msg := strings.TrimSpace(string(raw))
		if len(msg) > 200 {
			msg = msg[:200] + "..."
		}
		payload.Message = msg
	}
	return &Error{Status: resp.StatusCode, Message: payload.Message}
}

// NormalizeHost recorta espacios y añade https:// cuando falta esquema.
func NormalizeHost(host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}
	if !strings.Contains(host, "://") {
		host = "https://" + host
	}
	return host
}
