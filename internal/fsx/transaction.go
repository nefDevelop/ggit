package fsx

import (
	"fmt"
	"os"
	"path/filepath"
)

type opType int

const (
	opCreate opType = iota
	opDelete
)

type operation struct {
	typ       opType
	tempPath  string
	finalPath string
}

// Transaction coordinates filesystem operations that can be committed or rolled back.
type Transaction struct {
	ops []operation
}

// NewTransaction creates a new filesystem transaction.
func NewTransaction() *Transaction {
	return &Transaction{}
}

// StageCreate stages a creation operation. Changes should be made in tempPath.
// On commit, tempPath is moved to finalPath.
// On rollback, tempPath is deleted.
func (t *Transaction) StageCreate(tempPath, finalPath string) {
	t.ops = append(t.ops, operation{
		typ:       opCreate,
		tempPath:  tempPath,
		finalPath: finalPath,
	})
}

// StageDelete stages a deletion operation.
// It immediately moves finalPath to tempPath.
// On commit, tempPath is permanently deleted.
// On rollback, tempPath is moved back to finalPath.
func (t *Transaction) StageDelete(finalPath, tempPath string) error {
	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		return nil // Nothing to delete
	}

	if err := os.MkdirAll(filepath.Dir(tempPath), 0755); err != nil {
		return fmt.Errorf("mkdir temp path parent"+": %w", err)
	}

	if err := os.Rename(finalPath, tempPath); err != nil {
		return fmt.Errorf("move to temp for deletion"+": %w", err)
	}

	t.ops = append(t.ops, operation{
		typ:       opDelete,
		tempPath:  tempPath,
		finalPath: finalPath,
	})
	return nil
}

// Commit executes the final phase of all staged operations.
func (t *Transaction) Commit() error {
	for _, op := range t.ops {
		switch op.typ {
		case opCreate:
			if _, err := os.Stat(op.tempPath); os.IsNotExist(err) {
				continue
			}
			if err := os.MkdirAll(filepath.Dir(op.finalPath), 0755); err != nil {
				return fmt.Errorf("mkdir final path parent"+": %w", err)
			}
			if err := os.Rename(op.tempPath, op.finalPath); err != nil {
				return fmt.Errorf("rename temp to final"+": %w", err)
			}
		case opDelete:
			if err := os.RemoveAll(op.tempPath); err != nil {
				return fmt.Errorf("permanent delete"+": %w", err)
			}
		}
	}
	t.ops = nil
	return nil
}

// Rollback reverts all staged operations.
func (t *Transaction) Rollback() error {
	for _, op := range t.ops {
		switch op.typ {
		case opCreate:
			_ = os.RemoveAll(op.tempPath)
		case opDelete:
			if _, err := os.Stat(op.tempPath); err == nil {
				_ = os.Rename(op.tempPath, op.finalPath)
			}
		}
	}
	t.ops = nil
	return nil
}
