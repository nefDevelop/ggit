package context

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/errx"
	"gogs.io/gogs/internal/template"
)

type Locale interface {
	Tr(string, ...any) string
}

type Captcha interface {
	VerifyReq(req *http.Request) bool
}

type Cache interface {
	Get(key string) any
	Put(key string, val any, timeout int64) error
	Delete(key string) error
	IsExist(key string) bool
}

type Context struct {
	*gin.Context
	Data    map[string]any
	Cache   Cache
	Flash   *Flash
	Session sessions.Session
	csrf    any
	Locale  Locale

	Link        string // Current request URL
	User        *database.User
	IsLogged    bool
	IsBasicAuth bool
	IsTokenAuth bool
	Req         *http.Request // Alias for c.Request for legacy support

	Repo *Repository
	Org  *Organization
}

type Flash struct {
	ErrorMsg string
	SuccessMsg string
	InfoMsg string
}

func (f *Flash) Error(msg string, args ...any) { f.ErrorMsg = fmt.Sprintf(msg, args...) }
func (f *Flash) Success(msg string, args ...any) { f.SuccessMsg = fmt.Sprintf(msg, args...) }
func (f *Flash) Info(msg string, args ...any) { f.InfoMsg = fmt.Sprintf(msg, args...) }

// RawTitle sets the "Title" field in template data.
func (c *Context) RawTitle(title string) {
	c.Data["Title"] = title
}

// Title localizes the "Title" field in template data.
func (c *Context) Title(locale string) {
	c.RawTitle(c.Tr(locale))
}

// Tr returns the localized string for given key.
func (c *Context) Tr(key string, args ...any) string {
	if c.Locale != nil {
		return c.Locale.Tr(key, args...)
	}
	return key
}

// PageIs sets "PageIsxxx" field in template data.
func (c *Context) PageIs(name string) {
	c.Data["PageIs"+name] = true
}

// UserID returns ID of current logged user.
func (c *Context) UserID() int64 {
	if c.User != nil {
		return c.User.ID
	}
	return 0
}

func (c *Context) HasError() bool {
	_, ok := c.Data["HasError"]
	return ok
}

func (c *Context) Params(name string) string {
	return c.Param(strings.TrimPrefix(name, ":"))
}

func (c *Context) QueryInt(name string) int {
	v, _ := strconv.Atoi(c.Query(name))
	return v
}

func (c *Context) QueryInt64(name string) int64 {
	v, _ := strconv.ParseInt(c.Query(name), 10, 64)
	return v
}

func (c *Context) QueryStrings(name string) []string {
	return c.Request.URL.Query()[name]
}

func (c *Context) RawData() ([]byte, error) {
	return io.ReadAll(c.Request.Body)
}

func (c *Context) FormInt(name string) int {
	v, _ := strconv.Atoi(c.PostForm(name))
	return v
}

func (c *Context) FormInt64(name string) int64 {
	v, _ := strconv.ParseInt(c.PostForm(name), 10, 64)
	return v
}

func (c *Context) HTML(status int, name string) {
	c.Context.HTML(status, name, c.Data)
}

func (c *Context) Success(name string) {
	c.HTML(http.StatusOK, name)
}

func (c *Context) Redirect(location string, status ...int) {
	code := http.StatusFound
	if len(status) > 0 {
		code = status[0]
	}
	c.Context.Redirect(code, template.EscapePound(location))
}

func (c *Context) RedirectSubpath(location string, status ...int) {
	c.Redirect(conf.Server.Subpath+location, status...)
}

func (c *Context) RenderWithErr(msg string, status int, tpl string, form any) {
	if form != nil {
		formErr := "HasError"
		c.Data[formErr] = true
		c.Data["ErrorMsg"] = msg
		// TODO: reflect form to Data
	}
	c.HTML(status, tpl)
}

func (c *Context) NotFound() {
	c.HTML(http.StatusNotFound, "status/404")
}

func (c *Context) QueryBool(name string) bool {
	v := c.Query(name)
	return v == "true" || v == "1" || v == "on"
}

func (c *Context) Resp() http.ResponseWriter {
	return c.Context.Writer
}

func (c *Context) RequireHighlightJS() {
	c.Data["RequireHighlightJS"] = true
}

func (c *Context) Written() bool {
	return c.IsAborted()
}

func (c *Context) Error(err error, msg string) {
	log.ErrorDepth(4, "%s: %v", msg, err)
	c.HTML(http.StatusInternalServerError, "status/500")
}

func (c *Context) NotFoundOrError(err error, msg string) {
	if errx.IsNotFound(err) {
		c.HTML(http.StatusNotFound, "status/404")
		return
	}
	c.Error(err, msg)
}

func (c *Context) JSONSuccess(v any) {
	c.JSON(http.StatusOK, v)
}

func (c *Context) ServeContent(name string, r io.ReadSeeker, params ...any) {
	modtime := time.Now()
	for _, p := range params {
		if t, ok := p.(time.Time); ok {
			modtime = t
		}
	}
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+name)
	http.ServeContent(c.Writer, c.Request, name, modtime, r)
}

func (c *Context) ParamsInt64(name string) int64 {
	v, _ := strconv.ParseInt(c.Param(name), 10, 64)
	return v
}

func Contexter(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		gctx := &Context{
			Context: c,
			Data:    make(map[string]any),
			Session: sess,
			Flash:   &Flash{},
			Link:    conf.Server.Subpath + strings.TrimSuffix(c.Request.URL.Path, "/"),
			Req:     c.Request,
			Repo:    &Repository{PullRequest: &PullRequest{}},
			Org:     &Organization{},
		}
		gctx.Data["Flash"] = gctx.Flash
		
		gctx.User, gctx.IsBasicAuth, gctx.IsTokenAuth = authenticatedUser(store, c, sess)
		if gctx.User != nil {
			gctx.IsLogged = true
			gctx.Data["LoggedUser"] = gctx.User
			gctx.Data["LoggedUserID"] = gctx.User.ID
			gctx.Data["LoggedUserName"] = gctx.User.Name
			gctx.Data["IsAdmin"] = gctx.User.IsAdmin
		}

		c.Set("context", gctx)
		c.Next()
	}
}


func (c *Context) GetData() map[string]any {
	return c.Data
}

func (c *Context) GetLocale() Locale {
	return c.Locale
}
