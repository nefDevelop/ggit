package database

import (
	"context"
	"strings"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

// Topic represents a keyword assigned to a repository.
type Topic struct {
	ID         int64  `gorm:"primaryKey"`
	Name       string `xorm:"UNIQUE" gorm:"unique;not null"`
	UsageCount int    `xorm:"-" gorm:"-"`
}

// RepoTopic represents the relation between a repository and a topic.
type RepoTopic struct {
	ID      int64 `gorm:"primaryKey"`
	RepoID  int64 `xorm:"UNIQUE(s) INDEX" gorm:"uniqueIndex:repo_topic_unique;index"`
	TopicID int64 `xorm:"UNIQUE(s) INDEX" gorm:"uniqueIndex:repo_topic_unique;index"`
}

// GetTopicsByRepoID returns all topics associated with the given repository.
func (s *RepositoriesStore) GetTopicsByRepoID(ctx context.Context, repoID int64) ([]string, error) {
	var topics []string
	err := s.db.WithContext(ctx).
		Table("topic").
		Select("topic.name").
		Joins("JOIN repo_topic ON repo_topic.topic_id = topic.id").
		Where("repo_topic.repo_id = ?", repoID).
		Find(&topics).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "get topics by repo ID")
	}
	return topics, nil
}

// AddTopics adds topics to a repository. It creates new topics if they don't exist.
func (s *RepositoriesStore) AddTopics(ctx context.Context, repoID int64, topics []string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove existing relations
		err := tx.Where("repo_id = ?", repoID).Delete(&RepoTopic{}).Error
		if err != nil {
			return errors.Wrap(err, "delete existing repo topics")
		}

		for _, name := range topics {
			name = strings.ToLower(strings.TrimSpace(name))
			if name == "" {
				continue
			}

			var topic Topic
			err = tx.Where("name = ?", name).FirstOrCreate(&topic, Topic{Name: name}).Error
			if err != nil {
				return errors.Wrapf(err, "first or create topic %q", name)
			}

			err = tx.Create(&RepoTopic{
				RepoID:  repoID,
				TopicID: topic.ID,
			}).Error
			if err != nil {
				return errors.Wrapf(err, "create repo topic relation for %q", name)
			}
		}
		return nil
	})
}

// ListPopularTopics returns the most used topics.
func (s *RepositoriesStore) ListPopularTopics(ctx context.Context, limit int) ([]*Topic, error) {
	var topics []*Topic
	err := s.db.WithContext(ctx).
		Table("topic").
		Select("topic.*, COUNT(repo_topic.id) as usage_count").
		Joins("LEFT JOIN repo_topic ON repo_topic.topic_id = topic.id").
		Group("topic.id").
		Order("usage_count DESC").
		Limit(limit).
		Find(&topics).
		Error
	if err != nil {
		return nil, errors.Wrap(err, "list popular topics")
	}
	return topics, nil
}
