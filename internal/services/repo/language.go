package repo

import (
	"context"
	"encoding/json"

	"github.com/go-enry/go-enry/v2"
	"github.com/gogs/git-module"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/database"
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

		if enry.IsVendor(entry.Name()) || enry.IsDotFile(entry.Name()) || enry.IsDocumentation(entry.Name()) || enry.IsConfiguration(entry.Name()) {
			continue
		}

		lang, _ := enry.GetLanguageByExtension(entry.Name())
		if lang == "" {
			continue
		}

		languages[lang] += entry.Size()
	}
	return nil
}
