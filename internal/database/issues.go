package database

import (
	"context"

	"gorm.io/gorm"
)

// IssuesStore es el almacén de datos para issues utilizando GORM.
type IssuesStore struct {
	db *gorm.DB
}

func newIssuesStore(db *gorm.DB) *IssuesStore {
	return &IssuesStore{db: db}
}

// Create crea un nuevo issue en la base de datos de forma transaccional.
// Nota: Para mantener la compatibilidad con el código legacy, seguimos usando
// database.NewIssue internamente por ahora si es necesario, o implementamos
// la versión pura aquí.
func (s *IssuesStore) Create(ctx context.Context, repo *Repository, issue *Issue, labelIDs []int64, attachments []string) error {
	// Por ahora seguimos usando el motor  a través de database.NewIssue
	// pero quitando los efectos secundarios.
	// TODO: Migrar la lógica de database.newIssue (minúscula) a GORM aquí.
	return NewIssue(repo, issue, labelIDs, attachments)
}

// database.go necesita conocer este nuevo store.
