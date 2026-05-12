package process

import (
	"testing"
	"strings"
)

func TestExec_Env(t *testing.T) {
	stdout, _, err := Exec("test env", "env")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if !strings.Contains(stdout, "GIT_TERMINAL_PROMPT=0") {
		t.Errorf("GIT_TERMINAL_PROMPT=0 not found in env output: %s", stdout)
	}
}
