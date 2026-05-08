package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
	log "unknwon.dev/clog/v2"

	"github.com/gogs/git-module"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/domain"
	"github.com/nef734/ggit/internal/errx"
	"github.com/nef734/ggit/internal/osx"
	"github.com/nef734/ggit/internal/process"
	apiv2types "github.com/nef734/ggit/internal/route/api/v2/types"
	"github.com/nef734/ggit/internal/sync"
)

var PullRequestQueue = sync.NewUniqueQueue(1000)

type PullRequestType int

const (
	PullRequestTypeGogs PullRequestType = iota
	PullRequestTypeGit
)

type PullRequestStatus int

const (
	PullRequestStatusConflict PullRequestStatus = iota
	PullRequestStatusChecking
	PullRequestStatusMergeable
)

type AttentionState int

const (
	AttentionStateNone AttentionState = iota
	AttentionStateReviewRequired
	AttentionStateWaitingForAuthor
	AttentionStateApproved
	AttentionStateClosed
)

// PullRequest represents relation between pull request and repositories.
type PullRequest struct {
	ID             int64 `gorm:"primaryKey"`
	Type           PullRequestType
	Status         PullRequestStatus
	AttentionState AttentionState `gorm:"not null;default:0"`

	IssueID int64  `gorm:"index"`
	Issue   *Issue `json:"-" gorm:"-"`
	Index   int64

	HeadRepoID   int64
	HeadRepo     *Repository `json:"-" gorm:"-"`
	BaseRepoID   int64
	BaseRepo     *Repository `json:"-" gorm:"-"`
	HeadUserName string
	HeadBranch   string
	HeadCommitID string `gorm:"type:VARCHAR(40)"`
	BaseBranch   string
	MergeBase    string `gorm:"type:VARCHAR(40)"`

	HasMerged      bool
	MergedCommitID string `gorm:"type:VARCHAR(40)"`
	MergerID       int64
	Merger         *User     `json:"-" gorm:"-"`
	Merged         time.Time `json:"-" gorm:"-"`
	MergedUnix     int64
}

type PullRequestsStore struct {
	db *gorm.DB
}

func newPullRequestsStore(db *gorm.DB) *PullRequestsStore {
	return &PullRequestsStore{db: db}
}

func (s *PullRequestsStore) GetByID(ctx context.Context, id int64) (*domain.PullRequest, error) {
	pr, err := getPullRequestByID(s.db, id)
	if err != nil {
		return nil, err
	}
	return pr.ToDomain(), nil
}

func (s *PullRequestsStore) GetByIssueID(ctx context.Context, issueID int64) (*domain.PullRequest, error) {
	pr, err := getPullRequestByIssueID(s.db, issueID)
	if err != nil {
		return nil, err
	}
	return pr.ToDomain(), nil
}

func (pr *PullRequest) ToDomain() *domain.PullRequest {
	if pr == nil {
		return nil
	}
	return &domain.PullRequest{
		ID:             pr.ID,
		Type:           int(pr.Type),
		Status:         int(pr.Status),
		AttentionState: int(pr.AttentionState),
		IssueID:        pr.IssueID,
		Issue:          pr.Issue.ToDomain(),
		Index:          pr.Index,
		HeadRepoID:     pr.HeadRepoID,
		HeadRepo:       pr.HeadRepo.ToDomain(),
		BaseRepoID:     pr.BaseRepoID,
		BaseRepo:       pr.BaseRepo.ToDomain(),
		HeadUserName:   pr.HeadUserName,
		HeadBranch:     pr.HeadBranch,
		HeadCommitID:   pr.HeadCommitID,
		BaseBranch:     pr.BaseBranch,
		MergeBase:      pr.MergeBase,
		HasMerged:      pr.HasMerged,
		MergedCommitID: pr.MergedCommitID,
		MergerID:       pr.MergerID,
		Merger:         pr.Merger.ToDomain(),
		Merged:         pr.Merged,
	}
}

func FromDomainPullRequest(pr *domain.PullRequest) *PullRequest {
	if pr == nil {
		return nil
	}
	return &PullRequest{
		ID:             pr.ID,
		Type:           PullRequestType(pr.Type),
		Status:         PullRequestStatus(pr.Status),
		AttentionState: AttentionState(pr.AttentionState),
		IssueID:        pr.IssueID,
		Issue:          FromDomainIssue(pr.Issue),
		Index:          pr.Index,
		HeadRepoID:     pr.HeadRepoID,
		HeadRepo:       FromDomainRepository(pr.HeadRepo),
		BaseRepoID:     pr.BaseRepoID,
		BaseRepo:       FromDomainRepository(pr.BaseRepo),
		HeadUserName:   pr.HeadUserName,
		HeadBranch:     pr.HeadBranch,
		HeadCommitID:   pr.HeadCommitID,
		BaseBranch:     pr.BaseBranch,
		MergeBase:      pr.MergeBase,
		HasMerged:      pr.HasMerged,
		MergedCommitID: pr.MergedCommitID,
		MergerID:       pr.MergerID,
		Merger:         FromDomainUser(pr.Merger),
		Merged:         pr.Merged,
		MergedUnix:     pr.Merged.Unix(),
	}
}

// BeforeUpdate implements the GORM update hook.
func (pr *PullRequest) BeforeUpdate(tx *gorm.DB) error {
	pr.MergedUnix = pr.Merged.Unix()
	return nil
}

