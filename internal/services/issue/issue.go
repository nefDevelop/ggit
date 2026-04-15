package issue

import (
	"context"
	"fmt"

	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/database"
	apiv1types "gogs.io/gogs/internal/route/api/v1/types"
)

// IssueService define las operaciones de negocio para issues y pull requests.
type IssueService interface {
	CreateIssue(ctx context.Context, repo *database.Repository, issue *database.Issue, opts CreateIssueOptions) error
	CreatePullRequest(ctx context.Context, repo *database.Repository, pull *database.Issue, pr *database.PullRequest, opts CreatePullRequestOptions) error
}

type issueService struct {
	db *database.DB
}

// Issues es el singleton global para el servicio de issues.
var Issues IssueService

// Init inicializa el servicio de issues.
func Init(db *database.DB) {
	Issues = &issueService{
		db: db,
	}
}

func (s *issueService) CreateIssue(ctx context.Context, repo *database.Repository, issue *database.Issue, opts CreateIssueOptions) error {
	// 1. Crear el issue en la base de datos (puro)
	err := s.db.Issues().Create(ctx, repo, issue, opts.LabelIDs, opts.Attachments)
	if err != nil {
		return err
	}

	// 2. Orquestación de efectos secundarios
	if err = database.NotifyWatchers(&database.Action{
		ActUserID:    issue.Poster.ID,
		ActUserName:  issue.Poster.Name,
		OpType:       database.ActionCreateIssue,
		Content:      fmt.Sprintf("%d|%s", issue.Index, issue.Title),
		RepoID:       repo.ID,
		RepoUserName: repo.Owner.Name,
		RepoName:     repo.Name,
		IsPrivate:    repo.IsPrivate,
	}); err != nil {
		log.Error("NotifyWatchers: %v", err)
	}
	if err = issue.MailParticipants(); err != nil {
		log.Error("MailParticipants: %v", err)
	}

	if err = database.PrepareWebhooks(repo, database.HookEventTypeIssues, &apiv1types.WebhookIssuesPayload{
		Action:     apiv1types.WebhookIssueOpened,
		Index:      issue.Index,
		Issue:      issue.APIFormat(),
		Repository: repo.APIFormatLegacy(nil),
		Sender:     issue.Poster.APIFormat(),
	}); err != nil {
		log.Error("PrepareWebhooks: %v", err)
	}

	return nil
}

func (s *issueService) CreatePullRequest(ctx context.Context, repo *database.Repository, pull *database.Issue, pr *database.PullRequest, opts CreatePullRequestOptions) error {
	// 1. Crear el Pull Request en la base de datos (puro)
	err := database.NewPullRequest(repo, pull, opts.LabelIDs, opts.Attachments, pr, opts.Patch)
	if err != nil {
		return err
	}

	// 2. Orquestación de efectos secundarios
	if err = database.NotifyWatchers(&database.Action{
		ActUserID:    pull.Poster.ID,
		ActUserName:  pull.Poster.Name,
		OpType:       database.ActionCreatePullRequest,
		Content:      fmt.Sprintf("%d|%s", pull.Index, pull.Title),
		RepoID:       repo.ID,
		RepoUserName: repo.Owner.Name,
		RepoName:     repo.Name,
		IsPrivate:    repo.IsPrivate,
	}); err != nil {
		log.Error("NotifyWatchers: %v", err)
	}
	if err = pull.MailParticipants(); err != nil {
		log.Error("MailParticipants: %v", err)
	}

	if err = database.PrepareWebhooks(repo, database.HookEventTypePullRequest, &apiv1types.WebhookPullRequestPayload{
		Action:      apiv1types.WebhookIssueOpened,
		Index:       pull.Index,
		PullRequest: pr.APIFormat(),
		Repository:  repo.APIFormatLegacy(nil),
		Sender:      pull.Poster.APIFormat(),
	}); err != nil {
		log.Error("PrepareWebhooks: %v", err)
	}

	// 3. Sincronizar con el repositorio base (esto orquesta Git)
	if err := pr.PushToBaseRepo(); err != nil {
		return err
	}

	return nil
}
