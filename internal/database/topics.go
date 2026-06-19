package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/nef734/ggit/internal/domain"
	"gorm.io/gorm"
)

// Topic represents a repository topic.
type Topic struct {
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	NumRepos    int
	CreatedUnix int64
	UpdatedUnix int64
}

func (t *Topic) ToDomain() *domain.Topic {
	if t == nil {
		return nil
	}
	return &domain.Topic{
		ID:       t.ID,
		Name:     t.Name,
		NumRepos: t.NumRepos,
	}
}

// BeforeCreate implements the GORM create hook.
func (t *Topic) BeforeCreate(tx *gorm.DB) error {
	if t.CreatedUnix == 0 {
		t.CreatedUnix = tx.NowFunc().Unix()
		t.UpdatedUnix = t.CreatedUnix
	}
	return nil
}

// BeforeUpdate implements the GORM update hook.
func (t *Topic) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedUnix = tx.NowFunc().Unix()
	return nil
}

// RepoTopic represents a repository and its topic.
type RepoTopic struct {
	TopicID int64 `gorm:"uniqueIndex:repo_topic_unique;not null"`
	RepoID  int64 `gorm:"uniqueIndex:repo_topic_unique;not null"`
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
		return nil, fmt.Errorf("get topics by repo ID"+": %v", err)
	}
	return topics, nil
}

// AddTopics adds topics to a repository. It creates new topics if they don't exist.
func (s *RepositoriesStore) AddTopics(ctx context.Context, repoID int64, topics []string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove existing relations
		err := tx.Where("repo_id = ?", repoID).Delete(&RepoTopic{}).Error
		if err != nil {
			return fmt.Errorf("delete existing repo topics"+": %v", err)
		}

		for _, name := range topics {
			name = strings.ToLower(strings.TrimSpace(name))
			if name == "" {
				continue
			}

			var topic Topic
			err = tx.Where("name = ?", name).FirstOrCreate(&topic, Topic{Name: name}).Error
			if err != nil {
				return fmt.Errorf("first or create topic %q: %v", name, err)
			}

			err = tx.Create(&RepoTopic{
				RepoID:  repoID,
				TopicID: topic.ID,
			}).Error
			if err != nil {
				return fmt.Errorf("create repo topic relation for %q: %v", name, err)
			}
		}
		return nil
	})
}

// ListPopularTopics returns the most used topics.
func (s *RepositoriesStore) ListPopularTopics(ctx context.Context, limit int) ([]*domain.Topic, error) {
	var topics []*Topic
	err := s.db.WithContext(ctx).
		Table("topic").
		Select("topic.*, COUNT(repo_topic.repo_id) as usage_count").
		Joins("LEFT JOIN repo_topic ON repo_topic.topic_id = topic.id").
		Group("topic.id").
		Order("usage_count DESC").
		Limit(limit).
		Find(&topics).
		Error
	if err != nil {
		return nil, fmt.Errorf("list popular topics"+": %v", err)
	}

	dtops := make([]*domain.Topic, len(topics))
	for i := range topics {
		dtops[i] = topics[i].ToDomain()
	}
	return dtops, nil
}
