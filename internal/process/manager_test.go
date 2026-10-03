package process

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestExecDirTimeoutKillsCommand is the regression guard for the bug that made
// Kill a no-op. The guard used to require cmd.ProcessState to be non-nil, but
// that field is only populated once Wait returns, so a command that was still
// running never satisfied it. The kill was skipped and the caller blocked on the
// command forever, which is what hung the repository creation test for the full
// ten minute go test timeout and left the admin "kill process" button inert.
func TestExecDirTimeoutKillsCommand(t *testing.T) {
	const timeout = time.Second

	done := make(chan error, 1)
	go func() {
		_, _, err := ExecTimeout(timeout, "test-sleep", "sh", "-c", "sleep 300")
		done <- err
	}()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, ErrExecTimeout, "un comando colgado debe acabar en error de timeout")
	case <-time.After(20 * time.Second):
		t.Fatal("ExecTimeout no volvio: Kill no esta terminando el proceso")
	}
}

// TestExecDirTimeoutReturnsPromptly pins down the timing, not just the error: the
// call has to return shortly after the timeout, which only happens if the
// command is actually signalled.
func TestExecDirTimeoutReturnsPromptly(t *testing.T) {
	start := time.Now()
	_, _, err := ExecTimeout(500*time.Millisecond, "test-sleep-prompt", "sh", "-c", "sleep 300")
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, ErrExecTimeout)
	assert.Less(t, elapsed, 15*time.Second,
		"el comando debio terminar en cuanto se cumplio el timeout, tardo %s", elapsed)
}

// TestKillRemovesProcessFromList checks that a killed process does not stay in
// the global list, which would otherwise grow by one entry for every command that
// ever timed out.
func TestKillRemovesProcessFromList(t *testing.T) {
	start := len(Processes)

	_, _, err := ExecTimeout(300*time.Millisecond, "test-sleep-list", "sh", "-c", "sleep 30")
	assert.ErrorIs(t, err, ErrExecTimeout)

	counter.Lock()
	remaining := len(Processes)
	counter.Unlock()

	assert.Equal(t, start, remaining,
		"el proceso matado debe desaparecer de la lista global")
}

// TestKillOnUnknownPIDIsNoOp keeps the admin endpoint harmless when the process
// already finished on its own.
func TestKillOnUnknownPIDIsNoOp(t *testing.T) {
	assert.NoError(t, Kill(-1))
	assert.NoError(t, Kill(999999))
}

// TestKillTerminatesProcessGroup covers the case the single process kill missed:
// the command spawned a child of its own, and killing only the parent leaves the
// child holding the inherited stdout pipe.
func TestKillTerminatesProcessGroup(t *testing.T) {
	const timeout = 700 * time.Millisecond

	start := time.Now()
	// The shell backgrounds a sleep and then blocks. Only signalling the whole
	// group reaps the backgrounded child too.
	_, _, err := ExecTimeout(timeout, "test-group", "sh", "-c", "sleep 300 & wait")
	elapsed := time.Since(start)

	assert.ErrorIs(t, err, ErrExecTimeout)
	assert.Less(t, elapsed, 15*time.Second,
		"el grupo de procesos debio terminarse, tardo %s", elapsed)
}
