package dbx

import (
	"fmt"
)

// Quote adds surrounding backticks for all given arguments before being formatted.
func Quote(format string, args ...string) string {
	anys := make([]any, len(args))
	for i := range args {
		anys[i] = "`" + args[i] + "`"
	}
	return fmt.Sprintf(format, anys...)
}
