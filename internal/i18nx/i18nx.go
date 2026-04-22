package i18nx

import (
	"fmt"
)

// trMap contains the minimal translations needed by the backend.
// Ggit has moved most translations to the frontend.
var trMap = map[string]map[string]string{
	"en": {
		"tool.now":      "now",
		"tool.1s":       "1 second %s",
		"tool.seconds":  "%d seconds %s",
		"tool.1m":       "1 minute %s",
		"tool.minutes":  "%d minutes %s",
		"tool.1h":       "1 hour %s",
		"tool.hours":    "%d hours %s",
		"tool.1d":       "1 day %s",
		"tool.days":     "%d days %s",
		"tool.1w":       "1 week %s",
		"tool.weeks":    "%d weeks %s",
		"tool.1mon":     "1 month %s",
		"tool.months":   "%d months %s",
		"tool.1y":       "1 year %s",
		"tool.years":    "%d years %s",
		"tool.ago":      "ago",
		"tool.from_now": "from now",
	},
	"es": {
		"tool.now":      "ahora",
		"tool.1s":       "hace 1 segundo %s",
		"tool.seconds":  "hace %d segundos %s",
		"tool.1m":       "hace 1 minuto %s",
		"tool.minutes":  "hace %d minutos %s",
		"tool.1h":       "hace 1 hora %s",
		"tool.hours":    "hace %d horas %s",
		"tool.1d":       "hace 1 día %s",
		"tool.days":     "hace %d días %s",
		"tool.1w":       "hace 1 semana %s",
		"tool.weeks":    "hace %d semanas %s",
		"tool.1mon":     "hace 1 mes %s",
		"tool.months":   "hace %d meses %s",
		"tool.1y":       "hace 1 año %s",
		"tool.years":    "hace %d años %s",
		"tool.ago":      "", // En español el "hace" va delante, lo he incluido arriba
		"tool.from_now": "desde ahora",
	},
}

// Tr performs a simple translation for the given language and key.
func Tr(lang, key string, args ...any) string {
	l, ok := trMap[lang]
	if !ok {
		l = trMap["en"]
	}

	format, ok := l[key]
	if !ok {
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}