// AfterFind implements the GORM query hook.
func (pr *PullRequest) AfterFind(_ *gorm.DB) error {
	if pr.HasMerged {
		pr.Merged = time.Unix(pr.MergedUnix, 0).Local()
	}
	return nil
}

// Note: don't try to get Issue because will end up recursive querying.
func (pr *PullRequest) loadAttributes(tx *gorm.DB) (err error) {
	if pr.HeadRepo == nil {
		pr.HeadRepo, err = getRepositoryByID(tx, pr.HeadRepoID)
		if err != nil && !IsErrRepoNotExist(err) {
			return errors.Newf("get head repository by ID: %v", err)
		}
	}

	if pr.BaseRepo == nil {
		pr.BaseRepo, err = getRepositoryByID(tx, pr.BaseRepoID)
		if err != nil {
			return errors.Newf("get base repository by ID: %v", err)
		}
	}

	if pr.HasMerged && pr.Merger == nil {
		pr.Merger, err = getUserByID(tx, pr.MergerID)
		if IsErrUserNotExist(err) {
			pr.MergerID = -1
			pr.Merger = NewGhostUser()
		} else if err != nil {
			return errors.Newf("get merger by ID: %v", err)
		}
	}

	return nil
}

func (pr *PullRequest) LoadAttributes() error {
	return pr.loadAttributes(Handle.db)
}

func (pr *PullRequest) LoadIssue() (err error) {
	if pr.Issue != nil {
		return nil
	}

	pr.Issue, err = GetIssueByID(pr.IssueID)
	return err
}

// UpdateAttentionState updates the attention state of the pull request.
func (pr *PullRequest) UpdateAttentionState(tx *gorm.DB, state AttentionState) error {
	pr.AttentionState = state
	return tx.Model(pr).Update("attention_state", state).Error
}

// This method assumes following fields have been assigned with valid values:
// Required - Issue, BaseRepo
// Optional - HeadRepo, Merger
func (pr *PullRequest) APIFormat() *apiv2types.PullRequestResponse {
	// In case of head repo has been deleted.
	var apiHeadRepo *apiv2types.RepositoryResponse
	if pr.HeadRepo != nil {
		apiHeadRepo = pr.HeadRepo.APIFormatLegacy(nil)
	}

	if pr.Issue == nil {
		_ = pr.LoadIssue()
	}

	res := &apiv2types.PullRequestResponse{
		ID:         pr.ID,
		Index:      pr.Index,
		Poster:     pr.Issue.Poster.APIFormat(),
		Title:      pr.Issue.Title,
		Body:       pr.Issue.Content,
		State:      pr.Issue.State(),
		Comments:   pr.Issue.NumComments,
		HeadBranch: pr.HeadBranch,
		HeadRepo:   apiHeadRepo,
		BaseBranch: pr.BaseBranch,
		BaseRepo:   pr.BaseRepo.APIFormatLegacy(nil),
		HTMLURL:    pr.Issue.HTMLURL(),
		HasMerged:  pr.HasMerged,
		CreatedAt:  pr.Issue.Created,
		UpdatedAt:  pr.Issue.Updated,
		AttentionState: int(pr.AttentionState),
		CombinedStatus: string(Handle.CommitStatuses().GetCombinedStatusBySHA(context.TODO(), pr.BaseRepoID, pr.HeadCommitID)),
	}

	if pr.Status != PullRequestStatusChecking {
		mergeable := pr.Status != PullRequestStatusConflict
		res.Mergeable = &mergeable
	}
	if pr.HasMerged {
		res.MergedAt = &pr.Merged
		res.MergedCommitID = &pr.MergedCommitID
		res.MergedBy = pr.Merger.APIFormat()
	}

	return res
}

// IsChecking returns true if this pull request is still checking conflict.
func (pr *PullRequest) IsChecking() bool {
	return pr.Status == PullRequestStatusChecking
}

// CanAutoMerge returns true if pull request can be merged automatically (no technical conflicts).
func (pr *PullRequest) CanAutoMerge() bool {
	return pr.Status == PullRequestStatusMergeable
}

// CheckMergeChecks checks if the pull request meets both repository-level and branch protection requirements.
func (pr *PullRequest) CheckMergeChecks(ctx context.Context, doer *User) error {
	if pr.BaseRepo == nil {
		if err := pr.loadAttributes(Handle.db); err != nil {
			return err
		}
	}
	repo := pr.BaseRepo

	// 1. Repository-level checks
	if repo.RequiredApprovals > 0 || repo.BlockMergeOnPendingReviews || repo.BlockMergeOnUnresolvedConversations {
		stats, err := Handle.PullReviews().GetReviewStatsByIssueID(ctx, pr.IssueID)
		if err != nil {
			return err
		}

		if repo.RequiredApprovals > 0 && int(stats.Approvals) < repo.RequiredApprovals {
			return errors.Newf("Insufficient approvals (repo): %d/%d", stats.Approvals, repo.RequiredApprovals)
		}

		if repo.BlockMergeOnPendingReviews && stats.Rejections > 0 {
			return errors.New("All requested changes must be addressed before merging (repo).")
		}

		if repo.BlockMergeOnUnresolvedConversations {
			comments, err := Handle.PullReviewComments().ListByIssueID(ctx, pr.IssueID)
			if err != nil {
				return err
			}
			for _, c := range comments {
				if !c.IsResolved {
					return errors.New("All conversations must be resolved before merging (repo).")
				}
			}
		}
	}

	// 2. Branch-specific checks
	return pr.CheckBranchProtection(ctx, doer)
}

