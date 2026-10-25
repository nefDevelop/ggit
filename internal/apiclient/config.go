package apiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// EnvHost y EnvToken sobrescriben lo que haya en el archivo de configuración.
const (
	EnvHost  = "GGIT_HOST"
	EnvToken = "GGIT_TOKEN"
)

// Config es la configuración persistida del CLI.
type Config struct {
	Host  string `json:"host"`
	Token string `json:"token"`
}

// LoadConfig mezcla el archivo de configuración con las variables de entorno,
// que tienen prioridad. Un archivo inexistente no es un error.
func LoadConfig() (Config, error) {
	var cfg Config

	path, err := configPath()
	if err != nil {
		return cfg, err
	}

	raw, err := os.ReadFile(path)
	switch {
	case err == nil:
		if err := json.Unmarshal(raw, &cfg); err != nil {
			return cfg, fmt.Errorf("configuración inválida en %s: %w", path, err)
		}
	case errors.Is(err, fs.ErrNotExist):
		// Sin archivo: solo entorno y flags.
	default:
		return cfg, fmt.Errorf("leer %s: %w", path, err)
	}

	if v := os.Getenv(EnvHost); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv(EnvToken); v != "" {
		cfg.Token = v
	}
	return cfg, nil
}

// ConfigPath devuelve la ruta del archivo de configuración del CLI.
func ConfigPath() (string, error) {
	return configPath()
}

// SaveConfig escribe la configuración con permisos 0600, porque contiene el token.
func SaveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("crear %s: %w", filepath.Dir(path), err)
	}

	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serializar configuración: %w", err)
	}
	raw = append(raw, '\n')

	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return fmt.Errorf("escribir %s: %w", path, err)
	}
	// WriteFile respeta el modo de un archivo existente, así que lo forzamos.
	if err := os.Chmod(path, 0o600); err != nil {
		return fmt.Errorf("ajustar permisos de %s: %w", path, err)
	}
	return nil
}

// RemoveConfig borra el archivo de configuración local. No falla si no existe.
func RemoveConfig() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("borrar %s: %w", path, err)
	}
	return nil
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("localizar el directorio de configuración: %w", err)
	}
	return filepath.Join(dir, "ggit-cli", "config.json"), nil
}
