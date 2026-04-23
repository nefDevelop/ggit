package form

type AdminCrateUser struct {
	LoginType  string `valid:"Required"`
	LoginName  string
	UserName   string `valid:"Required;AlphaDashDot;MaxSize(35)"`
	Email      string `valid:"Required;Email;MaxSize(254)"`
	Password   string `valid:"MaxSize(255)"`
	SendNotify bool
}

func (f *AdminCrateUser) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type AdminEditUser struct {
	LoginType        string `valid:"Required"`
	LoginName        string
	FullName         string `valid:"MaxSize(100)"`
	Email            string `valid:"Required;Email;MaxSize(254)"`
	Password         string `valid:"MaxSize(255)"`
	Website          string `valid:"MaxSize(50)"`
	Location         string `valid:"MaxSize(50)"`
	MaxRepoCreation  int
	Active           bool
	Admin            bool
	AllowGitHook     bool
	AllowImportLocal bool
	ProhibitLogin    bool
}

func (f *AdminEditUser) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}
