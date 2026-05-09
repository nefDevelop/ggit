package repo

import (
	"context"
	json "github.com/nef734/ggit/internal/jsonx"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/database"
	log "github.com/nef734/ggit/internal/log"
)

// UpdateRepositoryLanguages analyzes the repository and updates its language statistics.
func UpdateRepositoryLanguages(repo *database.Repository) error {
	gitRepo, err := git.Open(repo.RepoPath())
	if err != nil {
		return err
	}

	commit, err := gitRepo.BranchCommit(repo.DefaultBranch)
	if err != nil {
		return nil
	}

	languages := make(map[string]int64)
	err = walkTree(commit.Tree, "", languages)
	if err != nil {
		return err
	}

	var primary string
	var maxBytes int64
	for lang, bytes := range languages {
		if bytes > maxBytes {
			maxBytes = bytes
			primary = lang
		}
	}

	if primary != "" {
		log.Trace("Updating languages for %s: %s", repo.FullName(), primary)
		langJSON, _ := json.Marshal(languages)
		langStr := string(langJSON)

		err := database.Handle.Repositories().Update(context.Background(), repo.ID, database.UpdateRepoOptions{
			PrimaryLanguage: &primary,
			Languages:       &langStr,
		})
		if err != nil {
			log.Error("Failed to update language stats: %v", err)
			return err
		}
		repo.PrimaryLanguage = primary
		repo.Languages = langStr
	}

	return nil
}

func walkTree(tree *git.Tree, path string, languages map[string]int64) error {
	entries, err := tree.Entries()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsTree() {
			subtree, err := tree.Subtree(entry.Name())
			if err != nil {
				return err
			}
			err = walkTree(subtree, path+entry.Name()+"/", languages)
			if err != nil {
				return err
			}
			continue
		}

		if isVendor(entry.Name()) || isDotFile(entry.Name()) || isDocumentation(entry.Name()) || isConfiguration(entry.Name()) {
			continue
		}

		lang := getLanguageByExtension(entry.Name())
		if lang == "" {
			continue
		}

		languages[lang] += entry.Size()
	}
	return nil
}
