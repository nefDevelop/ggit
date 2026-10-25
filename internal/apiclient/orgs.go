package apiclient

import (
	"context"
	"net/url"
)

// Org es el subconjunto de GET /orgs/{org} que usa el CLI.
type Org struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

// GetOrg devuelve la organización por nombre. Hace falta su ID para crear
// repositorios a su nombre.
func (c *Client) GetOrg(ctx context.Context, name string) (*Org, error) {
	var org Org
	if err := c.Get(ctx, "/orgs/"+url.PathEscape(name), nil, &org); err != nil {
		return nil, err
	}
	return &org, nil
}
