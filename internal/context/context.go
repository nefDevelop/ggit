package context

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"

	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/errx"
	"github.com/nef734/ggit/internal/template"
)

type Locale interface {
	Tr(string, ...any) string
	Language() string
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
	Resp    http.ResponseWriter

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
	ErrorMsg   string
	SuccessMsg string
	InfoMsg    string
	WarningMsg string
}

func (f *Flash) Error(msg string, args ...any)   { f.ErrorMsg = fmt.Sprintf(msg, args...) }
func (f *Flash) Success(msg string, args ...any) { f.SuccessMsg = fmt.Sprintf(msg, args...) }
func (f *Flash) Info(msg string, args ...any)    { f.InfoMsg = fmt.Sprintf(msg, args...) }
func (f *Flash) Warning(msg string, args ...any) { f.WarningMsg = fmt.Sprintf(msg, args...) }

// Tr returns the localized string for given key.
func (c *Context) Tr(key string, args ...any) string {
	if c.Locale != nil {
		return c.Locale.Tr(key, args...)
	}
	return key
}

func (c *Context) Language() string {
	if c.Locale != nil {
		return c.Locale.Language()
	}
	return "en-US"
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

func (c *Context) HasValue(name string) bool {
	_, ok := c.Data[name]
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

func (c *Context) QueryTrim(name string) string {
	return strings.TrimSpace(c.Query(name))
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

// HTML is now lobotomized (Zero Template). It logs and returns JSON if an error occurred.
func (c *Context) HTML(status int, name string) {
	log.Trace("Legacy HTML call (disconnected): %s", name)
	if status >= 400 {
		c.JSON(status, map[string]any{"error": name})
		c.Abort()
	}
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
	c.Abort()
}

func (c *Context) RedirectSubpath(location string, status ...int) {
	c.Redirect(conf.Server.Subpath+location, status...)
}

func (c *Context) RenderWithErr(msg string, status int, tpl string, form any) {
	log.Trace("Legacy RenderWithErr call (disconnected): %s - %s", tpl, msg)
	c.Data["ErrorMsg"] = msg
	c.HTML(status, tpl)
}

func (c *Context) NotFound() {
	if strings.HasPrefix(c.Request.URL.Path, "/api/v2") {
		c.JSON(http.StatusNotFound, map[string]any{"error": "Resource not found"})
		c.Abort()
		return
	}
	c.HTML(http.StatusNotFound, "status/404")
	c.Abort()
}

func (c *Context) QueryBool(name string) bool {
	v := c.Query(name)
	return v == "true" || v == "1" || v == "on"
}

func (c *Context) Title(name string) {
	c.Data["Title"] = c.Tr(name)
}

func (c *Context) PageIs(name string) {
	c.Data["PageIs"+name] = true
}

func (c *Context) RequireHighlightJS()   {}
func (c *Context) RequireSimpleMDE()    {}
func (c *Context) RequireSimpleMDECDE() {}
func (c *Context) RequireAutosize()     {}
func (c *Context) RequireDropzone()     {}

func (c *Context) FormErr(names ...string) {
	for _, name := range names {
		c.Data["Err_"+name] = true
	}
	c.Data["HasError"] = true
}

func (c *Context) Written() bool {
	return c.IsAborted()
}

func (c *Context) RawRedirect(location string) {
	c.Context.Redirect(http.StatusFound, location)
}

func (c *Context) PlainText(status int, body string) {
	c.String(status, body)
}

func (c *Context) SetCookiePath(name, value string, maxAge int, path string) {
	c.Context.SetCookie(name, value, maxAge, path, "", false, false)
}

func (c *Context) GetCookie(name string) string {
	v, _ := c.Cookie(name)
	return v
}

func (c *Context) GetSuperSecureCookie(key, name string) (string, bool) {
	v, err := c.Cookie(name)
	return v, err == nil && v != ""
}

func (c *Context) SetCookie(name, value string, maxAge int, path string, args ...any) {
	domain := ""
	secure := false
	httpOnly := false
	if len(args) > 0 {
		if v, ok := args[0].(string); ok {
			domain = v
		}
	}
	if len(args) > 1 {
		if v, ok := args[1].(bool); ok {
			secure = v
		}
	}
	if len(args) > 2 {
		if v, ok := args[2].(bool); ok {
			httpOnly = v
		}
	}
	c.Context.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *Context) SetSuperSecureCookie(key, name, value string, maxAge int, path string, args ...any) {
	c.SetCookie(name, value, maxAge, path, args...)
}

func (c *Context) Status(code int) {
	c.Context.Status(code)
}

func (c *Context) ServeFile(file string, names ...string) {
	var name string
	if len(names) > 0 {
		name = names[0]
	} else {
		name = filepath.Base(file)
	}
	c.FileAttachment(file, name)
}

func (c *Context) Error(err error, msg string) {
	log.ErrorDepth(4, "%s: %v", msg, err)
	if isAPIPath(c.Request.URL.Path) {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"error":   msg,
			"details": err.Error(),
		})
		c.Abort()
		return
	}
	c.HTML(http.StatusInternalServerError, "status/500")
	c.Abort()
}

func (c *Context) Errorf(err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.ErrorDepth(4, "%s: %v", msg, err)
	c.Error(err, msg)
}

func (c *Context) StatusErrorf(status int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.ErrorDepth(4, "%s", msg)
	if strings.HasPrefix(c.Request.URL.Path, "/api/v2") {
		c.JSON(status, map[string]any{"error": msg})
		c.Abort()
		return
	}
	c.HTML(status, "status/500")
	c.Abort()
}

func (c *Context) BaseURL() string {
	return conf.Server.ExternalURL + c.Link
}

func (c *Context) SetParams(name, val string) {
	// Not easily supported in Gin without hacking, but let's see
	// Gin params are immutable in the context usually.
}

func (c *Context) Write(p []byte) (int, error) {
	return c.Writer.Write(p)
}

func (c *Context) GetErrMsg() string {
	return c.Data["ErrorMsg"].(string)
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
			Resp:    c.Writer,
			Repo:    &Repository{PullRequest: &PullRequest{}},
			Org:     &Organization{},
		}
		gctx.Data["Flash"] = gctx.Flash
		
		if database.HasEngine && database.Handle != nil {
			gctx.User, gctx.IsBasicAuth, gctx.IsTokenAuth = authenticatedUser(store, c, sess)
			if gctx.User != nil {
				gctx.IsLogged = true
			}
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
