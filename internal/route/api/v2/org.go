package v2

import (
	"net/http"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/errx"
)

// GetOrg returns the organization profile.
func GetOrg(c *context.APIContext) {
	orgName := c.Params("org")
	org, err := database.Handle.Users().GetByUsername(c.Req.Context(), orgName)
	if err != nil {
		if errx.IsNotFound(err) {
			c.NotFound()
		} else {
			c.Error(err, "get organization")
		}
		return
	}

	if !org.IsOrganization() {
		c.NotFound()
		return
	}

	c.JSON(http.StatusOK, org)
}

// UpdateOrg updates the organization profile.
func UpdateOrg(c *context.APIContext) {
	orgName := c.Params("org")
	org, err := database.Handle.Users().GetByUsername(c.Req.Context(), orgName)
	if err != nil {
		c.NotFoundOrError(err, "get organization")
		return
	}

	var f struct {
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		Website     string `json:"website"`
		Location    string `json:"location"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	opts := database.UpdateUserOptions{
		FullName:    &f.FullName,
		Description: &f.Description,
		Website:     &f.Website,
		Location:    &f.Location,
	}

	if err := database.Handle.Users().Update(c.Req.Context(), org.ID, opts); err != nil {
		c.Error(err, "update organization")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// DeleteOrg deletes the organization.
func DeleteOrg(c *context.APIContext) {
	orgName := c.Params("org")
	org, err := database.Handle.Users().GetByUsername(c.Req.Context(), orgName)
	if err != nil {
		c.NotFoundOrError(err, "get organization")
		return
	}

	if err := database.DeleteOrganization(org); err != nil {
		c.Error(err, "delete organization")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// ListOrgMembers returns all members of the organization.
func ListOrgMembers(c *context.APIContext) {
	orgName := c.Params("org")
	org, err := database.Handle.Users().GetByUsername(c.Req.Context(), orgName)
	if err != nil {
		c.NotFoundOrError(err, "get organization")
		return
	}

	members, err := database.GetOrgUsersByOrgID(org.ID, 0)
	if err != nil {
		c.Error(err, "list members")
		return
	}

	c.JSON(http.StatusOK, members)
}

// RemoveOrgMember removes a user from the organization.
func RemoveOrgMember(c *context.APIContext) {
	orgName := c.Params("org")
	userName := c.Params("username")

	org, err := database.Handle.Users().GetByUsername(c.Req.Context(), orgName)
	if err != nil {
		c.NotFoundOrError(err, "get organization")
		return
	}

	u, err := database.Handle.Users().GetByUsername(c.Req.Context(), userName)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}

	if err := database.RemoveOrgUser(org.ID, u.ID); err != nil {
		if database.IsErrLastOrgOwner(err) {
			c.StatusError(http.StatusUnprocessableEntity, "Cannot remove the last owner.", err)
		} else {
			c.Error(err, "remove organization user")
		}
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// ListOrgTeams returns all teams of the organization.
func ListOrgTeams(c *context.APIContext) {
	orgName := c.Params("org")
	org, err := database.Handle.Users().GetByUsername(c.Req.Context(), orgName)
	if err != nil {
		c.NotFoundOrError(err, "get organization")
		return
	}

	teams, err := database.GetTeamsByOrgID(org.ID)
	if err != nil {
		c.Error(err, "list teams")
		return
	}

	for _, t := range teams {
		_ = t.GetMembers()
		_ = t.GetRepositories()
	}

	c.JSON(http.StatusOK, teams)
}

// GetTeam returns a single team with its members and repositories.
func GetTeam(c *context.APIContext) {
	teamID := c.ParamsInt64("teamID")
	t, err := database.GetTeamByID(teamID)
	if err != nil {
		c.NotFoundOrError(err, "get team")
		return
	}

	_ = t.GetMembers()
	_ = t.GetRepositories()

	// Populate API format for members
	for _, m := range t.Members {
		// Ensure AvatarURL is loaded or formatted
		_ = m.AvatarURL()
	}

	c.JSON(http.StatusOK, t)
}

// CreateOrgTeam creates a new team in the organization.
func CreateOrgTeam(c *context.APIContext) {
	orgName := c.Params("org")
	org, err := database.Handle.Users().GetByUsername(c.Req.Context(), orgName)
	if err != nil {
		c.NotFoundOrError(err, "get organization")
		return
	}

	var f struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Permission  string `json:"permission"` // read, write, admin
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	t := &database.Team{
		OrgID:       org.ID,
		Name:        f.Name,
		Description: f.Description,
		Authorize:   database.ParseAccessMode(f.Permission),
	}

	if err := database.NewTeam(t); err != nil {
		c.Error(err, "create team")
		return
	}

	c.JSON(http.StatusCreated, t)
}

// UpdateOrgTeam updates a team's information.
func UpdateOrgTeam(c *context.APIContext) {
	teamID := c.ParamsInt64("teamID")
	t, err := database.GetTeamByID(teamID)
	if err != nil {
		c.NotFoundOrError(err, "get team")
		return
	}

	var f struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Permission  string `json:"permission"`
	}
	if err := c.ShouldBindJSON(&f); err != nil {
		c.ErrorStatus(http.StatusBadRequest, err)
		return
	}

	var authChanged bool
	if f.Name != "" {
		t.Name = f.Name
	}
	if f.Description != "" {
		t.Description = f.Description
	}
	if f.Permission != "" {
		newMode := database.ParseAccessMode(f.Permission)
		if newMode != t.Authorize {
			t.Authorize = newMode
			authChanged = true
		}
	}

	if err := database.UpdateTeam(t, authChanged); err != nil {
		c.Error(err, "update team")
		return
	}

	c.JSON(http.StatusOK, t)
}

// DeleteOrgTeam deletes a team from the organization.
func DeleteOrgTeam(c *context.APIContext) {
	teamID := c.ParamsInt64("teamID")
	t, err := database.GetTeamByID(teamID)
	if err != nil {
		c.NotFoundOrError(err, "get team")
		return
	}

	if err := database.DeleteTeam(t); err != nil {
		c.Error(err, "delete team")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// AddTeamMember adds a user to a team.
func AddTeamMember(c *context.APIContext) {
	teamID := c.ParamsInt64("teamID")
	userName := c.Params("username")

	t, err := database.GetTeamByID(teamID)
	if err != nil {
		c.NotFoundOrError(err, "get team")
		return
	}

	u, err := database.Handle.Users().GetByUsername(c.Req.Context(), userName)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}

	if err := t.AddMember(u.ID); err != nil {
		c.Error(err, "add team member")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// RemoveTeamMember removes a user from a team.
func RemoveTeamMember(c *context.APIContext) {
	teamID := c.ParamsInt64("teamID")
	userName := c.Params("username")

	t, err := database.GetTeamByID(teamID)
	if err != nil {
		c.NotFoundOrError(err, "get team")
		return
	}

	u, err := database.Handle.Users().GetByUsername(c.Req.Context(), userName)
	if err != nil {
		c.NotFoundOrError(err, "get user")
		return
	}

	if err := t.RemoveMember(u.ID); err != nil {
		c.Error(err, "remove team member")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// AddTeamRepo adds a repository to a team.
func AddTeamRepo(c *context.APIContext) {
	teamID := c.ParamsInt64("teamID")
	repoName := c.Params("reponame")

	t, err := database.GetTeamByID(teamID)
	if err != nil {
		c.NotFoundOrError(err, "get team")
		return
	}

	repo, err := database.GetRepositoryByName(t.OrgID, repoName)
	if err != nil {
		c.NotFoundOrError(err, "get repository")
		return
	}

	if err := t.AddRepository(repo); err != nil {
		c.Error(err, "add team repository")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// RemoveTeamRepo removes a repository from a team.
func RemoveTeamRepo(c *context.APIContext) {
	teamID := c.ParamsInt64("teamID")
	repoID := c.ParamsInt64("repoID")

	t, err := database.GetTeamByID(teamID)
	if err != nil {
		c.NotFoundOrError(err, "get team")
		return
	}

	if err := t.RemoveRepository(repoID); err != nil {
		c.Error(err, "remove team repository")
		return
	}

	c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
