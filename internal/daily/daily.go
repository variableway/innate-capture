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

func ReadForDate(cfg *config.Config, date string) (string, error) {
	return DefaultService().ReadForDate(cfg, date)
}

// BootstrapFromTemplate copies daily/_template/day.md to today.md,
// substituting __DATE__ with today's date. It does not touch checkbox
// state — daily stays a read-only view in capture (--open hands off to
// the user's editor).
func BootstrapFromTemplate(cfg *config.Config) error {
	return DefaultService().BootstrapFromTemplate(cfg)
}

func BootstrapFromTemplateForDate(cfg *config.Config, date string) error {
	return DefaultService().BootstrapFromTemplateForDate(cfg, date)
}

// PrintSection returns a single section of today.md. A section begins at
// a "## " heading whose title starts with the prefix mapped from section,
// and ends at the next "## " heading or end of file. The heading line is
// included in the output.
func PrintSection(cfg *config.Config, section string) (string, error) {
	return DefaultService().PrintSection(cfg, section)
}

func PrintSectionForDate(cfg *config.Config, section, date string) (string, error) {
	return DefaultService().PrintSectionForDate(cfg, section, date)
}

func DailyFilePath(cfg *config.Config, date string) (string, error) {
	return DefaultService().DailyFilePath(cfg, date)
}
