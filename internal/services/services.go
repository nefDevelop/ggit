package services

import (
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/services/auth"
	"gogs.io/gogs/internal/services/issue"
	"gogs.io/gogs/internal/services/repo"
)

// Init inicializa todos los servicios de la aplicación.
func Init() {
	repo.Init(database.Handle)
	issue.Init(database.Handle)
	auth.Init(database.Handle)
}
