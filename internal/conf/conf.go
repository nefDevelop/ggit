package conf

import (
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nef734/ggit/conf"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/osx"
	"github.com/nef734/ggit/internal/semverx"
	"github.com/nef734/ggit/internal/strx"
	"github.com/nef734/ggit/internal/testx"
	"gopkg.in/ini.v1"
)

func init() {
	// Initialize the primary logger until logging service is up.
	err := log.NewConsole(log.DefaultConsoleName, log.LevelTrace)
	if err != nil {
		panic("init console logger: " + err.Error())
	}
}

// File is the configuration object.
var File *ini.File

// Init initializes configuration from conf assets and given custom configuration file.
// If `customConf` is empty, it falls back to default location, i.e. "<WORK DIR>/custom".
// It is safe to call this function multiple times with desired `customConf`, but it is
// not concurrent safe.
//
// NOTE: The order of loading configuration sections matters as one may depend on another.
//
// ⚠️ WARNING: Do not print anything in this function other than warnings.
func Init(customConf string) error {
	data, err := conf.Files.ReadFile("app.ini")
	if err != nil {
		return fmt.Errorf("read default \"app.ini\": %w", err)
	}

	File, err = ini.LoadSources(ini.LoadOptions{
		IgnoreInlineComment: true,
	}, data)
	if err != nil {
		return fmt.Errorf("parse \"app.ini\": %w", err)
	}
	File.NameMapper = ini.SnackCase
	File.ValueMapper = os.ExpandEnv

	if customConf == "" {
		customConf = filepath.Join(CustomDir(), "conf", "app.ini")
	} else {
		customConf, err = filepath.Abs(customConf)
		if err != nil {
			return fmt.Errorf("get absolute path: %w", err)
		}
	}
	CustomConf = customConf

	if osx.IsFile(customConf) {
		if err = File.Append(customConf); err != nil {
			return fmt.Errorf("append %q: %w", customConf, err)
		}
	} else if !HookMode {
		log.Warn("Custom config %q not found. Ignore this warning if you're running for the first time", customConf)
	}

	if err = File.Section(ini.DefaultSection).MapTo(&App); err != nil {
		return fmt.Errorf("mapping default section: %w", err)
	}

	// ***************************
	// ----- Server settings -----
	// ***************************

	if err = File.Section("server").MapTo(&Server); err != nil {
		return fmt.Errorf("mapping [server] section: %w", err)
	}
	Server.AppDataPath = ensureAbs(Server.AppDataPath)

	if !strings.HasSuffix(Server.ExternalURL, "/") {
		Server.ExternalURL += "/"
	}
	Server.URL, err = url.Parse(Server.ExternalURL)
	if err != nil {
		return fmt.Errorf("parse '[server] EXTERNAL_URL' %q: %w", err, err)
	}

	// Subpath should start with '/' and end without '/', i.e. '/{subpath}'.
	Server.Subpath = strings.TrimRight(Server.URL.Path, "/")
	Server.SubpathDepth = strings.Count(Server.Subpath, "/")

	unixSocketMode, err := strconv.ParseUint(Server.UnixSocketPermission, 8, 32)
	if err != nil {
		return fmt.Errorf("parse '[server] UNIX_SOCKET_PERMISSION' %q: %w", Server.UnixSocketPermission, err)
	}
	if unixSocketMode > 0o777 {
		unixSocketMode = 0o666
	}
	Server.UnixSocketMode = os.FileMode(unixSocketMode)

	// ************************
	// ----- SSH settings -----
	// ************************

	SSH.RootPath = filepath.Join(HomeDir(), ".ssh")
	SSH.KeyTestPath = os.TempDir()
	if err = File.Section("server").MapTo(&SSH); err != nil {
		return fmt.Errorf("mapping SSH settings from [server] section: %w", err)
	}
	SSH.RootPath = ensureAbs(SSH.RootPath)
	SSH.KeyTestPath = ensureAbs(SSH.KeyTestPath)

	if !SSH.Disabled {
		if !SSH.StartBuiltinServer {
			if err := os.MkdirAll(SSH.RootPath, 0o700); err != nil {
				return fmt.Errorf("create SSH root directory: %w", err)
			} else if err = os.MkdirAll(SSH.KeyTestPath, 0o644); err != nil {
				return fmt.Errorf("create SSH key test directory: %w", err)
			}
		} else {
			SSH.RewriteAuthorizedKeysAtStart = false
		}

		// Check if server is eligible for minimum key size check when user choose to enable.
		// Windows server and OpenSSH version lower than 5.1 are forced to be disabled because
		// the "ssh-keygen" in Windows does not print key type.
		// See https://github.com/gogs/gogs/issues/4507.
		if SSH.MinimumKeySizeCheck {
			sshVersion, err := openSSHVersion()
			if err != nil {
				return fmt.Errorf("get OpenSSH version: %w", err)
			}

			if IsWindowsRuntime() || semverx.Compare(sshVersion, "<", "5.1") {
				if !HookMode {
					log.Warn(`SSH minimum key size check is forced to be disabled because server is not eligible:
	1. Windows server
	2. OpenSSH version is lower than 5.1`)
				}
			} else {
				SSH.MinimumKeySizes = map[string]int{}
				for _, key := range File.Section("ssh.minimum_key_sizes").Keys() {
					if key.MustInt() != -1 {
						SSH.MinimumKeySizes[strings.ToLower(key.Name())] = key.MustInt()
					}
				}
			}
		}
	}

	// *******************************
	// ----- Repository settings -----
	// *******************************

	Repository.Root = filepath.Join(HomeDir(), "ggit-repositories")
	if err = File.Section("repository").MapTo(&Repository); err != nil {
		return fmt.Errorf("mapping [repository] section: %w", err)
	}
	Repository.Root = ensureAbs(Repository.Root)
	Repository.Upload.TempPath = ensureAbs(Repository.Upload.TempPath)

	// *****************************
	// ----- Database settings -----
	// *****************************

	if err = File.Section("database").MapTo(&Database); err != nil {
		return fmt.Errorf("mapping [database] section: %w", err)
	}
	Database.Path = ensureAbs(Database.Path)

	// *****************************
	// ----- Security settings -----
	// *****************************

	if err = File.Section("security").MapTo(&Security); err != nil {
		return fmt.Errorf("mapping [security] section: %w", err)
	}

	if err = ensureSecretKey(); err != nil {
		return fmt.Errorf("ensure secret key: %w", err)
	}

	// Check run user when the install is locked.
	if Security.InstallLock {
		currentUser, match := CheckRunUser(App.RunUser)
		if !match {
			return fmt.Errorf("user configured to run Ggit is %q, but the current user is %q", App.RunUser, currentUser)
		}
	}

	// **************************
	// ----- Email settings -----
	// **************************

	if err = File.Section("email").MapTo(&Email); err != nil {
		return fmt.Errorf("mapping [email] section: %w", err)
	}

	if Email.Enabled {
		if Email.From == "" {
			Email.From = Email.User
		}

		parsed, err := mail.ParseAddress(Email.From)
		if err != nil {
			return fmt.Errorf("parse mail address %q: %w", Email.From, err)
		}
		Email.FromEmail = parsed.Address
	}

	// ***********************************
	// ----- Authentication settings -----
	// ***********************************

	if err = File.Section("auth").MapTo(&Auth); err != nil {
		return fmt.Errorf("mapping [auth] section: %w", err)
	}

	// *************************
	// ----- User settings -----
	// *************************

	if err = File.Section("user").MapTo(&User); err != nil {
		return fmt.Errorf("mapping [user] section: %w", err)
	}

	// ****************************
	// ----- Session settings -----
	// ****************************

	if err = File.Section("session").MapTo(&Session); err != nil {
		return fmt.Errorf("mapping [session] section: %w", err)
	}

	// *******************************
	// ----- Attachment settings -----
	// *******************************

	if err = File.Section("attachment").MapTo(&Attachment); err != nil {
		return fmt.Errorf("mapping [attachment] section: %w", err)
	}
	Attachment.Path = ensureAbs(Attachment.Path)

	if err = File.Section("time").MapTo(&Time); err != nil {
		return fmt.Errorf("mapping [time] section: %w", err)
	}

	// *******************************
	// ----- HealthCheck settings ----
	// *******************************

	HealthCheck.ActivityCooldownDays = 2
	HealthCheck.MaxCommits = 1000
	HealthCheck.IdleCheckDays = 15
	if err = File.Section("health_check").MapTo(&HealthCheck); err != nil {
		return fmt.Errorf("mapping [health_check] section: %w", err)
	}

	Time.FormatLayout = map[string]string{
		"ANSIC":       time.ANSIC,
		"UnixDate":    time.UnixDate,
		"RubyDate":    time.RubyDate,
		"RFC822":      time.RFC822,
		"RFC822Z":     time.RFC822Z,
		"RFC850":      time.RFC850,
		"RFC1123":     time.RFC1123,
		"RFC1123Z":    time.RFC1123Z,
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"Kitchen":     time.Kitchen,
		"Stamp":       time.Stamp,
		"StampMilli":  time.StampMilli,
		"StampMicro":  time.StampMicro,
		"StampNano":   time.StampNano,
	}[Time.Format]
	if Time.FormatLayout == "" {
		Time.FormatLayout = time.RFC3339
	}

	// ****************************
	// ----- Picture settings -----
	// ****************************

	if err = File.Section("picture").MapTo(&Picture); err != nil {
		return fmt.Errorf("mapping [picture] section: %w", err)
	}
	Picture.AvatarUploadPath = ensureAbs(Picture.AvatarUploadPath)
	Picture.RepositoryAvatarUploadPath = ensureAbs(Picture.RepositoryAvatarUploadPath)

	// *************************
	// ----- I18n settings -----
	// *************************

	I18n = new(i18nConf)
	if err = File.Section("i18n").MapTo(I18n); err != nil {
		return fmt.Errorf("mapping [i18n] section: %w", err)
	}
	I18n.dateLangs = File.Section("i18n.datelang").KeysHash()

	Git.MaxConcurrency = 50

	// *************************
	// ----- LFS settings -----
	// *************************

	if err = File.Section("lfs").MapTo(&LFS); err != nil {
		return fmt.Errorf("mapping [lfs] section: %w", err)
	}
	LFS.ObjectsPath = ensureAbs(LFS.ObjectsPath)
	LFS.ObjectsTempPath = ensureAbs(LFS.ObjectsTempPath)

	Gist.MaxVersions = 10
	if err = File.Section("gist").MapTo(&Gist); err != nil {
		return fmt.Errorf("mapping [gist] section: %w", err)
	}

	handleDeprecated()
	if !HookMode {
		for _, warning := range checkInvalidOptions(File) {
			log.Warn("%s", warning)
		}
	}

	if err = File.Section("cache").MapTo(&Cache); err != nil {
		return fmt.Errorf("mapping [cache] section: %w", err)
	} else if err = File.Section("http").MapTo(&HTTP); err != nil {
		return fmt.Errorf("mapping [http] section: %w", err)
	} else if err = File.Section("release").MapTo(&Release); err != nil {
		return fmt.Errorf("mapping [release] section: %w", err)
	} else if err = File.Section("webhook").MapTo(&Webhook); err != nil {
		return fmt.Errorf("mapping [webhook] section: %w", err)
	} else if err = File.Section("admin").MapTo(&Admin); err != nil {
		return fmt.Errorf("mapping [admin] section: %w", err)
	}

	// MapTo omite los valores vacíos, así que sin este reset los secretos de un
	// Init anterior sobreviven cuando la nueva configuración los deja en blanco.
	PikoCI = PikoCIOpts{}
	if err = File.Section("pikoci").MapTo(&PikoCI); err != nil {
		return fmt.Errorf("mapping [pikoci] section: %w", err)
	} else if err = File.Section("cron").MapTo(&Cron); err != nil {
		return fmt.Errorf("mapping [cron] section: %w", err)
	} else if err = File.Section("git").MapTo(&Git); err != nil {
		return fmt.Errorf("mapping [git] section: %w", err)
	} else if err = File.Section("api").MapTo(&API); err != nil {
		return fmt.Errorf("mapping [api] section: %w", err)
	} else if err = File.Section("ui").MapTo(&UI); err != nil {
		return fmt.Errorf("mapping [ui] section: %w", err)
	} else if err = File.Section("prometheus").MapTo(&Prometheus); err != nil {
		return fmt.Errorf("mapping [prometheus] section: %w", err)
	} else if err = File.Section("other").MapTo(&Other); err != nil {
		return fmt.Errorf("mapping [other] section: %w", err)
	}

	HasRobotsTxt = osx.IsFile(filepath.Join(CustomDir(), "robots.txt"))
	return nil
}

