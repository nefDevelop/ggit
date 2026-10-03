package conf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/ini.v1"
)

// TestPersistSecretKeyCreatesConfig covers the first-start path: no custom
// configuration file exists yet, so both the directory and the file have to be
// created before the key can be stored.
func TestPersistSecretKeyCreatesConfig(t *testing.T) {
	dir := t.TempDir()
	restore := swapCustomConf(filepath.Join(dir, "nested", "conf", "app.ini"))
	defer restore()

	File = ini.Empty()
	File.NameMapper = ini.SnackCase

	const key = "clave-de-prueba-con-suficiente-entropia-1234"
	if err := persistSecretKey(key); err != nil {
		t.Fatalf("persistSecretKey() = %v", err)
	}

	if _, err := os.Stat(CustomConf); err != nil {
		t.Fatalf("el archivo de configuracion no se creo: %v", err)
	}

	written, err := ini.Load(CustomConf)
	if err != nil {
		t.Fatalf("no se pudo releer el archivo: %v", err)
	}
	if got := written.Section("security").Key("SECRET_KEY").String(); got != key {
		t.Fatalf("SECRET_KEY en disco = %q, se esperaba %q", got, key)
	}
	if got := File.Section("security").Key("SECRET_KEY").String(); got != key {
		t.Fatalf("SECRET_KEY en memoria = %q, se esperaba %q", got, key)
	}
}

// TestPersistSecretKeyPreservesExistingConfig makes sure adding the key does not
// discard the rest of a deployment's custom configuration.
func TestPersistSecretKeyPreservesExistingConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.ini")
	existing := "[database]\nDB_TYPE = sqlite3\n\n[security]\nINSTALL_LOCK = true\n"
	if err := os.WriteFile(path, []byte(existing), 0o600); err != nil {
		t.Fatalf("no se pudo crear la config: %v", err)
	}

	restore := swapCustomConf(path)
	defer restore()

	File = ini.Empty()
	File.NameMapper = ini.SnackCase

	const key = "otra-clave-de-prueba-1234567890"
	if err := persistSecretKey(key); err != nil {
		t.Fatalf("persistSecretKey() = %v", err)
	}

	written, err := ini.Load(CustomConf)
	if err != nil {
		t.Fatalf("no se pudo releer el archivo: %v", err)
	}
	if got := written.Section("database").Key("DB_TYPE").String(); got != "sqlite3" {
		t.Fatalf("DB_TYPE = %q, se conservaba sqlite3", got)
	}
	if got := written.Section("security").Key("INSTALL_LOCK").String(); got != "true" {
		t.Fatalf("INSTALL_LOCK = %q, se conservaba true", got)
	}
	if got := written.Section("security").Key("SECRET_KEY").String(); got != key {
		t.Fatalf("SECRET_KEY = %q, se esperaba %q", got, key)
	}
}

// TestEnsureSecretKeyKeepsConfiguredKey is the upgrade guard: a deployment that
// already has its own key must not be rotated behind the operator's back, since
// rotating invalidates every session and every stored 2FA secret.
func TestEnsureSecretKeyKeepsConfiguredKey(t *testing.T) {
	restore := swapCustomConf(filepath.Join(t.TempDir(), "app.ini"))
	defer restore()

	Security.SecretKey = "una-clave-ya-configurada-por-el-operador"
	if err := ensureSecretKey(); err != nil {
		t.Fatalf("ensureSecretKey() = %v", err)
	}
	if Security.SecretKey != "una-clave-ya-configurada-por-el-operador" {
		t.Fatalf("la clave se sustituyo por %q, deberia conservarse", Security.SecretKey)
	}
	if _, err := os.Stat(CustomConf); err == nil {
		t.Fatal("no deberia haberse escrito la config cuando la clave ya era valida")
	}
}

// TestEnsureSecretKeyRejectsEmbeddedDefault covers the original problem: the
// key that ships inside the executable must never be the one in use.
func TestEnsureSecretKeyRejectsEmbeddedDefault(t *testing.T) {
	restore := swapCustomConf(filepath.Join(t.TempDir(), "app.ini"))
	defer restore()

	Security.SecretKey = embeddedSecretKey
	if err := ensureSecretKey(); err != nil {
		t.Fatalf("ensureSecretKey() = %v", err)
	}
	if Security.SecretKey == embeddedSecretKey {
		t.Fatal("la clave embebida en el binario seguia en uso")
	}
	if strings.TrimSpace(Security.SecretKey) == "" {
		t.Fatal("la clave quedo vacia")
	}
}

// swapCustomConf points the package level CustomConf at a temporary path and
// returns a function that restores every global the tests below touch.
//
// Security is restored too: these tests mutate Security.SecretKey, and the
// golden test in this package serialises whatever value it finds there, so a
// leaked key makes TestInit fail depending on the shuffled test order.
func swapCustomConf(path string) func() {
	previousConf := CustomConf
	previousFile := File
	previousSecurity := Security
	CustomConf = path
	return func() {
		CustomConf = previousConf
		File = previousFile
		Security = previousSecurity
	}
}
