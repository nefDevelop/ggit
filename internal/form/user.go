package form

import (
	"mime/multipart"
)

//nolint:staticcheck // Reason: needed for legacy code
type Install struct {
	DbType   string `json:"db_type" valid:"Required"`
	DbHost   string `json:"db_host"`
	DbUser   string `json:"db_user"`
	DbPasswd string `json:"db_password"`
	DbName   string `json:"db_name"`
	DbSchema string `json:"db_schema"`
	SSLMode  string `json:"ssl_mode"`
	DbPath   string `json:"db_path"`

	AppName             string `json:"app_name" valid:"Required" locale:"install.app_name"`
	RepoRootPath        string `json:"repo_root_path" valid:"Required"`
	RunUser             string `json:"run_user" valid:"Required"`
	Domain              string `json:"domain" valid:"Required"`
	SSHPort             int    `json:"ssh_port"`
	UseBuiltinSSHServer bool   `json:"use_builtin_ssh"`
	HTTPPort            string `json:"http_port" valid:"Required"`
	AppUrl              string `json:"app_url" valid:"Required"`
	LogRootPath         string `json:"log_root_path" valid:"Required"`
	EnableConsoleMode   bool   `json:"enable_console_mode"`
	DefaultBranch       string `json:"default_branch"`

	SMTPHost        string `json:"smtp_host"`
	SMTPFrom        string `json:"smtp_from"`
	SMTPUser        string `json:"smtp_user" valid:"OmitEmpty;MaxSize(254)" locale:"install.mailer_user"`
	SMTPPasswd      string `json:"smtp_password"`
	RegisterConfirm bool   `json:"register_confirm"`
	MailNotify      bool   `json:"mail_notify"`

	OfflineMode         bool `json:"offline_mode"`
	DisableRegistration bool `json:"disable_registration"`
	RequireSignInView   bool `json:"require_sign_in_view"`

	AdminName          string `json:"admin_name" valid:"OmitEmpty;AlphaDashDot;MaxSize(30)" locale:"install.admin_name"`
	AdminPasswd        string `json:"admin_password" valid:"OmitEmpty;MaxSize(255)" locale:"install.admin_password"`
	AdminConfirmPasswd string `json:"admin_confirm_password"`
	AdminEmail         string `json:"admin_email" valid:"OmitEmpty;MinSize(3);MaxSize(254);Include(@)" locale:"install.admin_email"`
}

func (f *Install) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

//    _____   ____ _________________ ___
//   /  _  \ |    |   \__    ___/   |   \
//  /  /_\  \|    |   / |    | /    ~    \
// /    |    \    |  /  |    | \    Y    /
// \____|__  /______/   |____|  \___|_  /
//         \/                         \/

type Register struct {
	UserName string `json:"user_name" valid:"Required;AlphaDashDot;MaxSize(35)"`
	Email    string `json:"email" valid:"Required;Email;MaxSize(254)"`
	Password string `json:"password" valid:"Required;MaxSize(255)"`
	Retype   string `json:"retype"`
}

func (f *Register) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type SignIn struct {
	UserName    string `valid:"Required;MaxSize(254)"`
	Password    string `valid:"Required;MaxSize(255)"`
	LoginSource int64
	Remember    bool
}

func (f *SignIn) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

//   __________________________________________.___ _______    ________  _________
//  /   _____/\_   _____/\__    ___/\__    ___/|   |\      \  /  _____/ /   _____/
//  \_____  \  |    __)_   |    |     |    |   |   |/   |   \/   \  ___ \_____  \
//  /        \ |        \  |    |     |    |   |   /    |    \    \_\  \/        \
// /_______  //_______  /  |____|     |____|   |___\____|__  /\______  /_______  /
//         \/         \/                                   \/        \/        \/

type UpdateProfile struct {
	Name            string `valid:"Required;AlphaDashDot;MaxSize(35)"`
	FullName        string `valid:"MaxSize(100)"`
	Website         string `valid:"Url;MaxSize(100)"`
	Location        string `valid:"MaxSize(50)"`
	DisableNewsFeed bool   `json:"disable_news_feed"`
}

func (f *UpdateProfile) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

const (
	AvatarLocal string = "local"
)

type Avatar struct {
	Source string
	Avatar *multipart.FileHeader
}

func (f *Avatar) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type AddEmail struct {
	Email string `valid:"Required;Email;MaxSize(254)"`
}

func (f *AddEmail) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type ChangePassword struct {
	OldPassword string `valid:"Required;MinSize(1);MaxSize(255)"`
	Password    string `valid:"Required;MaxSize(255)"`
	Retype      string
}

func (f *ChangePassword) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type AddSSHKey struct {
	Title   string `valid:"Required;MaxSize(50)"`
	Content string `valid:"Required"`
}

func (f *AddSSHKey) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type NewAccessToken struct {
	Name string `valid:"Required"`
}

func (f *NewAccessToken) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}