// CalculateAttentionState calculates the current attention state of the pull request.
func (pr *PullRequest) CalculateAttentionState(ctx context.Context, tx *gorm.DB) (AttentionState, error) {
	db := Handle
	if tx != nil {
		db = Handle.WithTransaction(tx)
	}

	if pr.Issue == nil {
		if err := pr.LoadIssue(); err != nil {
			return AttentionStateNone, err
		}
	}
	if pr.Issue.IsClosed {
		return AttentionStateClosed, nil
	}

	stats, err := db.PullReviews().GetReviewStatsByIssueID(ctx, pr.IssueID)
	if err != nil {
		return AttentionStateNone, err
	}

	// 1. If there are rejections, author needs to take action.
	if stats.Rejections > 0 {
		return AttentionStateWaitingForAuthor, nil
	}

	// 2. If there are unresolved conversations, author needs to take action.
	comments, err := db.PullReviewComments().ListByIssueID(ctx, pr.IssueID)
	if err != nil {
		return AttentionStateNone, err
	}
	for _, c := range comments {
		if !c.IsResolved {
			return AttentionStateWaitingForAuthor, nil
		}
	}

	// 3. If everything is resolved and approved (based on repo/branch rules).
	if pr.BaseRepo == nil {
		_ = pr.loadAttributes(db.db)
	}
	
	requiredApprovals := pr.BaseRepo.RequiredApprovals
	// Check branch protection for overrides
	pb, _ := GetProtectBranchOfRepoByName(pr.BaseRepoID, pr.BaseBranch)
	if pb != nil && pb.Protected && pb.RequiredApprovals > requiredApprovals {
		requiredApprovals = pb.RequiredApprovals
	}

	if int(stats.Approvals) >= requiredApprovals && requiredApprovals > 0 {
		return AttentionStateApproved, nil
	}

	// 4. Default to review required.
	return AttentionStateReviewRequired, nil
}

// RefreshAttentionState recalculates and updates the attention state of the pull request.
func (pr *PullRequest) RefreshAttentionState(ctx context.Context, tx *gorm.DB) error {
	state, err := pr.CalculateAttentionState(ctx, tx)
	if err != nil {
		return err
	}
	return pr.UpdateAttentionState(tx, state)
}

// CheckBranchProtection checks if the pull request meets branch protection requirements.
func (pr *PullRequest) CheckBranchProtection(ctx context.Context, doer *User) error {
	if pr.BaseRepo == nil {
		err := pr.loadAttributes(Handle.db)
		if err != nil {
			return err
		}
	}

	pb, err := GetProtectBranchOfRepoByName(pr.BaseRepoID, pr.BaseBranch)
	if err != nil {
		if IsErrBranchNotExist(err) {
			return nil // No protection
		}
		return err
	}

	if !pb.Protected {
		return nil
	}

	// 1. Check for required approvals
	if pb.RequiredApprovals > 0 || pb.BlockOnPendingReviews {
		stats, err := Handle.PullReviews().GetReviewStatsByIssueID(ctx, pr.IssueID)
		if err != nil {
			return err
		}
		if stats.Approvals < int64(pb.RequiredApprovals) {
			return errors.Newf("Insufficient approvals: %d/%d", stats.Approvals, pb.RequiredApprovals)
		}
		if stats.Rejections > 0 && (pb.BlockOnPendingReviews || pb.RequiredApprovals > 0) {
			return errors.New("All requested changes must be addressed before merging.")
		}
	}

	// 2. Check for unresolved conversations
	if pb.BlockOnUnresolvedConversations {
		comments, err := Handle.PullReviewComments().ListByIssueID(ctx, pr.IssueID)
		if err != nil {
			return err
		}
		for _, c := range comments {
			if !c.IsResolved {
				return errors.New("All conversations must be resolved before merging.")
			}
		}
	}

	// 2. Check for merge whitelist
	if pb.EnableMergeWhitelist {
		// This logic usually involves checking if doer.ID is in pb.MergeWhitelistUserIDs or belongs to a team in pb.MergeWhitelistTeamIDs.
		// For now, let's keep it simple or implement a helper.
	}

	return nil
}


// MergeStyle represents the approach to merge commits into base branch.
type MergeStyle string

const (
	MergeStyleRegular MergeStyle = "create_merge_commit"
	MergeStyleRebase  MergeStyle = "rebase_before_merging"
)

