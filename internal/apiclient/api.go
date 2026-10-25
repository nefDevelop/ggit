package apiclient

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	apiv2 "github.com/nef734/ggit/internal/route/api/v2/types"
)

// AuthStatus valida las credenciales y devuelve el usuario autenticado.
func (c *Client) AuthStatus(ctx context.Context) (*apiv2.UserResponse, error) {
	var res struct {
		Status string              `json:"status"`
		User   *apiv2.UserResponse `json:"user"`
	}
	if err := c.Get(ctx, "/user/check", nil, &res); err != nil {
		return nil, err
	}
	if res.User == nil {
		return nil, fmt.Errorf("la API no devolvió el usuario autenticado")
	}
	return res.User, nil
}

// ListUserRepos devuelve los repositorios propios y las colaboraciones.
func (c *Client) ListUserRepos(ctx context.Context) ([]*apiv2.RepositoryResponse, error) {
	var repos []*apiv2.RepositoryResponse
	if err := c.Get(ctx, "/user/repos", nil, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

// SearchRepos busca repositorios visibles para el usuario autenticado.
func (c *Client) SearchRepos(ctx context.Context, keyword string, page int) ([]*apiv2.RepositoryResponse, int64, error) {
	q := url.Values{}
	if keyword != "" {
		q.Set("q", keyword)
	}
	if page > 1 {
		q.Set("page", strconv.Itoa(page))
	}

	var res struct {
		Repos []*apiv2.RepositoryResponse `json:"repos"`
		Total int64                       `json:"total"`
	}
	if err := c.Get(ctx, "/explore/repos", q, &res); err != nil {
		return nil, 0, err
	}
	return res.Repos, res.Total, nil
}

// GetRepo resuelve owner/name con el buscador: la API v2 no publica metadatos
// de un repositorio suelto en /{owner}/{repo}.
func (c *Client) GetRepo(ctx context.Context, fullName string) (*apiv2.RepositoryResponse, error) {
	_, name, err := SplitRepoName(fullName)
	if err != nil {
		return nil, err
	}

	repos, _, err := c.SearchRepos(ctx, name, 1)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		if strings.EqualFold(repo.FullName, fullName) {
			return repo, nil
		}
	}
	return nil, fmt.Errorf("repositorio %q no encontrado o sin acceso", fullName)
}

// ListIssuesOptions filtra el listado de incidencias de un repositorio.
type ListIssuesOptions struct {
	State       string
	Page        int
	Labels      string
	MilestoneID int64
	AssigneeID  int64
	PosterID    int64
	Sort        string
}

// ListIssues devuelve las incidencias de un repositorio.
func (c *Client) ListIssues(ctx context.Context, owner, name string, opts ListIssuesOptions) ([]*apiv2.IssueResponse, error) {
	q := url.Values{}
	if opts.State != "" {
		q.Set("state", opts.State)
	}
	if opts.Page > 1 {
		q.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Labels != "" {
		q.Set("labels", opts.Labels)
	}
	if opts.MilestoneID > 0 {
		q.Set("milestone_id", strconv.FormatInt(opts.MilestoneID, 10))
	}
	if opts.AssigneeID > 0 {
		q.Set("assignee_id", strconv.FormatInt(opts.AssigneeID, 10))
	}
	if opts.PosterID > 0 {
		q.Set("poster_id", strconv.FormatInt(opts.PosterID, 10))
	}
	if opts.Sort != "" {
		q.Set("sort", opts.Sort)
	}

	var issues []*apiv2.IssueResponse
	if err := c.Get(ctx, repoPath(owner, name, "/issues"), q, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

// GetIssue devuelve una incidencia por su índice dentro del repositorio.
func (c *Client) GetIssue(ctx context.Context, owner, name string, index int64) (*apiv2.IssueResponse, error) {
	var issue apiv2.IssueResponse
	path := repoPath(owner, name, "/issues/"+strconv.FormatInt(index, 10))
	if err := c.Get(ctx, path, nil, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

// SplitRepoName valida y separa un argumento "owner/name".
func SplitRepoName(fullName string) (owner, name string, err error) {
	parts := strings.Split(strings.Trim(fullName, "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("formato de repositorio inválido %q: se espera owner/name", fullName)
	}
	return parts[0], parts[1], nil
}

func repoPath(owner, name, suffix string) string {
	return "/" + url.PathEscape(owner) + "/" + url.PathEscape(name) + suffix
}
