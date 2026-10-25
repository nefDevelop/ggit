package apiclient

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"

	apiv2 "github.com/nef734/ggit/internal/route/api/v2/types"
)

// ListPullRequestsOptions filtra el listado de pull requests de un repositorio.
type ListPullRequestsOptions struct {
	State string
	Page  int
}

// ListPullRequests devuelve los pull requests de un repositorio.
func (c *Client) ListPullRequests(ctx context.Context, owner, name string, opts ListPullRequestsOptions) ([]*apiv2.PullRequestResponse, error) {
	q := url.Values{}
	if opts.State != "" {
		q.Set("state", opts.State)
	}
	if opts.Page > 1 {
		q.Set("page", strconv.Itoa(opts.Page))
	}

	var pulls []*apiv2.PullRequestResponse
	if err := c.Get(ctx, repoPath(owner, name, "/pulls"), q, &pulls); err != nil {
		return nil, err
	}
	return pulls, nil
}

// GetPullRequest devuelve un pull request por su índice en el repositorio.
func (c *Client) GetPullRequest(ctx context.Context, owner, name string, index int64) (*apiv2.PullRequestResponse, error) {
	var pull apiv2.PullRequestResponse
	path := repoPath(owner, name, "/pulls/"+strconv.FormatInt(index, 10))
	if err := c.Get(ctx, path, nil, &pull); err != nil {
		return nil, err
	}
	return &pull, nil
}

// CreatePullRequestOptions son los campos de un pull request nuevo.
type CreatePullRequestOptions struct {
	Title      string `json:"title"`
	Content    string `json:"content,omitempty"`
	HeadBranch string `json:"head_branch"`
	BaseBranch string `json:"base_branch"`
}

// CreatePullRequest abre un pull request. El servidor responde con el modelo
// interno de la incidencia, así que devolvemos el JSON tal cual.
func (c *Client) CreatePullRequest(ctx context.Context, owner, name string, opts CreatePullRequestOptions) (json.RawMessage, error) {
	var raw json.RawMessage
	path := repoPath(owner, name, "/pulls")
	if err := c.Post(ctx, path, opts, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// ListPullReviews devuelve las revisiones de un pull request.
func (c *Client) ListPullReviews(ctx context.Context, owner, name string, index int64) ([]*apiv2.PullReviewResponse, error) {
	var reviews []*apiv2.PullReviewResponse
	path := repoPath(owner, name, "/pulls/"+strconv.FormatInt(index, 10)+"/reviews")
	if err := c.Get(ctx, path, nil, &reviews); err != nil {
		return nil, err
	}
	return reviews, nil
}

// GetPullRequestDiff devuelve el diff parseado de un pull request.
func (c *Client) GetPullRequestDiff(ctx context.Context, owner, name string, index int64) (*Diff, error) {
	var diff Diff
	path := repoPath(owner, name, "/pulls/"+strconv.FormatInt(index, 10)+"/diff")
	if err := c.Get(ctx, path, nil, &diff); err != nil {
		return nil, err
	}
	return &diff, nil
}

// ReviewType es el tipo de revisión tal como lo espera el servidor.
type ReviewType int

// Tipos de revisión aceptados por POST /pulls/{index}/reviews.
const (
	ReviewApprove ReviewType = 1
	ReviewComment ReviewType = 2
	ReviewReject  ReviewType = 3
)

// CreatePullReviewOptions son los campos de una revisión nueva.
type CreatePullReviewOptions struct {
	Type    ReviewType `json:"type"`
	Content string     `json:"content,omitempty"`
}

// CreatePullReview envía una revisión. El servidor responde con el modelo
// interno de la revisión, así que devolvemos el JSON tal cual.
func (c *Client) CreatePullReview(ctx context.Context, owner, name string, index int64, opts CreatePullReviewOptions) (json.RawMessage, error) {
	var raw json.RawMessage
	path := repoPath(owner, name, "/pulls/"+strconv.FormatInt(index, 10)+"/reviews")
	if err := c.Post(ctx, path, opts, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

// MergePullRequest fusiona un pull request y devuelve el estado informado por
// el servidor (normalmente "merged").
func (c *Client) MergePullRequest(ctx context.Context, owner, name string, index int64) (string, error) {
	var res struct {
		Status string `json:"status"`
	}
	path := repoPath(owner, name, "/pulls/"+strconv.FormatInt(index, 10)+"/merge")
	if err := c.Post(ctx, path, nil, &res); err != nil {
		return "", err
	}
	return res.Status, nil
}

// Diff refleja el JSON de /pulls/{index}/diff. El servidor serializa
// internal/gitx.Diff (gogs/git-module) sin tags JSON, así que los nombres de
// campo son los de Go y solo los exportados llegan al cliente.
type Diff struct {
	Files []*DiffFile
}

// DiffFile es un archivo dentro del diff.
type DiffFile struct {
	Name     string
	Type     DiffFileType
	Index    string
	OldIndex string
	Sections []*DiffSection
}

// DiffSection agrupa las líneas de un bloque del diff.
type DiffSection struct {
	Lines []*DiffLine
}

// DiffLine es una línea del diff. Content ya incluye el prefijo (+, -, espacio).
type DiffLine struct {
	Type      DiffLineType
	Content   string
	LeftLine  int
	RightLine int
}

// DiffFileType y DiffLineType replican los enums de gogs/git-module que el
// servidor envía como enteros.
type DiffFileType int

// DiffLineType es el tipo de una línea del diff.
type DiffLineType int

const (
	DiffFileAdd    DiffFileType = 1
	DiffFileChange DiffFileType = 2
	DiffFileDelete DiffFileType = 3
	DiffFileRename DiffFileType = 4

	DiffLinePlain   DiffLineType = 1
	DiffLineAdd     DiffLineType = 2
	DiffLineDelete  DiffLineType = 3
	DiffLineSection DiffLineType = 4
)

// Stats cuenta las líneas añadidas y eliminadas del archivo.
func (f *DiffFile) Stats() (additions, deletions int) {
	for _, section := range f.Sections {
		for _, line := range section.Lines {
			switch line.Type {
			case DiffLineAdd:
				additions++
			case DiffLineDelete:
				deletions++
			}
		}
	}
	return additions, deletions
}

// TypeName devuelve el nombre legible del tipo de cambio.
func (f *DiffFile) TypeName() string {
	switch f.Type {
	case DiffFileAdd:
		return "add"
	case DiffFileChange:
		return "change"
	case DiffFileDelete:
		return "delete"
	case DiffFileRename:
		return "rename"
	default:
		return "unknown"
	}
}