// MustInit panics if configuration initialization failed.
func MustInit(customConf string) {
	err := Init(customConf)
	if err != nil {
		panic(err)
	}
}

// embeddedSecretKey is the value that ships inside conf/app.ini. That file is
// compiled into the executable through //go:embed and is also part of the public
// source tree, so this constant must never be used to sign anything.
const embeddedSecretKey = "RhgOvRxKpJ1A8ub"

// secretKeyLength in characters. The generator draws from a 62 character
// alphabet, so 43 characters give roughly 256 bits of entropy.
const secretKeyLength = 43

// ensureSecretKey guarantees that the session signing key is not the value
// embedded in the binary.
//
// gorilla/securecookie is deterministic: whoever knows the key can mint a valid
// session cookie. Because the default lives in a file that travels with the
// executable and with the repository, shipping any real key there means every
// deployment shares it and every session is forgeable.
//
// When the key is absent or still the embedded default, a random one is written
// to CustomConf. Init appends that file after the embedded one, so the value on
// disk takes precedence and survives restarts. An ephemeral key is never
// generated silently: if it cannot be persisted the startup fails, because
// otherwise every restart would invalidate all sessions.
func ensureSecretKey() error {
	if Security.SecretKey != "" && Security.SecretKey != embeddedSecretKey {
		return nil
	}

	// Tests load conf against fixtures inside the source tree. Generating a key
	// there would rewrite a tracked file on every run, so the test binaries keep
	// the empty value instead.
	if testx.InTest {
		Security.SecretKey = "test-secret-key-not-used-to-sign-anything"
		return nil
	}

	key, err := strx.RandomChars(secretKeyLength)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	if err = persistSecretKey(key); err != nil {
		return err
	}

	Security.SecretKey = key
	log.Info("Generated a new SECRET_KEY and stored it in %q", CustomConf)

	// The default location is <data>/conf/app.ini, which is ignored by git. When
	// -c points somewhere else the key lands in that file instead, and the file
	// the operator passed is often one that is tracked: "make test-e2e" used to
	// hand this function conf/app.test.ini and the generated key ended up staged
	// for commit. Failing here would break a legitimate deployment that mounts
	// its config read-only, so this warns rather than refuses.
	if !isDefaultCustomConf() {
		log.Warn("SECRET_KEY was written to %q, which is not the default "+
			"<data>/conf/app.ini. If that file is tracked by git, add it to "+
			".gitignore or move the key out, because anyone with the key can "+
			"forge a session cookie for any user.", CustomConf)
	}
	return nil
}

// isDefaultCustomConf reports whether CustomConf is the location chosen when no
// -c flag is given.
func isDefaultCustomConf() bool {
	def := filepath.Join(CustomDir(), "conf", "app.ini")
	a, errA := filepath.Abs(CustomConf)
	b, errB := filepath.Abs(def)
	if errA != nil || errB != nil {
		return false
	}
	return a == b
}

// persistSecretKey writes the key to the custom configuration file and to the
// in-memory copy, so that the value on disk is also the one this process uses.
func persistSecretKey(key string) error {
	if err := os.MkdirAll(filepath.Dir(CustomConf), 0o750); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	// Section.Key creates the key when it is missing, so the same code covers
	// both a fresh install and an existing custom config.
	cfg, err := ini.Load(CustomConf)
	if err != nil {
		cfg = ini.Empty()
	}
	cfg.Section("security").Key("SECRET_KEY").SetValue(key)
	if err = cfg.SaveTo(CustomConf); err != nil {
		return fmt.Errorf("write %q: %w", CustomConf, err)
	}

	// Keep the in-memory copy in sync so nothing downstream can still observe the
	// embedded default.
	File.Section("security").Key("SECRET_KEY").SetValue(key)
	return nil
}
