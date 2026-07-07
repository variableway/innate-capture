package daily

import "github.com/variableway/innate/capture/internal/config"

// Service defines daily read domain API.
type Service interface {
	ValidSections() []string
	IsValidSection(s string) bool
	Read(cfg *config.Config) (string, error)
	ReadForDate(cfg *config.Config, date string) (string, error)
	BootstrapFromTemplate(cfg *config.Config) error
	BootstrapFromTemplateForDate(cfg *config.Config, date string) error
	PrintSection(cfg *config.Config, section string) (string, error)
	PrintSectionForDate(cfg *config.Config, section, date string) (string, error)
	DailyFilePath(cfg *config.Config, date string) (string, error)
}
