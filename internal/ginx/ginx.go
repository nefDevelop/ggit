package ginx

import (
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nef734/ggit/internal/conf"
	"github.com/nef734/ggit/internal/context"
	"github.com/nef734/ggit/internal/database"
	"github.com/nef734/ggit/internal/form"
	"github.com/nef734/ggit/internal/osx"
	"github.com/nef734/ggit/public"
)

type Engine struct {
	*gin.Engine
}

func New() *Engine {
	if conf.IsProdMode() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := &Engine{Engine: gin.New()}

	if !conf.Server.DisableRouterLog {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())
	r.Use(CORS())
	if conf.Server.EnableGzip {
		r.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	store := cookie.NewStore([]byte(conf.Security.SecretKey))
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   86400,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("ggit_session", store))
	r.Use(context.Contexter(context.NewStore()))
	r.Use(CSRF())

	setupStatic(r)
	return r
}

func Handle(h any) gin.HandlerFunc {
	v := reflect.ValueOf(h)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	return func(c *gin.Context) {
		if c.IsAborted() {
			return
		}
		gctx := c.MustGet("context").(*context.Context)

		inputs := make([]reflect.Value, t.NumIn())
		for i := 0; i < t.NumIn(); i++ {
			argType := t.In(i)

			// Detect known types
			switch argType.String() {
			case "*context.Context":
				inputs[i] = reflect.ValueOf(gctx)
			case "*context.APIContext":
				actx, _ := c.Get("apiContext")
				inputs[i] = reflect.ValueOf(actx)
			case "*context.ParamsUser":
				pu, _ := c.Get("paramsUser")
				inputs[i] = reflect.ValueOf(pu)
			default:
				// If it's a pointer to a struct, try to bind it (legacy binding)
				if argType.Kind() == reflect.Ptr && argType.Elem().Kind() == reflect.Struct {
					val := reflect.New(argType.Elem())
					if err := c.ShouldBind(val.Interface()); err != nil {
						// Macaron ignored binding errors in many places or handled them in the handler
						// For now, let's just proceed to stay close to legacy behavior
					}
					if f, ok := val.Interface().(form.Form); ok {
						f.Validate(gctx, nil)
					}
					inputs[i] = val
				} else if argType.Kind() == reflect.Struct {
					val := reflect.New(argType).Elem()
					if err := c.ShouldBind(val.Addr().Interface()); err != nil {
					}
					if f, ok := val.Addr().Interface().(form.Form); ok {
						f.Validate(gctx, nil)
					}
					inputs[i] = val
				} else {
					panic("unsupported argument type: " + argType.String())
				}
			}
		}
		v.Call(inputs)
	}
}

func (r *Engine) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return r.Engine.Group(path, handlers...)
}

type Combo struct {
	rg   *gin.RouterGroup
	path string
}

func (c *Combo) GET(h any, handlers ...gin.HandlerFunc) *Combo {
	handlers = append(handlers, Handle(h))
	c.rg.GET(c.path, handlers...)
	return c
}

func (c *Combo) POST(h any, handlers ...gin.HandlerFunc) *Combo {
	handlers = append(handlers, Handle(h))
	c.rg.POST(c.path, handlers...)
	return c
}

func (c *Combo) DELETE(h any, handlers ...gin.HandlerFunc) *Combo {
	handlers = append(handlers, Handle(h))
	c.rg.DELETE(c.path, handlers...)
	return c
}

func UseCombo(rg *gin.RouterGroup, path string) *Combo {
	return &Combo{rg: rg, path: path}
}

func setupStatic(r *Engine) {
	// Ensure WOFF2 MIME type is registered (vital for Alpine-based containers)
	_ = mime.AddExtensionType(".woff2", "font/woff2")

	customPublic := filepath.Join(conf.CustomDir(), "public")

	serveStatic := func(prefix, path string) {
		r.Static(prefix, path)
	}
	serveStaticFS := func(prefix string, fs http.FileSystem) {
		r.StaticFS(prefix, fs)
	}

	if osx.IsDir(customPublic) {
		serveStatic("/public", customPublic)
		if osx.IsDir(filepath.Join(customPublic, "img")) {
			serveStatic("/img", filepath.Join(customPublic, "img"))
		}
	} else if conf.Server.LoadAssetsFromDisk {
		publicDir := filepath.Join(conf.WorkDir(), "public")
		serveStatic("/public", publicDir)
		serveStatic("/img", filepath.Join(publicDir, "img"))
	} else {
		serveStaticFS("/public", http.FS(public.Files))
		sub, err := fs.Sub(public.Files, "img")
		if err == nil {
			serveStaticFS("/img", http.FS(sub))
		}
	}

	r.Static(conf.UsersAvatarPathPrefix, conf.Picture.AvatarUploadPath)
	r.Static(database.RepoAvatarURLPrefix, conf.Picture.RepositoryAvatarUploadPath)

	// Svelte Build Static Files
	svelteBuildDir := filepath.Join(conf.WorkDir(), "frontend/build")
	r.Static("/_app", filepath.Join(svelteBuildDir, "_app"))

	// Static Assets (Fonts, etc.)
	r.Static("/fonts", filepath.Join("frontend", "build", "fonts"))
}

// ServeSPA serves the Svelte frontend entry point (index.html).
func ServeSPA(c *gin.Context) {
	c.File(filepath.Join(conf.WorkDir(), "frontend/build/index.html"))
}
