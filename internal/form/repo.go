package form

import (
	"net/url"
	"strings"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/netx"
	"github.com/nef734/ggit/internal/osx"
)

// _______________________________________    _________.______________________ _______________.___.
// \______   \_   _____/\______   \_____  \  /   _____/|   \__    ___/\_____  \\______   \__  |   |
//  |       _/|    __)_  |     ___//   |   \ \_____  \ |   | |    |    /   |   \|       _//   |   |
//  |    |   \|        \ |    |   /    |    \/        \|   | |    |   /    |    \    |   \\____   |
//  |____|_  /_______  / |____|   \_______  /_______  /|___| |____|   \_______  /____|_  // ______|
//         \/        \/                   \/        \/                        \/       \/ \/

type CreateRepo struct {
	UserID      int64  `valid:"Required"`
	RepoName    string `valid:"Required;AlphaDashDot;MaxSize(100)"`
	Private     bool
	Unlisted    bool
	Description string `valid:"MaxSize(512)"`
	AutoInit    bool
	Gitignores  string
	License     string
	Readme      string
}

func (f *CreateRepo) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type MigrateRepo struct {
	CloneAddr    string `json:"clone_addr" valid:"Required"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	UID          int64  `json:"uid" valid:"Required"`
	RepoName     string `json:"repo_name" valid:"Required;AlphaDashDot;MaxSize(100)"`
	Private      bool   `json:"private"`
	Unlisted     bool   `json:"unlisted"`
	Description  string `json:"description" valid:"MaxSize(512)"`
}

func (f *MigrateRepo) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

// ParseRemoteAddr checks if given remote address is valid,
// and returns composed URL with needed username and password.
// It also checks if given user has permission when remote address
// is actually a local path.
func (f MigrateRepo) ParseRemoteAddr(user *database.User) (string, error) {
	remoteAddr := strings.TrimSpace(f.CloneAddr)

	// Remote address can be HTTP/HTTPS/Git URL or local path.
	if strings.HasPrefix(remoteAddr, "http://") ||
		strings.HasPrefix(remoteAddr, "https://") ||
		strings.HasPrefix(remoteAddr, "git://") {
		u, err := url.Parse(remoteAddr)
		if err != nil {
			return "", database.ErrInvalidCloneAddr{IsURLError: true}
		}

		if netx.IsBlockedLocalHostname(u.Hostname(), conf.Security.LocalNetworkAllowlist) {
			return "", database.ErrInvalidCloneAddr{IsBlockedLocalAddress: true}
		}

		if len(f.AuthUsername)+len(f.AuthPassword) > 0 {
			u.User = url.UserPassword(f.AuthUsername, f.AuthPassword)
		}
		// To prevent CRLF injection in git protocol, see https://github.com/gogs/gogs/issues/6413
		if u.Scheme == "git" && (strings.Contains(remoteAddr, "%0d") || strings.Contains(remoteAddr, "%0a")) {
			return "", database.ErrInvalidCloneAddr{IsURLError: true}
		}
		remoteAddr = u.String()
	} else if !user.CanImportLocal() {
		return "", database.ErrInvalidCloneAddr{IsPermissionDenied: true}
	} else if !osx.IsDir(remoteAddr) {
		return "", database.ErrInvalidCloneAddr{IsInvalidPath: true}
	}

	return remoteAddr, nil
}

type RepoSetting struct {
	RepoName      string `valid:"Required;AlphaDashDot;MaxSize(100)"`
	Description   string `valid:"MaxSize(512)"`
	Website       string `valid:"Url;MaxSize(100)"`
	Branch        string
	Interval      int
	Private       bool
	Unlisted      bool
	EnablePrune   bool

	// Advanced settings
	EnableWiki            bool
	AllowPublicWiki       bool
	EnableExternalWiki    bool
	ExternalWikiURL       string
	EnableIssues          bool
	AllowPublicIssues     bool
	EnableExternalTracker bool
	ExternalTrackerURL    string
	TrackerURLFormat      string
	TrackerIssueStyle     string
	EnablePulls           bool
	PullsIgnoreWhitespace bool
	PullsAllowRebase      bool
}

func (f *RepoSetting) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

// __________                             .__
// \______   \____________    ____   ____ |  |__
//  |    |  _/\_  __ \__  \  /    \_/ ___\|  |  \
//  |    |   \ |  | \// __ \|   |  \  \___|   Y  \
//  |______  / |__|  (____  /___|  /\___  >___|  /
//         \/             \/     \/     \/     \/

type ProtectBranch struct {
	Protected          bool
	RequirePullRequest bool
	EnableWhitelist    bool
	WhitelistUsers     string
	WhitelistTeams     string
}

func (f *ProtectBranch) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

//  __      __      ___.   .__    .__            __
// /  \    /  \ ____\_ |__ |  |__ |  |__   ____ |  | __
// \   \/\/   // __ \| __ \|  |  \|  |  \ /  _ \|  |/ /
//  \        /\  ___/| \_\ \   Y  \   Y  (  <_> )    <
//   \__/\  /  \___  >___  /___|  /___|  /\____/|__|_ \
//        \/       \/    \/     \/     \/            \/

type Webhook struct {
	Events       string
	Create       bool
	Delete       bool
	Fork         bool
	Push         bool
	Issues       bool
	IssueComment bool
	PullRequest  bool
	Release      bool
	Active       bool
}

func (f Webhook) PushOnly() bool {
	return f.Events == "push_only"
}

func (f Webhook) SendEverything() bool {
	return f.Events == "send_everything"
}

func (f Webhook) ChooseEvents() bool {
	return f.Events == "choose_events"
}

type NewWebhook struct {
	PayloadURL  string `valid:"Required;Url"`
	ContentType int    `valid:"Required"`
	Secret      string
	Webhook
}

func (f *NewWebhook) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

// .___
// |   | ______ ________ __   ____
// |   |/  ___//  ___/  |  \_/ __ \
// |   |\___ \ \___ \|  |  /\  ___/
// |___/____  >____  >____/  \___  >
//          \/     \/            \/

type NewIssue struct {
	Title       string `valid:"Required;MaxSize(255)"`
	LabelIDs    string `form:"label_ids"`
	MilestoneID int64
	AssigneeID  int64
	Content     string
	Files       []string
}

func (f *NewIssue) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type CreateComment struct {
	Content string
	Status  string `valid:"OmitEmpty;In(reopen,close)"`
	Files   []string
}

func (f *CreateComment) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

//    _____  .__.__                   __
//   /     \ |__|  |   ____   _______/  |_  ____   ____   ____
//  /  \ /  \|  |  | _/ __ \ /  ___/\   __\/  _ \ /    \_/ __ \
// /    Y    \  |  |_\  ___/ \___ \  |  | (  <_> )   |  \  ___/
// \____|__  /__|____/\___  >____  > |__|  \____/|___|  /\___  >
//         \/             \/     \/                   \/     \/

type CreateMilestone struct {
	Title    string `valid:"Required;MaxSize(50)"`
	Content  string
	Deadline string
}

func (f *CreateMilestone) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

// .____          ___.          .__
// |    |   _____ \_ |__   ____ |  |
// |    |   \__  \ | __ \_/ __ \|  |
// |    |___ / __ \| \_\ \  ___/|  |__
// |_______ (____  /___  /\___  >____/
//         \/    \/    \/     \/

type CreateLabel struct {
	ID    int64
	Title string `valid:"Required;MaxSize(50)" locale:"repo.issues.label_title"`
	Color string `valid:"Required;Size(7)" locale:"repo.issues.label_color"`
}

func (f *CreateLabel) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type InitializeLabels struct {
	TemplateName string `valid:"Required"`
}

func (f *InitializeLabels) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

// __________       .__
// \______   \ ____ |  |   ____ _____    ______ ____
//  |       _// __ \|  | _/ __ \\__  \  /  ___// __ \
//  |    |   \  ___/|  |_\  ___/ / __ \_\___ \\  ___/
//  |____|_  /\___  >____/\___  >____  /____  >\___  >
//         \/     \/          \/     \/     \/     \/

type NewRelease struct {
	TagName    string `valid:"Required"`
	Target     string `form:"tag_target" valid:"Required"`
	Title      string `valid:"Required"`
	Content    string
	Draft      string
	Prerelease bool
	Files      []string
	CreatedAt  string `json:"created_at"`
}

func (f *NewRelease) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type EditRelease struct {
	Title      string `valid:"Required"`
	Content    string
	Draft      string
	Prerelease bool
	Files      []string
}

func (f *EditRelease) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

//  __      __.__ __   .__
// /  \    /  \__|  | _|__|
// \   \/\/   /  |  |/ /  |
//  \        /|  |    <|  |
//   \__/\  / |__|__|_ \__|
//        \/          \/

type NewWiki struct {
	OldTitle string
	Title    string `valid:"Required"`
	Content  string `valid:"Required"`
	Message  string
	CreatedAt string `json:"created_at"`
}

// FIXME: use code generation to generate this method.
func (f *NewWiki) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

// ___________    .___.__  __
// \_   _____/  __| _/|__|/  |_
//  |    __)_  / __ | |  \   __\
//  |        \/ /_/ | |  ||  |
// /_______  /\____ | |__||__|
//         \/      \/

type EditRepoFile struct {
	TreePath      string `valid:"Required;MaxSize(500)"`
	Content       string `valid:"Required"`
	CommitSummary string `valid:"MaxSize(100)"`
	CommitMessage string
	CommitChoice  string `valid:"Required;MaxSize(50)"`
	NewBranchName string `valid:"AlphaDashDotSlash;MaxSize(100)"`
	LastCommit    string
	CreatedAt     string `json:"created_at"`
}

func (f *EditRepoFile) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

func (f *EditRepoFile) IsNewBrnach() bool {
	return f.CommitChoice == "commit-to-new-branch"
}

type EditPreviewDiff struct {
	Content string
}

func (f *EditPreviewDiff) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

//  ____ ___        .__                    .___
// |    |   \______ |  |   _________     __| _/
// |    |   /\____ \|  |  /  _ \__  \   / __ |
// |    |  / |  |_> >  |_(  <_> ) __ \_/ /_/ |
// |______/  |   __/|____/\____(____  /\____ |
//           |__|                   \/      \/
//

type UploadRepoFile struct {
	TreePath      string `valid:"MaxSize(500)"`
	CommitSummary string `valid:"MaxSize(100)"`
	CommitMessage string
	CommitChoice  string `valid:"Required;MaxSize(50)"`
	NewBranchName string `valid:"AlphaDashDot;MaxSize(100)"`
	Files         []string
}

func (f *UploadRepoFile) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

func (f *UploadRepoFile) IsNewBrnach() bool {
	return f.CommitChoice == "commit-to-new-branch"
}

type RemoveUploadFile struct {
	File string `valid:"Required;MaxSize(50)"`
}

func (f *RemoveUploadFile) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

// ________         .__          __
// \______ \   ____ |  |   _____/  |_  ____
// |    |  \_/ __ \|  | _/ __ \   __\/ __ \
// |    `   \  ___/|  |_\  ___/|  | \  ___/
// /_______  /\___  >____/\___  >__|  \___  >
//         \/     \/          \/          \/

type DeleteRepoFile struct {
	CommitSummary string `valid:"MaxSize(100)"`
	CommitMessage string
	CommitChoice  string `valid:"Required;MaxSize(50)"`
	NewBranchName string `valid:"AlphaDashDot;MaxSize(100)"`
}

func (f *DeleteRepoFile) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

func (f *DeleteRepoFile) IsNewBrnach() bool {
	return f.CommitChoice == "commit-to-new-branch"
}