// Merge merges pull request to base repository.
// FIXME: add repoWorkingPull make sure two merges does not happen at same time.
func (pr *PullRequest) Merge(doer *User, baseGitRepo *git.Repository, mergeStyle MergeStyle, commitDescription string) (err error) {
	ctx := context.TODO()

	defer func() {
		go HookQueue.Add(pr.BaseRepo.ID)
		go AddTestPullRequestTask(doer, pr.BaseRepo.ID, pr.BaseBranch, false)
	}()

	var headGitRepo *git.Repository
	err = Handle.db.Transaction(func(tx *gorm.DB) error {
		if err = pr.Issue.changeStatus(tx, doer, pr.Issue.Repo, true); err != nil {
			return errors.Newf("Issue.changeStatus: %v", err)
		}

		headRepoPath := RepoPath(pr.HeadUserName, pr.HeadRepo.Name)
		headGitRepo, err = git.Open(headRepoPath)
		if err != nil {
			return errors.Newf("open repository: %v", err)
		}

		// Create temporary directory to store temporary copy of the base repository,
		// and clean it up when operation finished regardless of succeed or not.
		tmpBasePath := filepath.Join(conf.Server.AppDataPath, "tmp", "repos", strconv.Itoa(time.Now().Nanosecond())+".git")
		if err = os.MkdirAll(filepath.Dir(tmpBasePath), os.ModePerm); err != nil {
			return err
		}
		defer func() {
			_ = os.RemoveAll(filepath.Dir(tmpBasePath))
		}()

		// Clone the base repository to the defined temporary directory,
		// and checks out to base branch directly.
		var stderr string
		if _, stderr, err = process.ExecTimeout(5*time.Minute,
			fmt.Sprintf("PullRequest.Merge (git clone): %s", tmpBasePath),
			"git", "clone", "-b", pr.BaseBranch, baseGitRepo.Path(), tmpBasePath); err != nil {
			return errors.Newf("git clone: %s", stderr)
		}

		// Add remote which points to the head repository.
		if _, stderr, err = process.ExecDir(-1, tmpBasePath,
			fmt.Sprintf("PullRequest.Merge (git remote add): %s", tmpBasePath),
			"git", "remote", "add", "head_repo", headRepoPath); err != nil {
			return errors.Newf("git remote add [%s -> %s]: %s", headRepoPath, tmpBasePath, stderr)
		}

		// Fetch information from head repository to the temporary copy.
		if _, stderr, err = process.ExecDir(-1, tmpBasePath,
			fmt.Sprintf("PullRequest.Merge (git fetch): %s", tmpBasePath),
			"git", "fetch", "head_repo"); err != nil {
			return errors.Newf("git fetch [%s -> %s]: %s", headRepoPath, tmpBasePath, stderr)
		}

		remoteHeadBranch := "head_repo/" + pr.HeadBranch

		// Check if merge style is allowed, reset to default style if not
		if mergeStyle == MergeStyleRebase && !pr.BaseRepo.PullsAllowRebase {
			mergeStyle = MergeStyleRegular
		}

		switch mergeStyle {
		case MergeStyleRegular: // Create merge commit

			// Merge changes from head branch.
			if _, stderr, err = process.ExecDir(-1, tmpBasePath,
				fmt.Sprintf("PullRequest.Merge (git merge --no-ff --no-commit): %s", tmpBasePath),
				"git", "merge", "--no-ff", "--no-commit", remoteHeadBranch); err != nil {
				return errors.Newf("git merge --no-ff --no-commit [%s]: %v - %s", tmpBasePath, err, stderr)
			}

			// Create a merge commit for the base branch.
			if _, stderr, err = process.ExecDir(-1, tmpBasePath,
				fmt.Sprintf("PullRequest.Merge (git merge): %s", tmpBasePath),
				"git", "commit", fmt.Sprintf("--author='%s <%s>'", doer.DisplayName(), doer.Email),
				"-m", fmt.Sprintf("Merge branch '%s' of %s/%s into %s", pr.HeadBranch, pr.HeadUserName, pr.HeadRepo.Name, pr.BaseBranch),
				"-m", commitDescription); err != nil {
				return errors.Newf("git commit [%s]: %v - %s", tmpBasePath, err, stderr)
			}

		case MergeStyleRebase: // Rebase before merging

			// Rebase head branch based on base branch, this creates a non-branch commit state.
			if _, stderr, err = process.ExecDir(-1, tmpBasePath,
				fmt.Sprintf("PullRequest.Merge (git rebase): %s", tmpBasePath),
				"git", "rebase", "--quiet", pr.BaseBranch, remoteHeadBranch); err != nil {
				return errors.Newf("git rebase [%s on %s]: %s", remoteHeadBranch, pr.BaseBranch, stderr)
			}

			// Name non-branch commit state to a new temporary branch in order to save changes.
			tmpBranch := strconv.FormatInt(time.Now().UnixNano(), 10)
			if _, stderr, err = process.ExecDir(-1, tmpBasePath,
				fmt.Sprintf("PullRequest.Merge (git checkout): %s", tmpBasePath),
				"git", "checkout", "-b", tmpBranch); err != nil {
				return errors.Newf("git checkout '%s': %s", tmpBranch, stderr)
			}

			// Check out the base branch to be operated on.
			if _, stderr, err = process.ExecDir(-1, tmpBasePath,
				fmt.Sprintf("PullRequest.Merge (git checkout): %s", tmpBasePath),
				"git", "checkout", pr.BaseBranch); err != nil {
				return errors.Newf("git checkout '%s': %s", pr.BaseBranch, stderr)
			}

			// Merge changes from temporary branch to the base branch.
			if _, stderr, err = process.ExecDir(-1, tmpBasePath,
				fmt.Sprintf("PullRequest.Merge (git merge): %s", tmpBasePath),
				"git", "merge", tmpBranch); err != nil {
				return errors.Newf("git merge [%s]: %v - %s", tmpBasePath, err, stderr)
			}

		default:
			return errors.Newf("unknown merge style: %s", mergeStyle)
		}

		// Push changes on base branch to upstream.
		if _, stderr, err = process.ExecDir(-1, tmpBasePath,
			fmt.Sprintf("PullRequest.Merge (git push): %s", tmpBasePath),
			"git", "push", baseGitRepo.Path(), pr.BaseBranch); err != nil {
			return errors.Newf("git push: %s", stderr)
		}

		pr.MergedCommitID, err = headGitRepo.BranchCommitID(pr.HeadBranch)
		if err != nil {
			return errors.Newf("get head branch %q commit ID: %v", pr.HeadBranch, err)
		}

		pr.HasMerged = true
		pr.Merged = time.Now()
		pr.MergerID = doer.ID
		if err = tx.Model(pr).Updates(pr).Error; err != nil {
			return errors.Newf("update pull request: %v", err)
		}

		return nil
	})
	if err != nil {
		return err
	}
	if err = Handle.Actions().MergePullRequest(ctx, doer.ToDomain(), pr.Issue.Repo.Owner.ToDomain(), pr.Issue.Repo.ToDomain(), pr.Issue.ToDomain()); err != nil {
		log.Error("Failed to create action for merge pull request, pull_request_id: %d, error: %v", pr.ID, err)
	}

	// Reload pull request information.
	if err = pr.LoadAttributes(); err != nil {
		log.Error("LoadAttributes: %v", err)
		return nil
	}
	if err = PrepareWebhooks(pr.Issue.Repo, HookEventTypePullRequest, &apiv2types.WebhookPullRequestPayload{
		Action:      apiv2types.WebhookIssueClosed,
		Index:       pr.Index,
		PullRequest: pr.APIFormat(),
		Repository:  pr.Issue.Repo.APIFormatLegacy(nil),
		Sender:      doer.APIFormat(),
	}); err != nil {
		log.Error("PrepareWebhooks: %v", err)
		return nil
	}

	commits, err := headGitRepo.RevList([]string{pr.MergeBase + "..." + pr.MergedCommitID})
	if err != nil {
		log.Error("Failed to list commits [merge_base: %s, merged_commit_id: %s]: %v", pr.MergeBase, pr.MergedCommitID, err)
		return nil
	}

	// NOTE: It is possible that head branch is not fully sync with base branch
	// for merge commits, so we need to get latest head commit and append merge
	// commit manually to avoid strange diff commits produced.
	mergeCommit, err := baseGitRepo.BranchCommit(pr.BaseBranch)
	if err != nil {
		log.Error("Failed to get base branch %q commit: %v", pr.BaseBranch, err)
		return nil
	}
	if mergeStyle == MergeStyleRegular {
		commits = append([]*git.Commit{mergeCommit}, commits...)
	}

	pcs, err := CommitsToPushCommits(commits).APIFormat(ctx, Handle.Users(), pr.BaseRepo.RepoPath(), pr.BaseRepo.HTMLURL())
	if err != nil {
		log.Error("Failed to convert to API payload commits: %v", err)
		return nil
	}

	p := &apiv2types.WebhookPushPayload{
		Ref:        git.RefsHeads + pr.BaseBranch,
		Before:     pr.MergeBase,
		After:      mergeCommit.ID.String(),
		CompareURL: conf.Server.ExternalURL + pr.BaseRepo.ComposeCompareURL(pr.MergeBase, pr.MergedCommitID),
		Commits:    pcs,
		Repo:       pr.BaseRepo.APIFormatLegacy(nil),
		Pusher:     pr.HeadRepo.MustOwner().APIFormat(),
		Sender:     doer.APIFormat(),
	}
	if err = PrepareWebhooks(pr.BaseRepo, HookEventTypePush, p); err != nil {
		log.Error("Failed to prepare webhooks: %v", err)
		return nil
	}
	return nil
}

