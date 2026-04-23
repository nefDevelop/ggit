package form

type Authentication struct {
	ID                int64
	Type              int    `valid:"Range(2,6)"`
	Name              string `valid:"Required;MaxSize(30)"`
	Host              string
	Port              int
	BindDN            string
	BindPassword      string
	UserBase          string
	UserDN            string
	AttributeUsername string
	AttributeName     string
	AttributeSurname  string
	AttributeMail     string
	AttributesInBind  bool
	Filter            string
	AdminFilter       string
	GroupEnabled      bool
	GroupDN           string
	GroupFilter       string
	GroupMemberUID    string
	UserUID           string
	IsActive          bool
	IsDefault         bool
	SMTPAuth          string
	SMTPHost          string
	SMTPPort          int
	AllowedDomains    string
	SecurityProtocol  int `valid:"Range(0,2)"`
	TLS               bool
	SkipVerify        bool
	PAMServiceName    string
	GitHubAPIEndpoint string `form:"github_api_endpoint" valid:"Url"`
}

func (f *Authentication) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}
