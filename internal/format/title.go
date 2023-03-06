package format

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func FormatTitle(a string, b string) string {
	aLength := utf8.RuneCountInString(a)
	bLength := utf8.RuneCountInString(b)
	if aLength+bLength > 30 {
		if aLength > 15 {
			aLength = 15
			a = fmt.Sprintf("%s...", strings.TrimSpace(substring(a, 0, aLength)))
		}

		if bLength > 15 {
			bLength = 15
			b = fmt.Sprintf("%s...", strings.TrimSpace(substring(b, 0, bLength)))
		}

	}

	return fmt.Sprintf(
		"%s | %s",
		a, b,
	)
}