// testPatch checks if patch can be merged to base repository without conflict.
// FIXME: make a mechanism to clean up stable local copies.
func (pr *PullRequest) testPatch() (err error) {
	if pr.BaseRepo == nil {
		pr.BaseRepo, err = GetRepositoryByID(pr.BaseRepoID)
		if err != nil {
			return errors.Newf("GetRepositoryByID: %v", err)
		}
	}

	patchPath, err := pr.BaseRepo.PatchPath(pr.Index)
	if err != nil {
		return errors.Newf("BaseRepo.PatchPath: %v", err)
	}

	// Fast fail if patch does not exist, this assumes data is corrupted.
	if !osx.IsFile(patchPath) {
		log.Trace("PullRequest[%d].testPatch: ignored corrupted data", pr.ID)
		return nil
	}

	repoWorkingPool.CheckIn(strconv.FormatInt(pr.BaseRepoID, 10))
	defer repoWorkingPool.CheckOut(strconv.FormatInt(pr.BaseRepoID, 10))

	log.Trace("PullRequest[%d].testPatch (patchPath): %s", pr.ID, patchPath)

	if err := pr.BaseRepo.UpdateLocalCopyBranch(pr.BaseBranch); err != nil {
		return errors.Newf("UpdateLocalCopy [%d]: %v", pr.BaseRepoID, err)
	}

	args := []string{"apply", "--check"}
	if pr.BaseRepo.PullsIgnoreWhitespace {
		args = append(args, "--ignore-whitespace")
	}
	args = append(args, patchPath)

	pr.Status = PullRequestStatusChecking
	_, stderr, err := process.ExecDir(-1, pr.BaseRepo.LocalCopyPath(),
		fmt.Sprintf("testPatch (git apply --check): %d", pr.BaseRepo.ID),
		"git", args...)
	if err != nil {
		log.Trace("PullRequest[%d].testPatch (apply): has conflict\n%s", pr.ID, stderr)
		pr.Status = PullRequestStatusConflict
		return nil
	}
	return nil
}

