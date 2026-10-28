package v2

import (
	"net/http"
	"sort"
	"testing"

	"github.com/gin-gonic/gin"
)

// reviewedMutatingRoutes is the inventory of every state-changing route this
// package registers, together with the authorization middleware that guards it.
//
// This exists because the gap it closes is exactly how the API shipped without
// access control: routes kept being added to RegisterRoutes and nothing forced
// anyone to state which middleware protects them. Adding a mutating route now
// fails this test until it is listed here, which makes the decision explicit and
// reviewable in the diff.
//
// Entry format is "METHOD /path". The middleware is not introspectable from
// gin.Routes(), so this test cannot assert the guard itself; it asserts that the
// surface is deliberate. The guards themselves are covered by the unit tests in
// internal/context (authz_test.go).
var reviewedMutatingRoutes = map[string]bool{
	"DELETE /api/v2/:username/:reponame/branches/:branch":                        true,
	"DELETE /api/v2/:username/:reponame/branches/:branch/protection":             true,
	"DELETE /api/v2/:username/:reponame/issues/:index/labels":                    true,
	"DELETE /api/v2/:username/:reponame/labels/:id":                              true,
	"DELETE /api/v2/:username/:reponame/milestones/:id":                          true,
	"DELETE /api/v2/:username/:reponame/releases/:id":                            true,
	"DELETE /api/v2/:username/:reponame/settings":                                true,
	"DELETE /api/v2/:username/:reponame/settings/avatar":                         true,
	"DELETE /api/v2/:username/:reponame/settings/collaborators/:collaborator":    true,
	"DELETE /api/v2/:username/:reponame/settings/hooks/:id":                      true,
	"DELETE /api/v2/:username/:reponame/settings/keys/:id":                       true,
	"DELETE /api/v2/:username/:reponame/settings/pikoci":                         true,
	"DELETE /api/v2/:username/:reponame/settings/remotes/:name":                  true,
	"DELETE /api/v2/:username/:reponame/settings/security/guardian/patterns/:id": true,
	"DELETE /api/v2/:username/:reponame/tags/:tag":                               true,
	"DELETE /api/v2/:username/:reponame/wiki/:page":                              true,
	"DELETE /api/v2/admin/auths/:id":                                             true,
	"DELETE /api/v2/admin/notices":                                               true,
	"DELETE /api/v2/admin/notices/:id":                                           true,
	"DELETE /api/v2/admin/pikoci/repos/:repoID":                                  true,
	"DELETE /api/v2/admin/repos/:id":                                             true,
	"DELETE /api/v2/admin/security/guardian/patterns/:id":                        true,
	"DELETE /api/v2/admin/users/:id":                                             true,
	"DELETE /api/v2/gists/:uuid":                                                 true,
	"DELETE /api/v2/orgs/:org":                                                   true,
	"DELETE /api/v2/orgs/:org/avatar":                                            true,
	"DELETE /api/v2/orgs/:org/members/:username":                                 true,
	"DELETE /api/v2/orgs/teams/:teamID":                                          true,
	"DELETE /api/v2/orgs/teams/:teamID/members/:username":                        true,
	"DELETE /api/v2/orgs/teams/:teamID/repos/:repoID":                            true,
	"DELETE /api/v2/user/settings/account":                                       true,
	"DELETE /api/v2/user/settings/avatar":                                        true,
	"DELETE /api/v2/user/settings/blocks/:username":                              true,
	"DELETE /api/v2/user/settings/emails":                                        true,
	"DELETE /api/v2/user/settings/orgs/:id":                                      true,
	"DELETE /api/v2/user/settings/repos/:id":                                     true,
	"DELETE /api/v2/user/settings/security/two_factor":                           true,
	"DELETE /api/v2/user/settings/sessions":                                      true,
	"DELETE /api/v2/user/settings/sessions/:id":                                  true,
	"DELETE /api/v2/user/settings/ssh_keys/:id":                                  true,
	"DELETE /api/v2/user/settings/tokens/:id":                                    true,
	"PATCH /api/v2/:username/:reponame/labels/:id":                               true,
	"PATCH /api/v2/:username/:reponame/milestones/:id":                           true,
	"PATCH /api/v2/:username/:reponame/releases/:id":                             true,
	"PATCH /api/v2/:username/:reponame/settings":                                 true,
	"PATCH /api/v2/:username/:reponame/settings/collaborators/:collaborator":     true,
	"PATCH /api/v2/:username/:reponame/settings/hooks/:id":                       true,
	"PATCH /api/v2/:username/:reponame/settings/hooks/git/:name":                 true,
	"PATCH /api/v2/:username/:reponame/settings/security/guardian":               true,
	"PATCH /api/v2/admin/auths/:id":                                              true,
	"PATCH /api/v2/admin/users/:id":                                              true,
	"PATCH /api/v2/orgs/:org":                                                    true,
	"PATCH /api/v2/orgs/teams/:teamID":                                           true,
	"POST /api/v2/:username/:reponame/branches":                                  true,
	"POST /api/v2/:username/:reponame/editor/delete":                             true,
	"POST /api/v2/:username/:reponame/editor/edit":                               true,
	"POST /api/v2/:username/:reponame/issues":                                    true,
	"POST /api/v2/:username/:reponame/issues/:index/comments":                    true,
	"POST /api/v2/:username/:reponame/issues/:index/labels":                      true,
	"POST /api/v2/:username/:reponame/labels":                                    true,
	"POST /api/v2/:username/:reponame/milestones":                                true,
	"POST /api/v2/:username/:reponame/pulls":                                     true,
	"POST /api/v2/:username/:reponame/pulls/:index/merge":                        true,
	"POST /api/v2/:username/:reponame/pulls/:index/reviews":                      true,
	"POST /api/v2/:username/:reponame/pulls/:index/reviews/:id/comments":         true,
	"POST /api/v2/:username/:reponame/pulls/comments/:id/resolve":                true,
	"POST /api/v2/:username/:reponame/releases":                                  true,
	"POST /api/v2/:username/:reponame/releases/attachments":                      true,
	"POST /api/v2/:username/:reponame/settings/avatar":                           true,
	"POST /api/v2/:username/:reponame/settings/collaborators/:collaborator":      true,
	"POST /api/v2/:username/:reponame/settings/hooks":                            true,
	"POST /api/v2/:username/:reponame/settings/hooks/:id/redeliver/:uuid":        true,
	"POST /api/v2/:username/:reponame/settings/hooks/:id/test":                   true,
	"POST /api/v2/:username/:reponame/settings/janitor/cleanup":                  true,
	"POST /api/v2/:username/:reponame/settings/janitor/scan":                     true,
	"POST /api/v2/:username/:reponame/settings/keys":                             true,
	"POST /api/v2/:username/:reponame/settings/pikoci":                           true,
	"POST /api/v2/:username/:reponame/settings/remotes":                          true,
	"POST /api/v2/:username/:reponame/settings/security/guardian/patterns":       true,
	"POST /api/v2/:username/:reponame/settings/transfer":                         true,
	"POST /api/v2/:username/:reponame/statuses/:sha":                             true,
	"POST /api/v2/:username/:reponame/tags":                                      true,
	"POST /api/v2/:username/:reponame/wiki":                                      true,
	"POST /api/v2/admin/auths":                                                   true,
	"POST /api/v2/admin/config/test_email":                                       true,
	"POST /api/v2/admin/janitor/scan-all":                                        true,
	"POST /api/v2/admin/monitor/cron":                                            true,
	"POST /api/v2/admin/monitor/processes/:pid":                                  true,
	"POST /api/v2/admin/operations":                                              true,
	"POST /api/v2/admin/orgs":                                                    true,
	"POST /api/v2/admin/pikoci/repos/:repoID/sync":                               true,
	"POST /api/v2/admin/repos/:id/check-integrity":                               true,
	"POST /api/v2/admin/repos/adopt":                                             true,
	"POST /api/v2/admin/security/guardian/patterns":                              true,
	"POST /api/v2/admin/users":                                                   true,
	"POST /api/v2/gists":                                                         true,
	"POST /api/v2/install":                                                       true,
	"POST /api/v2/orgs/:org/avatar":                                              true,
	"POST /api/v2/orgs/:org/teams":                                               true,
	"POST /api/v2/orgs/teams/:teamID/members/:username":                          true,
	"POST /api/v2/orgs/teams/:teamID/repos/:reponame":                            true,
	"POST /api/v2/user/activate":                                                 true,
	"POST /api/v2/user/forgot_password":                                          true,
	"POST /api/v2/user/login":                                                    true,
	"POST /api/v2/user/login/two_factor":                                         true,
	"POST /api/v2/user/logout":                                                   true,
	"POST /api/v2/user/register":                                                 true,
	"POST /api/v2/user/repos":                                                    true,
	"POST /api/v2/user/repos/migrate":                                            true,
	"POST /api/v2/user/reset_password":                                           true,
	"POST /api/v2/user/settings/avatar":                                          true,
	"POST /api/v2/user/settings/avatar/magic":                                    true,
	"POST /api/v2/user/settings/blocks/:username":                                true,
	"POST /api/v2/user/settings/emails":                                          true,
	"POST /api/v2/user/settings/orgs":                                            true,
	"POST /api/v2/user/settings/security/two_factor/confirm":                     true,
	"POST /api/v2/user/settings/security/two_factor/enable":                      true,
	"POST /api/v2/user/settings/security/two_factor/recovery_codes":              true,
	"POST /api/v2/user/settings/ssh_keys":                                        true,
	"POST /api/v2/user/settings/tokens":                                          true,
	"PUT /api/v2/:username/:reponame/branches/:branch/protection":                true,
	"PUT /api/v2/:username/:reponame/issues/:index":                              true,
	"PUT /api/v2/:username/:reponame/issues/:index/labels":                       true,
	"PUT /api/v2/:username/:reponame/wiki/:page":                                 true,
	"PUT /api/v2/admin/config":                                                   true,
	"PUT /api/v2/admin/pikoci/repos/:repoID":                                     true,
	"PUT /api/v2/gists/:uuid":                                                    true,
	"PUT /api/v2/user/settings/emails/primary":                                   true,
	"PUT /api/v2/user/settings/password":                                         true,
	"PUT /api/v2/user/settings/profile":                                          true,
}

func TestEveryMutatingRouteIsReviewed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	RegisterRoutes(router.Group("/api/v2"))

	var registered []string
	for _, route := range router.Routes() {
		switch route.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			registered = append(registered, route.Method+" "+route.Path)
		}
	}
	sort.Strings(registered)

	var unreviewed []string
	for _, key := range registered {
		if !reviewedMutatingRoutes[key] {
			unreviewed = append(unreviewed, key)
		}
	}

	if len(unreviewed) > 0 {
		t.Errorf("rutas mutantes sin revisar (%d). Anadelas a reviewedMutatingRoutes y confirma que su middleware de autorizacion es el adecuado:\n  %s",
			len(unreviewed), joinLines(unreviewed))
	}

	var removed []string
	for key := range reviewedMutatingRoutes {
		found := false
		for _, got := range registered {
			if got == key {
				found = true
				break
			}
		}
		if !found {
			removed = append(removed, key)
		}
	}
	if len(removed) > 0 {
		sort.Strings(removed)
		t.Logf("rutas retiradas del inventario (quitarlas de reviewedMutatingRoutes): %s", joinLines(removed))
	}
}

func joinLines(items []string) string {
	out := ""
	for _, item := range items {
		out += "\n    - " + item
	}
	return out
}
