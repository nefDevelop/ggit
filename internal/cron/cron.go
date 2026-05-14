package cron

import (
	"time"

	"github.com/gogs/cron"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	log "github.com/nef734/ggit/internal/log"
)

var (
	c             = cron.New()
	taskFunctions = make(map[string]func())
)

func addTask(name, schedule string, task func()) {
	entry, err := c.AddFunc(name, schedule, task)
	if err != nil {
		log.Fatal("Cron.(%s): %v", name, err)
	}
	taskFunctions[name] = task

	// Check if task should run at start
	switch name {
	case "Repository health check":
		if conf.Cron.RepoHealthCheck.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go task()
		}
	case "Check repository statistics":
		if conf.Cron.CheckRepoStats.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go task()
		}
	case "Repository archive cleanup":
		if conf.Cron.RepoArchiveCleanup.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go task()
		}
	case "Sync repository hooks":
		if conf.Cron.SyncRepoHooks.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go task()
		}
	}
}

func NewContext() {
	if conf.Cron.RepoHealthCheck.Enabled {
		addTask("Repository health check", conf.Cron.RepoHealthCheck.Schedule, database.GitFsck)
	}
	if conf.Cron.CheckRepoStats.Enabled {
		addTask("Check repository statistics", conf.Cron.CheckRepoStats.Schedule, database.CheckRepoStats)
	}
	if conf.Cron.RepoArchiveCleanup.Enabled {
		addTask("Repository archive cleanup", conf.Cron.RepoArchiveCleanup.Schedule, database.DeleteOldRepositoryArchives)
	}
	if conf.Cron.SyncRepoHooks.Enabled {
		addTask("Sync repository hooks", conf.Cron.SyncRepoHooks.Schedule, database.SyncRepoHooks)
	}

	// Orphaned attachments cleanup (Every 24h)
	addTask("Cleanup orphaned attachments", "@every 24h", func() {
		count, err := database.DeleteOrphanedAttachments(24 * time.Hour)
		if err != nil {
			log.Error("Cron.(cleanup orphaned attachments): %v", err)
		} else if count > 0 {
			log.Trace("Cron.(cleanup orphaned attachments): deleted %d files", count)
		}
	})

	c.Start()
}

// ListTasks returns all running cron tasks.
func ListTasks() []*cron.Entry {
	return c.Entries()
}

// RunTask runs a cron task manually by its description.
func RunTask(name string) {
	if task, ok := taskFunctions[name]; ok {
		go task()
	}
}