// NewPullRequest creates new pull request with labels for repository.
func NewPullRequest(repo *Repository, pull *Issue, labelIDs []int64, uuids []string, pr *PullRequest, patch []byte) (err error) {
	return Handle.db.Transaction(func(tx *gorm.DB) error {
		if err = newIssue(tx, NewIssueOptions{
			Repo:        repo,
			Issue:       pull,
			LableIDs:    labelIDs,
			Attachments: uuids,
			IsPull:      true,
		}); err != nil {
			return errors.Newf("newIssue: %v", err)
		}

		pr.Index = pull.Index
		if err = repo.SavePatch(pr.Index, patch); err != nil {
			return errors.Newf("SavePatch: %v", err)
		}

		pr.BaseRepo = repo
		if err = pr.testPatch(); err != nil {
			return errors.Newf("testPatch: %v", err)
		}
		// No conflict appears after test means mergeable.
		if pr.Status == PullRequestStatusChecking {
			pr.Status = PullRequestStatusMergeable
		}

		pr.IssueID = pull.ID
		pr.AttentionState = AttentionStateReviewRequired
		
		headGitRepo, err := git.Open(pr.HeadRepo.RepoPath())
		if err == nil {
			pr.HeadCommitID, _ = headGitRepo.BranchCommitID(pr.HeadBranch)
		}

		return tx.Create(pr).Error
	})
}

// GetUnmergedPullRequest returns a pull request that is open and has not been merged
// by given head/base and repo/branch.
func GetUnmergedPullRequest(headRepoID, baseRepoID int64, headBranch, baseBranch string) (*PullRequest, error) {
	pr := new(PullRequest)
	err := Handle.db.Joins("JOIN issue ON issue.id = pull_request.issue_id").
		Where("head_repo_id = ? AND head_branch = ? AND base_repo_id = ? AND base_branch = ? AND has_merged = ? AND issue.is_closed = ?",
			headRepoID, headBranch, baseRepoID, baseBranch, false, false).
		First(pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPullRequestNotExist{args: map[string]any{
				"headRepoID": headRepoID,
				"baseRepoID": baseRepoID,
				"headBranch": headBranch,
				"baseBranch": baseBranch,
			}}
		}
		return nil, err
	}

	return pr, nil
}

// GetUnmergedPullRequestsByHeadInfo returns all pull requests that are open and has not been merged
// by given head information (repo and branch).
func GetUnmergedPullRequestsByHeadInfo(repoID int64, branch string) ([]*PullRequest, error) {
	prs := make([]*PullRequest, 0, 2)
	err := Handle.db.Joins("JOIN issue ON issue.id = pull_request.issue_id").
		Where("head_repo_id = ? AND head_branch = ? AND has_merged = ? AND issue.is_closed = ?",
			repoID, branch, false, false).
		Find(&prs).Error
	return prs, err
}

// GetUnmergedPullRequestsByBaseInfo returns all pull requests that are open and has not been merged
// by given base information (repo and branch).
func GetUnmergedPullRequestsByBaseInfo(repoID int64, branch string) ([]*PullRequest, error) {
	prs := make([]*PullRequest, 0, 2)
	err := Handle.db.Joins("JOIN issue ON issue.id = pull_request.issue_id").
		Where("base_repo_id = ? AND base_branch = ? AND has_merged = ? AND issue.is_closed = ?",
			repoID, branch, false, false).
		Find(&prs).Error
	return prs, err
}

var _ errx.NotFound = (*ErrPullRequestNotExist)(nil)

type ErrPullRequestNotExist struct {
	args map[string]any
}

func IsErrPullRequestNotExist(err error) bool {
	_, ok := err.(ErrPullRequestNotExist)
	return ok
}

func (err ErrPullRequestNotExist) Error() string {
	return fmt.Sprintf("pull request does not exist: %v", err.args)
}

func (ErrPullRequestNotExist) NotFound() bool {
	return true
}

func getPullRequestByID(tx *gorm.DB, id int64) (*PullRequest, error) {
	pr := new(PullRequest)
	err := tx.Where("id = ?", id).First(pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPullRequestNotExist{args: map[string]any{"pullRequestID": id}}
		}
		return nil, err
	}
	return pr, pr.loadAttributes(tx)
}

// GetPullRequestByID returns a pull request by given ID.
func GetPullRequestByID(id int64) (*PullRequest, error) {
	return getPullRequestByID(Handle.db, id)
}

func getPullRequestByIssueID(tx *gorm.DB, issueID int64) (*PullRequest, error) {
	pr := new(PullRequest)
	err := tx.Where("issue_id = ?", issueID).First(pr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPullRequestNotExist{args: map[string]any{"issueID": issueID}}
		}
		return nil, err
	}
	return pr, pr.loadAttributes(tx)
}

// GetPullRequestByIssueID returns pull request by given issue ID.
func GetPullRequestByIssueID(issueID int64) (*PullRequest, error) {
	return getPullRequestByIssueID(Handle.db, issueID)
}

