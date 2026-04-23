package form

type CreateOrg struct {
	OrgName string `valid:"Required;AlphaDashDot;MaxSize(35)" locale:"org.org_name_holder"`
}

func (f *CreateOrg) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type UpdateOrgSetting struct {
	Name            string `valid:"Required;AlphaDashDot;MaxSize(35)" locale:"org.org_name_holder"`
	FullName        string `valid:"MaxSize(100)"`
	Description     string `valid:"MaxSize(255)"`
	Website         string `valid:"Url;MaxSize(100)"`
	Location        string `valid:"MaxSize(50)"`
	MaxRepoCreation int
}

func (f *UpdateOrgSetting) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type CreateTeam struct {
	TeamName    string `valid:"Required;AlphaDashDot;MaxSize(30)"`
	Description string `valid:"MaxSize(255)"`
	Permission  string
}

func (f *CreateTeam) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}
