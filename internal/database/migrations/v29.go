package migrations

import (
	"gorm.io/gorm"
)

func addActivityTimestampsAndSentinelToRepository(db *gorm.DB) error {
	type repository struct {
		LastHealthCheckUnix int64 `gorm:"not null;default:0"`
		CommitsSinceCheck   int   `gorm:"not null;default:0"`
		LastCommitUnix      int64 `gorm:"not null;default:0"`
		LastIssueUnix       int64 `gorm:"not null;default:0"`
		LastPullUnix        int64 `gorm:"not null;default:0"`
		LastBranchUnix      int64 `gorm:"not null;default:0"`
		IsMaintenance       bool  `gorm:"not null;default:FALSE"`
	}

	return db.Migrator().AutoMigrate(&repository{})
}
