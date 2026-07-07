package idea

import (
	"regexp"
	"strings"
	"unicode"
)

var multiDash = regexp.MustCompile(`-+`)

// Slug converts a title into a filesystem-safe slug per workspace-io-v1.
func Slug(title string) string {
	var b strings.Builder
	prevDash := false

	for _, r := range strings.ToLower(strings.TrimSpace(title)) {
		switch {
		case r == ' ' || r == '_':
			if !prevDash && b.Len() > 0 {
				b.WriteByte('-')
				prevDash = true
			}
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '.' || r == '-':
			b.WriteRune(r)
			prevDash = false
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			// non-ascii letters/digits: skip (contract allows only [a-z0-9._-])
		}
	}

	slug := strings.Trim(multiDash.ReplaceAllString(b.String(), "-"), "-.")
	if slug == "" {
		return "idea"
	}
	return slug
}
