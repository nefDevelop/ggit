package pikoci

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/gogs/git-module"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
)

// ErrPipelineFileNotFound indica que el repositorio no define pipeline en la
// ruta configurada.
var ErrPipelineFileNotFound = errors.New("pipeline file not found")

// PipelineName deriva el nombre canónico del pipeline a partir de owner/name.
func PipelineName(owner, name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(owner + "-" + name) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

// SyncRepoPipeline publica el pipeline del repositorio con sus variables y
// guarda el webhook token. Devuelve el error real: lo llama una acción
// explícita del administrador, no un push.
func (s *Service) SyncRepoPipeline(ctx context.Context, repoID int64, owner, name, branch string) (*database.PikoCIRepo, error) {
	pipeline := PipelineName(owner, name)

	hcl, err := readPipelineFile(owner, name, branch, conf.PikoCI.PipelineFile)
	if err != nil {
		return nil, err
	}

	externalURL := strings.TrimRight(conf.Server.ExternalURL, "/")
	vars := map[string]string{
		"ggit_url":     externalURL,
		"ggit_token":   conf.PikoCI.BotToken,
		"repo":         owner + "/" + name,
		"git_url":      externalURL + "/" + owner + "/" + name + ".git",
		"branch":       branch,
		"pikoci_url":   conf.PikoCI.PublicURL,
		"serial_group": pipeline,
	}

	if err := s.PublishPipeline(ctx, pipeline, hcl, vars); err != nil {
		return nil, fmt.Errorf("publicar pipeline %q: %w", pipeline, err)
	}

	token, err := s.ResourceWebhookToken(ctx, pipeline, "repo")
	if err != nil {
		return nil, fmt.Errorf("obtener webhook token de %q: %w", pipeline, err)
	}

	link := &database.PikoCIRepo{
		RepoID:       repoID,
		PipelineName: pipeline,
		WebhookToken: token,
	}
	if err := s.repoLinks().Upsert(ctx, link); err != nil {
		return nil, fmt.Errorf("guardar vínculo de PikoCI: %w", err)
	}
	return link, nil
}

// DisableRepoPipeline borra el pipeline en PikoCI (si se conoce su nombre) y
// el vínculo local.
func (s *Service) DisableRepoPipeline(ctx context.Context, repoID int64, pipeline string) error {
	if pipeline == "" {
		link, err := s.repoLinks().GetByRepoID(ctx, repoID)
		switch {
		case err == nil:
			pipeline = link.PipelineName
		case database.IsErrPikoCIRepoNotExist(err):
		default:
			return err
		}
	}

	if pipeline != "" {
		if err := s.DeletePipeline(ctx, pipeline); err != nil {
			return err
		}
	}
	return s.repoLinks().Delete(ctx, repoID)
}

// RemoveRepoPipeline intenta borrar el pipeline en PikoCI y elimina siempre el
// vínculo local. Se usa al borrar, renombrar o transferir un repositorio.
func (s *Service) RemoveRepoPipeline(ctx context.Context, repoID int64) error {
	link, err := s.repoLinks().GetByRepoID(ctx, repoID)

	var pipelineErr error
	switch {
	case err == nil && link.PipelineName != "":
		pipelineErr = s.DeletePipeline(ctx, link.PipelineName)
	case err != nil && !database.IsErrPikoCIRepoNotExist(err):
		return err
	}

	if err := database.Handle.PikoCIRepos().Delete(ctx, repoID); err != nil {
		return err
	}
	return pipelineErr
}

// PipelineFileChanged informa si el fichero del pipeline cambió entre dos
// commits. Un error del commit antiguo (alta de rama) o del nuevo (borrado)
// cuentan como cambio: el sync decidirá qué hacer.
func PipelineFileChanged(owner, name, oldSHA, newSHA, path string) (bool, error) {
	if path == "" {
		return false, nil
	}

	gitRepo, err := git.Open(database.RepoPath(owner, name))
	if err != nil {
		return false, fmt.Errorf("abrir repositorio: %w", err)
	}

	oldID, err := fileBlobID(gitRepo, oldSHA, path)
	if err != nil {
		return true, nil
	}
	newID, err := fileBlobID(gitRepo, newSHA, path)
	if err != nil {
		return true, nil
	}
	return oldID != newID, nil
}

func fileBlobID(gitRepo *git.Repository, rev, path string) (string, error) {
	commit, err := gitRepo.CatFileCommit(rev)
	if err != nil {
		return "", err
	}
	entry, err := commit.TreeEntry(path)
	if err != nil {
		return "", err
	}
	if !entry.IsBlob() {
		return "", fmt.Errorf("%s no es un archivo", path)
	}
	return entry.ID().String(), nil
}

func readPipelineFile(owner, name, branch, path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("falta [pikoci] PIPELINE_FILE")
	}

	gitRepo, err := git.Open(database.RepoPath(owner, name))
	if err != nil {
		return nil, fmt.Errorf("abrir repositorio: %w", err)
	}
	commit, err := gitRepo.BranchCommit(branch)
	if err != nil {
		return nil, fmt.Errorf("leer la rama %q: %w", branch, err)
	}
	entry, err := commit.TreeEntry(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s en %s", ErrPipelineFileNotFound, path, branch)
	}
	if !entry.IsBlob() {
		return nil, fmt.Errorf("%s no es un archivo", path)
	}
	return entry.Blob().Bytes()
}
