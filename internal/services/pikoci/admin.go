package pikoci

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
)

// PublishPipeline crea o actualiza un pipeline en PikoCI. La API espera el HCL
// en base64 dentro de "config".
func (s *Service) PublishPipeline(ctx context.Context, name string, hcl []byte, vars map[string]string) error {
	body := map[string]any{
		"config": base64.StdEncoding.EncodeToString(hcl),
		"vars":   vars,
	}
	path := fmt.Sprintf("/teams/%s/pipelines/%s", url.PathEscape(s.team), url.PathEscape(name))
	return s.doAdmin(ctx, http.MethodPut, path, body, nil)
}

// ResourceWebhookToken devuelve el token con el que ggit dispara el check del
// resource. Requiere rol maintain, que cubre el API_TOKEN de servicio.
func (s *Service) ResourceWebhookToken(ctx context.Context, pipeline, resource string) (string, error) {
	path := fmt.Sprintf("/teams/%s/pipelines/%s/resources/%s/webhook_token",
		url.PathEscape(s.team), url.PathEscape(pipeline), url.PathEscape(resource))

	var res struct {
		Token string `json:"token"`
	}
	if err := s.doAdmin(ctx, http.MethodPost, path, nil, &res); err != nil {
		return "", err
	}
	if res.Token == "" {
		return "", fmt.Errorf("PikoCI no devolvió el token del webhook de %q", pipeline)
	}
	return res.Token, nil
}

// DeletePipeline borra un pipeline. Un 404 se trata como éxito: el objetivo es
// que no quede nada.
func (s *Service) DeletePipeline(ctx context.Context, name string) error {
	path := fmt.Sprintf("/teams/%s/pipelines/%s", url.PathEscape(s.team), url.PathEscape(name))
	err := s.doAdmin(ctx, http.MethodDelete, path, nil, nil)
	if err != nil && strings.Contains(err.Error(), "respondió 404") {
		return nil
	}
	return err
}

func (s *Service) doAdmin(ctx context.Context, method, path string, body, out any) error {
	if s.apiToken == "" {
		return fmt.Errorf("falta [pikoci] API_TOKEN para hablar con la API de administración")
	}

	var payload io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("serializar petición: %w", err)
		}
		payload = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+path, payload)
	if err != nil {
		return fmt.Errorf("construir petición: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+s.apiToken)

	resp, err := s.http.Do(req)
	if err != nil {
		return fmt.Errorf("no se pudo contactar %s: %w", s.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		message := strings.TrimSpace(string(raw))
		var parsed struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		if json.Unmarshal(raw, &parsed) == nil {
			if parsed.Message != "" {
				message = parsed.Message
			} else if parsed.Error != "" {
				message = parsed.Error
			}
		}
		if len(message) > 200 {
			message = message[:200] + "..."
		}
		return fmt.Errorf("PikoCI respondió %d: %s", resp.StatusCode, message)
	}

	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("respuesta inválida de %s: %w", path, err)
	}
	return nil
}
