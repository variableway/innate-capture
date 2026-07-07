package daily

import (
	"github.com/variableway/innate/capture/internal/config"
)

// ValidSections returns the accepted --section values.
func ValidSections() []string {
	return DefaultService().ValidSections()
}

// IsValidSection reports whether s is a recognized section name.
func IsValidSection(s string) bool {
	return DefaultService().IsValidSection(s)
}

// Read returns the full contents of today.md. If today.md does not yet
// exist, it is bootstrapped from daily/_template/day.md first.
func Read(cfg *config.Config) (string, error) {
	return DefaultService().Read(cfg)
}

// BootstrapFromTemplate copies daily/_template/day.md to today.md,
// substituting __DATE__ with today's date. It does not touch checkbox
// state — daily stays a read-only view in capture (--open hands off to
// the user's editor).
func BootstrapFromTemplate(cfg *config.Config) error {
	return DefaultService().BootstrapFromTemplate(cfg)
}

// PrintSection returns a single section of today.md. A section begins at
// a "## " heading whose title starts with the prefix mapped from section,
// and ends at the next "## " heading or end of file. The heading line is
// included in the output.
func PrintSection(cfg *config.Config, section string) (string, error) {
	return DefaultService().PrintSection(cfg, section)
}
