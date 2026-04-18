package v2

import (
	"net/http"
	"strings"
	"sync"

	"gopkg.in/ini.v1"
	"gogs.io/gogs/internal/context"
	gogsconf "gogs.io/gogs/conf"
)

var (
	i18nCache = make(map[string]map[string]map[string]string)
	i18nMutex sync.RWMutex
)

// GetI18nStrings returns the localized strings for the requested language as JSON.
func GetI18nStrings(c *context.APIContext) {
	lang := c.Params("lang")
	if lang == "" {
		lang = "en-US"
	}

	// Remove .json suffix if present for cleaner URLs
	lang = strings.TrimSuffix(lang, ".json")

	i18nMutex.RLock()
	if data, ok := i18nCache[lang]; ok {
		i18nMutex.RUnlock()
		c.JSON(http.StatusOK, data)
		return
	}
	i18nMutex.RUnlock()

	// Load from embedded assets
	path := "locale/locale_" + lang + ".ini"
	rawData, err := gogsconf.Files.ReadFile(path)
	if err != nil {
		// Try fallback to just lang part if it was like en-US -> en
		if strings.Contains(lang, "-") {
			baseLang := strings.Split(lang, "-")[0]
			path = "locale/locale_" + baseLang + ".ini"
			rawData, err = gogsconf.Files.ReadFile(path)
		}
		
		if err != nil {
			c.StatusError(http.StatusNotFound, "Locale not found", err)
			return
		}
	}

	f, err := ini.Load(rawData)
	if err != nil {
		c.Error(err, "parse locale ini")
		return
	}

	data := make(map[string]map[string]string)
	for _, section := range f.Sections() {
		name := section.Name()
		if name == ini.DefaultSection {
			name = "common"
		}
		
		secData := make(map[string]string)
		for _, key := range section.Keys() {
			secData[key.Name()] = key.Value()
		}
		data[name] = secData
	}

	i18nMutex.Lock()
	i18nCache[lang] = data
	i18nMutex.Unlock()

	c.JSON(http.StatusOK, data)
}
