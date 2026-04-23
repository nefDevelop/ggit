package services

import (
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/services/auth"
	"github.com/nef734/ggit/internal/services/issue"
	"github.com/nef734/ggit/internal/services/repo"
)

// Init inicializa todos los servicios de la aplicación.
func Init() {
	repo.Init(database.Handle)
	issue.Init(database.Handle)
	auth.Init(database.Handle)
}
