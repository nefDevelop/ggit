// Package pikoci habla con la API HTTP de PikoCI. En el modo recomendado
// (sidecar) PikoCI corre fuera del proceso de ggit y este servicio se limita a
// consultarlo; no gestiona su ciclo de vida.
package pikoci

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
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
	baseURL  string
	team     string
	apiToken string
	http     *http.Client
	links    database.PikoCIRepoStore // nil: se resuelve con database.Handle
}

// CI es el servicio global, inicializado desde services.Init.
var CI *Service

// Init crea el cliente con la URL y las credenciales configuradas. No hace
// ninguna petición.
func Init() {
	CI = NewWithAuth(conf.PikoCI.URL, conf.PikoCI.Team, conf.PikoCI.APIToken)
}

// New crea un cliente apuntando a baseURL, sin credenciales de administración.
func New(baseURL string) *Service {
	return NewWithAuth(baseURL, "main", "")
}

// NewWithAuth crea un cliente con el equipo y el token de la API de PikoCI.
func NewWithAuth(baseURL, team, apiToken string) *Service {
	if team == "" {
		team = "main"
	}
	return &Service{
		baseURL:  strings.TrimRight(baseURL, "/"),
		team:     team,
		apiToken: apiToken,
		http:     &http.Client{Timeout: 30 * time.Second},
	}
}

// repoLinks devuelve el store de vínculos, inyectable en los tests.
func (s *Service) repoLinks() database.PikoCIRepoStore {
	if s.links != nil {
		return s.links
	}
	return database.Handle.PikoCIRepos()
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

// TriggerWebhook lanza el check del resource asociado al token. PikoCI ignora
// el cuerpo y no firma: el token es el único secreto.
func (s *Service) TriggerWebhook(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}

	endpoint := s.baseURL + "/webhooks/" + url.PathEscape(token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("construir petición: %w", err)
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return fmt.Errorf("no se pudo avisar a %s: %w", s.baseURL, err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("PikoCI respondió %d al trigger del webhook", resp.StatusCode)
	}
	return nil
}

// NotifyPush avisa a PikoCI de que el repositorio tiene refs nuevas. Es
// best-effort: si el repositorio no tiene pipeline asociado no hace nada.
func (s *Service) NotifyPush(ctx context.Context, repoID int64) error {
	link, err := s.repoLinks().GetByRepoID(ctx, repoID)
	if err != nil {
		if database.IsErrPikoCIRepoNotExist(err) {
			return nil
		}
		return fmt.Errorf("buscar vínculo de PikoCI: %w", err)
	}
	return s.TriggerWebhook(ctx, link.WebhookToken)
}