// Update updates all fields of pull request.
func (pr *PullRequest) Update() error {
	return Handle.db.Model(pr).Updates(pr).Error
}

// Update updates specific fields of pull request.
func (pr *PullRequest) UpdateCols(cols ...string) error {
	return Handle.db.Model(pr).Select(cols).Updates(pr).Error
}

// UpdatePatch generates and saves a new patch.
func (pr *PullRequest) UpdatePatch() (err error) {
	headGitRepo, err := git.Open(pr.HeadRepo.RepoPath())
	if err != nil {
		return errors.Newf("open repository: %v", err)
	}

	// Add a temporary remote.
	tmpRemote := strconv.FormatInt(time.Now().UnixNano(), 10)
	baseRepoPath := RepoPath(pr.BaseRepo.MustOwner().Name, pr.BaseRepo.Name)
	err = headGitRepo.RemoteAdd(tmpRemote, baseRepoPath, git.RemoteAddOptions{Fetch: true})
	if err != nil {
		return errors.Newf("add remote %q [repo_id: %d]: %v", tmpRemote, pr.HeadRepoID, err)
	}
	defer func() {
		if err := headGitRepo.RemoteRemove(tmpRemote); err != nil {
			log.Error("Failed to remove remote %q [repo_id: %d]: %v", tmpRemote, pr.HeadRepoID, err)
		}
	}()

	remoteBranch := "remotes/" + tmpRemote + "/" + pr.BaseBranch
	pr.MergeBase, err = headGitRepo.MergeBase(remoteBranch, pr.HeadBranch)
	if err != nil {
		return errors.Newf("get merge base: %v", err)
	}
	pr.HeadCommitID, err = headGitRepo.BranchCommitID(pr.HeadBranch)
	if err != nil {
		return errors.Newf("get head commit ID: %v", err)
	}
	if err = pr.Update(); err != nil {
		return errors.Newf("update: %v", err)
	}

	patch, err := headGitRepo.DiffBinary(pr.MergeBase, pr.HeadBranch)
	if err != nil {
		return errors.Newf("get binary patch: %v", err)
	}

	if err = pr.BaseRepo.SavePatch(pr.Index, patch); err != nil {
		return errors.Newf("save patch: %v", err)
	}

	log.Trace("PullRequest[%d].UpdatePatch: patch saved", pr.ID)
	return nil
}

// PushToBaseRepo pushes commits from branches of head repository to
// corresponding branches of base repository.
// FIXME: Only push branches that are actually updates?
func (pr *PullRequest) PushToBaseRepo() (err error) {
	log.Trace("PushToBaseRepo[%d]: pushing commits to base repo 'refs/pull/%d/head'", pr.BaseRepoID, pr.Index)

	headRepoPath := pr.HeadRepo.RepoPath()
	headGitRepo, err := git.Open(headRepoPath)
	if err != nil {
		return errors.Newf("open repository: %v", err)
	}

	tmpRemote := fmt.Sprintf("tmp-pull-%d", pr.ID)
	if err = headGitRepo.RemoteAdd(tmpRemote, pr.BaseRepo.RepoPath()); err != nil {
		return errors.Newf("add remote %q [repo_id: %d]: %v", tmpRemote, pr.HeadRepoID, err)
	}

	// Make sure to remove the remote even if the push fails
	defer func() {
		if err := headGitRepo.RemoteRemove(tmpRemote); err != nil {
			log.Error("Failed to remove remote %q [repo_id: %d]: %v", tmpRemote, pr.HeadRepoID, err)
		}
	}()

	headRefspec := fmt.Sprintf("refs/pull/%d/head", pr.Index)
	headFile := filepath.Join(pr.BaseRepo.RepoPath(), headRefspec)
	if osx.Exist(headFile) {
		err = os.Remove(headFile)
		if err != nil {
			return errors.Newf("remove head file [repo_id: %d]: %v", pr.BaseRepoID, err)
		}
	}

	err = headGitRepo.Push(tmpRemote, fmt.Sprintf("%s:%s", pr.HeadBranch, headRefspec))
	if err != nil {
		return errors.Newf("push: %v", err)
	}

	return nil
}

// AddToTaskQueue adds itself to pull request test task queue.
func (pr *PullRequest) AddToTaskQueue() {
	go PullRequestQueue.AddFunc(pr.ID, func() {
		pr.Status = PullRequestStatusChecking
		if err := pr.UpdateCols("status"); err != nil {
			log.Error("AddToTaskQueue.UpdateCols[%d].(add to queue): %v", pr.ID, err)
		}
	})
}

type PullRequestList []*PullRequest

func (prs PullRequestList) loadAttributes(tx *gorm.DB) (err error) {
	if len(prs) == 0 {
		return nil
	}

	// Load issues
	set := make(map[int64]*Issue)
	for i := range prs {
		set[prs[i].IssueID] = nil
	}
	issueIDs := make([]int64, 0, len(prs))
	for issueID := range set {
		issueIDs = append(issueIDs, issueID)
	}
	issues := make([]*Issue, 0, len(issueIDs))
	if err = tx.Where("id IN (?)", issueIDs).Find(&issues).Error; err != nil {
		return errors.Newf("find issues: %v", err)
	}
	for i := range issues {
		set[issues[i].ID] = issues[i]
	}
	for i := range prs {
		prs[i].Issue = set[prs[i].IssueID]
	}

	// Load attributes
	for i := range prs {
		if err = prs[i].loadAttributes(tx); err != nil {
			return errors.Newf("loadAttributes [%d]: %v", prs[i].ID, err)
		}
	}

	return nil
}

