package apiclient

import (
	"context"
	"encoding/json"
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

// CreateAccessToken crea un token de acceso con el nombre dado y devuelve su
// valor en claro. Solo se puede leer en el momento de la creación.
func (c *Client) CreateAccessToken(ctx context.Context, name string) (string, error) {
	var res apiv2.UserAccessTokenResponse
	if err := c.Post(ctx, "/user/settings/tokens", map[string]string{"name": name}, &res); err != nil {
		return "", err
	}
	if res.Sha1 == "" {
		return "", fmt.Errorf("la API no devolvió el token creado")
	}
	return res.Sha1, nil
}

// CreateRepoOptions son los campos aceptados al crear un repositorio.
type CreateRepoOptions struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Private       bool   `json:"private"`
	Unlisted      bool   `json:"unlisted"`
	AutoInit      bool   `json:"auto_init"`
	Readme        string `json:"readme,omitempty"`
	Gitignores    string `json:"gitignores,omitempty"`
	License       string `json:"license,omitempty"`
	DefaultBranch string `json:"default_branch,omitempty"`
}

// CreateRepo crea un repositorio en la cuenta del usuario autenticado.
func (c *Client) CreateRepo(ctx context.Context, opts CreateRepoOptions) (*apiv2.RepositoryResponse, error) {
	var repo apiv2.RepositoryResponse
	if err := c.Post(ctx, "/user/repos", opts, &repo); err != nil {
		return nil, err
	}
	return &repo, nil
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

// CreateIssueOptions son los campos aceptados al crear una incidencia.
type CreateIssueOptions struct {
	Title       string  `json:"title"`
	Content     string  `json:"content,omitempty"`
	MilestoneID int64   `json:"milestone_id,omitempty"`
	AssigneeID  int64   `json:"assignee_id,omitempty"`
	LabelIDs    []int64 `json:"label_ids,omitempty"`
}

// CreateIssue crea una incidencia y devuelve la incidencia persistida.
func (c *Client) CreateIssue(ctx context.Context, owner, name string, opts CreateIssueOptions) (*apiv2.IssueResponse, error) {
	var issue apiv2.IssueResponse
	path := repoPath(owner, name, "/issues")
	if err := c.Post(ctx, path, opts, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

// CreateIssueComment añade un comentario a una incidencia del repositorio.
func (c *Client) CreateIssueComment(ctx context.Context, owner, name string, index int64, body string) (*apiv2.IssueCommentResponse, error) {
	var comment apiv2.IssueCommentResponse
	path := repoPath(owner, name, "/issues/"+strconv.FormatInt(index, 10)+"/comments")
	if err := c.Post(ctx, path, map[string]string{"body": body}, &comment); err != nil {
		return nil, err
	}
	return &comment, nil
}

// UpdateIssueOptions son los campos modificables de una incidencia.
type UpdateIssueOptions struct {
	IsClosed *bool `json:"is_closed"`
}

// SetIssueClosed cierra o reabre una incidencia. El servidor responde con el
// modelo interno (no con el DTO de la API), así que devolvemos el JSON tal cual.
func (c *Client) SetIssueClosed(ctx context.Context, owner, name string, index int64, closed bool) (json.RawMessage, error) {
	var raw json.RawMessage
	path := repoPath(owner, name, "/issues/"+strconv.FormatInt(index, 10))
	if err := c.Put(ctx, path, UpdateIssueOptions{IsClosed: &closed}, &raw); err != nil {
		return nil, err
	}
	return raw, nil
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
