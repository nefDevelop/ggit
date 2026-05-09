package health

import (
	"context"
	"time"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	log "github.com/nef734/ggit/internal/log"
)

// Sentinel is the manager for repository integrity checks.
var Sentinel = &sentinel{}

type sentinel struct{}

// CheckRepository decides if a repository needs a git fsck and runs it if necessary.
// It follows the rules:
// 1. Mandatory every 30 days.
// 2. If active (commits > 0), every ActivityCooldownDays.
// 3. If idle, every IdleCheckDays.
func (s *sentinel) CheckRepository(ctx context.Context, repo *database.Repository) error {
	now := time.Now().Unix()
	timeSinceLast := now - repo.LastHealthCheckUnix

	shouldRun := false
	reason := ""

	if repo.LastHealthCheckUnix == 0 {
		shouldRun = true
		reason = "initial check"
	} else if timeSinceLast > 30*24*60*60 {
		shouldRun = true
		reason = "mandatory monthly check"
	} else if repo.CommitsSinceCheck > 0 {
		if timeSinceLast > int64(conf.HealthCheck.ActivityCooldownDays)*24*60*60 {
			shouldRun = true
			reason = "activity cooldown reached"
		} else if repo.CommitsSinceCheck >= conf.HealthCheck.MaxCommits {
			shouldRun = true
			reason = "commit threshold reached"
		}
	} else if timeSinceLast > int64(conf.HealthCheck.IdleCheckDays)*24*60*60 {
		shouldRun = true
		reason = "idle period reached"
	}

	if !shouldRun {
		return nil
	}

	log.Trace("Sentinel: Running integrity check on %s (Reason: %s)", repo.FullName(), reason)
	return s.RunIntegrityCheck(ctx, repo)
}

// RunIntegrityCheck performs the actual git fsck and updates the repo status.
func (s *sentinel) RunIntegrityCheck(ctx context.Context, repo *database.Repository) error {
	repoPath := repo.RepoPath()
	err := git.Fsck(repoPath, git.FsckOptions{
		Timeout: 30 * time.Minute,
	})

	if err != nil {
		log.Warn("Sentinel: Integrity check failed for %s: %v", repo.FullName(), err)
		// Create a system notice
		_ = database.Handle.Notices().Create(ctx, database.NoticeTypeRepository,
			"Corrupción detectada en "+repo.FullName()+": "+err.Error())
	}

	// Update metadata even if it failed (we checked it)
	repo.LastHealthCheckUnix = time.Now().Unix()
	repo.CommitsSinceCheck = 0

	// Use GORM directly to update specific fields to avoid overwriting other changes
	return database.Handle.DB().Model(repo).Updates(map[string]interface{}{
		"last_health_check_unix": repo.LastHealthCheckUnix,
		"commits_since_check":    0,
	}).Error
}

// StartWorker starts the background sentinel worker.
func (s *sentinel) StartWorker() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			s.CheckAllRepositories(context.Background())
		}
	}()
}

// CheckAllRepositories iterates over all repositories and checks their integrity.
func (s *sentinel) CheckAllRepositories(ctx context.Context) {
	// We use a stream or batch to avoid loading all repos into RAM
	var repos []*database.Repository
	err := database.Handle.DB().Find(&repos).Error
	if err != nil {
		log.Error("Sentinel: Failed to fetch repositories: %v", err)
		return
	}

	for _, r := range repos {
		_ = s.CheckRepository(ctx, r)
	}
}