func (prs PullRequestList) LoadAttributes() error {
	return prs.loadAttributes(Handle.db)
}

func addHeadRepoTasks(prs []*PullRequest) {
	for _, pr := range prs {
		if pr.HeadRepo == nil {
			log.Trace("addHeadRepoTasks[%d]: missing head repository", pr.ID)
			continue
		}

		log.Trace("addHeadRepoTasks[%d]: composing new test task", pr.ID)
		if err := pr.UpdatePatch(); err != nil {
			log.Error("UpdatePatch: %v", err)
			continue
		} else if err := pr.PushToBaseRepo(); err != nil {
			log.Error("PushToBaseRepo: %v", err)
			continue
		}

		pr.AddToTaskQueue()
		_ = pr.RefreshAttentionState(context.TODO(), Handle.db)
	}
}

// AddTestPullRequestTask adds new test tasks by given head/base repository and head/base branch,
// and generate new patch for testing as needed.
func AddTestPullRequestTask(doer *User, repoID int64, branch string, isSync bool) {
	log.Trace("AddTestPullRequestTask [head_repo_id: %d, head_branch: %s]: finding pull requests", repoID, branch)
	prs, err := GetUnmergedPullRequestsByHeadInfo(repoID, branch)
	if err != nil {
		log.Error("Find pull requests [head_repo_id: %d, head_branch: %s]: %v", repoID, branch, err)
		return
	}

	if isSync {
		if err = PullRequestList(prs).LoadAttributes(); err != nil {
			log.Error("PullRequestList.LoadAttributes: %v", err)
		}

		if err == nil {
			for _, pr := range prs {
				pr.Issue.PullRequest = pr
				if err = pr.Issue.LoadAttributes(); err != nil {
					log.Error("LoadAttributes: %v", err)
					continue
				}
				if err = PrepareWebhooks(pr.Issue.Repo, HookEventTypePullRequest, &apiv2types.WebhookPullRequestPayload{
					Action:      apiv2types.WebhookIssueSynchronized,
					Index:       pr.Issue.Index,
					PullRequest: pr.Issue.PullRequest.APIFormat(),
					Repository:  pr.Issue.Repo.APIFormatLegacy(nil),
					Sender:      doer.APIFormat(),
				}); err != nil {
					log.Error("PrepareWebhooks [pull_id: %v]: %v", pr.ID, err)
					continue
				}
			}
		}
	}

	addHeadRepoTasks(prs)

	log.Trace("AddTestPullRequestTask [base_repo_id: %d, base_branch: %s]: finding pull requests", repoID, branch)
	prs, err = GetUnmergedPullRequestsByBaseInfo(repoID, branch)
	if err != nil {
		log.Error("Find pull requests [base_repo_id: %d, base_branch: %s]: %v", repoID, branch, err)
		return
	}
	for _, pr := range prs {
		pr.AddToTaskQueue()
	}
}

// checkAndUpdateStatus checks if pull request is possible to leaving checking status,
// and set to be either conflict or mergeable.
func (pr *PullRequest) checkAndUpdateStatus() {
	// Status is not changed to conflict means mergeable.
	if pr.Status == PullRequestStatusChecking {
		pr.Status = PullRequestStatusMergeable
	}

	// Make sure there is no waiting test to process before leaving the checking status.
	if !PullRequestQueue.Exist(pr.ID) {
		if err := pr.UpdateCols("status"); err != nil {
			log.Error("Update[%d]: %v", pr.ID, err)
		}
	}
}

// TestPullRequests checks and tests untested patches of pull requests.
// TODO: test more pull requests at same time.
func TestPullRequests() {
	var prs []*PullRequest
	err := Handle.db.Where("status = ?", PullRequestStatusChecking).Find(&prs).Error
	if err != nil {
		log.Error("Find checking pull requests: %v", err)
		return
	}

	validPRs := make([]*PullRequest, 0, len(prs))
	for _, pr := range prs {
		if err := pr.LoadAttributes(); err != nil {
			log.Error("LoadAttributes [%d]: %v", pr.ID, err)
			continue
		}

		if err := pr.testPatch(); err != nil {
			log.Error("testPatch [%d]: %v", pr.ID, err)
			continue
		}
		validPRs = append(validPRs, pr)
	}

	// Update pull request status.
	for _, pr := range validPRs {
		pr.checkAndUpdateStatus()
	}

	// Start listening on new test requests.
	for prID := range PullRequestQueue.Queue() {
		log.Trace("TestPullRequests[%v]: processing test task", prID)
		PullRequestQueue.Remove(prID)

		id, _ := strconv.ParseInt(prID, 10, 64)
		pr, err := GetPullRequestByID(id)
		if err != nil {
			log.Error("GetPullRequestByID[%s]: %v", prID, err)
			continue
		} else if err = pr.testPatch(); err != nil {
			log.Error("testPatch[%d]: %v", pr.ID, err)
			continue
		}

		pr.checkAndUpdateStatus()
	}
}

func InitTestPullRequests() {
	go TestPullRequests()
}
