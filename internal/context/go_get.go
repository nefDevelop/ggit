package context

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/repox"
	"gogs.io/gogs/internal/strx"
)

// ServeGoGet does quick responses for appropriate go-get meta with status OK
// regardless of whether the user has access to the repository, or the repository
// does exist at all. This is particular a workaround for "go get" command which
// does not respect .netrc file.
func ServeGoGet() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("go-get") != "1" {
			return
		}

		ownerName := c.Param("username")
		repoName := strings.TrimSuffix(c.Param("reponame"), ".git")
		branchName := "master"

		owner, err := database.Handle.Users().GetByUsername(c.Request.Context(), ownerName)
		if err == nil {
			repo, err := database.Handle.Repositories().GetByName(c.Request.Context(), owner.ID, repoName)
			if err == nil && repo.DefaultBranch != "" {
				branchName = repo.DefaultBranch
			}
		}

		prefix := conf.Server.ExternalURL + path.Join(ownerName, repoName, "src", branchName)
		insecureFlag := ""
		if !strings.HasPrefix(conf.Server.ExternalURL, "https://") {
			insecureFlag = "--insecure "
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(strx.Expand(`<!doctype html>
<html>
	<head>
		<meta name="go-import" content="{GoGetImport} git {CloneLink}">
		<meta name="go-source" content="{GoGetImport} _ {GoDocDirectory} {GoDocFile}">
	</head>
	<body>
		go get {InsecureFlag}{GoGetImport}
	</body>
</html>
`,
			map[string]string{
				"GoGetImport":    path.Join(conf.Server.URL.Host, conf.Server.Subpath, ownerName, repoName),
				"CloneLink":      repox.HTTPSCloneURL(ownerName, repoName),
				"GoDocDirectory": prefix + "{/dir}",
				"GoDocFile":      prefix + "{/dir}/{file}#L{line}",
				"InsecureFlag":   insecureFlag,
			},
		)))
		c.Abort()
	}
}
