package context

import (
	"github.com/nef734/ggit/internal/database"
)

type APIOrganization struct {
	Organization *database.User
	Team         *database.Team
}
