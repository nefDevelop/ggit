package cron

import (
	"time"

	log "unknwon.dev/clog/v2"

	"github.com/gogs/cron"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
)

var c = cron.New()

func NewContext() {
	var (
		entry *cron.Entry
		err   error
	)
	if conf.Cron.RepoHealthCheck.Enabled {
		entry, err = c.AddFunc("Repository health check", conf.Cron.RepoHealthCheck.Schedule, database.GitFsck)
		if err != nil {
			log.Fatal("Cron.(repository health check): %v", err)
		}
		if conf.Cron.RepoHealthCheck.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go database.GitFsck()
		}
	}
	if conf.Cron.CheckRepoStats.Enabled {
		entry, err = c.AddFunc("Check repository statistics", conf.Cron.CheckRepoStats.Schedule, database.CheckRepoStats)
		if err != nil {
			log.Fatal("Cron.(check repository statistics): %v", err)
		}
		if conf.Cron.CheckRepoStats.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go database.CheckRepoStats()
		}
	}
	if conf.Cron.RepoArchiveCleanup.Enabled {
		entry, err = c.AddFunc("Repository archive cleanup", conf.Cron.RepoArchiveCleanup.Schedule, database.DeleteOldRepositoryArchives)
		if err != nil {
			log.Fatal("Cron.(repository archive cleanup): %v", err)
		}
		if conf.Cron.RepoArchiveCleanup.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go database.DeleteOldRepositoryArchives()
		}
	}
	if conf.Cron.SyncRepoHooks.Enabled {
		entry, err = c.AddFunc("Sync repository hooks", conf.Cron.SyncRepoHooks.Schedule, database.SyncRepoHooks)
		if err != nil {
			log.Fatal("Cron.(sync repository hooks): %v", err)
		}
		if conf.Cron.SyncRepoHooks.RunAtStart {
			entry.Prev = time.Now()
			entry.ExecTimes++
			go database.SyncRepoHooks()
		}
	}
	c.Start()
}

// ListTasks returns all running cron tasks.
func ListTasks() []*cron.Entry {
	return c.Entries()
}
