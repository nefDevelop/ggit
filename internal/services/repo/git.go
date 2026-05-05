package repo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/gogs/git-module"

	embedConf "github.com/nef734/ggit/conf"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/osx"
	"github.com/nef734/ggit/internal/process"
	"github.com/nef734/ggit/internal/repox"
	"github.com/nef734/ggit/internal/strx"
)

func getRepoInitFile(tp, name string) ([]byte, error) {
	relPath := filepath.Join(tp, strings.TrimLeft(filepath.Clean("/"+name), "/"))

	// Use custom file when available.
	customPath := filepath.Join(conf.CustomDir(), "conf", relPath)
	if osx.IsFile(customPath) {
		return os.ReadFile(customPath)
	}
	return embedConf.Files.ReadFile(relPath)
}

func prepareRepoCommit(repo *database.Repository, tmpDir, repoPath string, opts CreateOptions) error {
	// Clone to temporary path and do the init commit.
	err := git.Clone(repoPath, tmpDir, git.CloneOptions{})
	if err != nil {
		return errors.Wrap(err, "clone")
	}

	// README
	readmeName := opts.Readme
	if readmeName == "" {
		readmeName = "Default"
	}
	data, err := getRepoInitFile("readme", readmeName)
	if err != nil {
		return errors.Newf("getRepoInitFile[%s]: %v", readmeName, err)
	}

	cloneLink := repox.NewCloneLink(repo.Owner.Name, repo.Name)
	match := map[string]string{
		"Name":           repo.Name,
		"Description":    repo.Description,
		"CloneURL.SSH":   cloneLink.SSH,
		"CloneURL.HTTPS": cloneLink.HTTPS,
	}
	if err = os.WriteFile(filepath.Join(tmpDir, "README.md"),
		[]byte(strx.Expand(string(data), match)), 0o644); err != nil {
		return errors.Newf("write README.md: %v", err)
	}

	// .gitignore
	if len(opts.Gitignores) > 0 {
		var buf bytes.Buffer
		names := strings.Split(opts.Gitignores, ",")
		for _, name := range names {
			data, err = getRepoInitFile("gitignore", name)
			if err != nil {
				return errors.Newf("getRepoInitFile[%s]: %v", name, err)
			}
			buf.WriteString("# ---> " + name + "\n")
			buf.Write(data)
			buf.WriteString("\n")
		}

		if buf.Len() > 0 {
			if err = os.WriteFile(filepath.Join(tmpDir, ".gitignore"), buf.Bytes(), 0o644); err != nil {
				return errors.Newf("write .gitignore: %v", err)
			}
		}
	}

	// LICENSE
	if len(opts.License) > 0 {
		data, err = getRepoInitFile("license", opts.License)
		if err != nil {
			return errors.Newf("getRepoInitFile[%s]: %v", opts.License, err)
		}

		if err = os.WriteFile(filepath.Join(tmpDir, "LICENSE"), data, 0o644); err != nil {
			return errors.Newf("write LICENSE: %v", err)
		}
	}

	return nil
}

func initRepoCommit(tmpDir string, sig *git.Signature, branch string) error {
	var stderr string
	var err error
	if _, stderr, err = process.ExecDir(-1,
		tmpDir, fmt.Sprintf("initRepoCommit (git add): %s", tmpDir),
		"git", "add", "--all"); err != nil {
		return errors.Newf("git add: %s", stderr)
	}

	if _, stderr, err = process.ExecDir(-1,
		tmpDir, fmt.Sprintf("initRepoCommit (git commit): %s", tmpDir),
		"git", "commit", fmt.Sprintf("--author='%s <%s>'", sig.Name, sig.Email),
		"-m", "Initial commit"); err != nil {
		return errors.Newf("git commit: %s", stderr)
	}

	if _, stderr, err = process.ExecDir(-1,
		tmpDir, fmt.Sprintf("initRepoCommit (git push): %s", tmpDir),
		"git", "push", "origin", branch); err != nil {
		return errors.Newf("git push: %s", stderr)
	}
	return nil
}

func initRepository(repoPath string, doer *database.User, repo *database.Repository, opts CreateOptions) error {
	// Init bare new repository.
	if err := git.Init(repoPath, git.InitOptions{Bare: true}); err != nil {
		return errors.Newf("init repository: %v", err)
	} else if err := repox.CreateDelegateHooks(repoPath); err != nil {
		return errors.Newf("CreateDelegateHooks: %v", err)
	}

	// Set default branch
	_, err := git.SymbolicRef(
		repoPath,
		git.SymbolicRefOptions{
			Name: "HEAD",
			Ref:  git.RefsHeads + repo.DefaultBranch,
		},
	)
	if err != nil {
		return errors.Wrap(err, "set default branch")
	}

	if opts.AutoInit {
		tmpDir := filepath.Join(os.TempDir(), "gogs-"+repo.Name+"-"+strconv.Itoa(time.Now().Nanosecond()))
		if err = os.MkdirAll(tmpDir, os.ModePerm); err != nil {
			return err
		}
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		if err = prepareRepoCommit(repo, tmpDir, repoPath, opts); err != nil {
			return errors.Newf("prepareRepoCommit: %v", err)
		}

		// Apply changes and commit.
		err = initRepoCommit(
			tmpDir,
			&git.Signature{
				Name:  doer.DisplayName(),
				Email: doer.Email,
				When:  time.Now(),
			},
			repo.DefaultBranch,
		)
		if err != nil {
			return errors.Newf("initRepoCommit: %v", err)
		}
	}

	return nil
}
