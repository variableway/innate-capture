package daily

import "github.com/variableway/innate/capture/internal/config"

// Service defines daily read domain API.
type Service interface {
	ValidSections() []string
	IsValidSection(s string) bool
	Read(cfg *config.Config) (string, error)
	BootstrapFromTemplate(cfg *config.Config) error
	PrintSection(cfg *config.Config, section string) (string, error)
}
