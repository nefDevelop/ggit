// Package pikoci habla con la API HTTP de PikoCI. En el modo recomendado
// (sidecar) PikoCI corre fuera del proceso de ggit y este servicio se limita a
// consultarlo; no gestiona su ciclo de vida.
package pikoci

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nef734/ggit/internal/conf"
)

// Health es la respuesta de GET /health.
type Health struct {
	Status  string `json:"status"`
	DB      string `json:"db"`
	Version string `json:"version"`
	Error   string `json:"error,omitempty"`
}

// Service consulta un servidor de PikoCI.
type Service struct {
	baseURL string
	http    *http.Client
}

// CI es el servicio global, inicializado desde services.Init.
var CI *Service

// Init crea el cliente con la URL configurada. No hace ninguna petición.
func Init() {
	CI = New(conf.PikoCI.URL)
}

// New crea un cliente apuntando a baseURL.
func New(baseURL string) *Service {
	return &Service{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

// BaseURL devuelve la URL configurada de PikoCI.
func (s *Service) BaseURL() string {
	return s.baseURL
}

// Health consulta /health. Un 503 se decodifica igual que un 200, porque sigue
// trayendo información útil (por ejemplo, la base de datos desconectada).
func (s *Service) Health(ctx context.Context) (*Health, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/health", nil)
	if err != nil {
		return nil, fmt.Errorf("construir petición: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("no se pudo contactar %s: %w", s.baseURL, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusServiceUnavailable:
	default:
		return nil, fmt.Errorf("%s/health respondió %d", s.baseURL, resp.StatusCode)
	}

	var health Health
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("respuesta inválida de %s/health: %w", s.baseURL, err)
	}
	return &health, nil
}

// Reachable informa si PikoCI responde y su base de datos está conectada.
func (s *Service) Reachable(ctx context.Context) bool {
	health, err := s.Health(ctx)
	return err == nil && health.Status == "ok"
}
