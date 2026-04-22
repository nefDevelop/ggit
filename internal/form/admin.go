package form

type AdminCrateUser struct {
	LoginType  string `binding:"Required"`
	LoginName  string
	UserName   string `binding:"Required;AlphaDashDot;MaxSize(35)"`
	Email      string `binding:"Required;Email;MaxSize(254)"`
	Password   string `binding:"MaxSize(255)"`
	SendNotify bool
}

func (f *AdminCrateUser) Validate(ctx ValidationContext, errs Errors) Errors {
	return validate(errs, ctx.GetData(), f, ctx.GetLocale())
}

type AdminEditUser struct {
	LoginType        string `binding:"Required"`
	LoginName        string
	FullName         string `binding:"MaxSize(100)"`
	Email            string `binding:"Required;Email;MaxSize(254)"`
	Password         string `binding:"MaxSize(255)"`
	Website          string `binding:"MaxSize(50)"`
	Location         string `binding:"MaxSize(50)"`
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
