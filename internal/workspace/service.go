package workspace

import "github.com/variableway/innate/capture/internal/config"

// Service defines the workspace domain API.
type Service interface {
	Root(cfg *config.Config) string
	InboxDir(cfg *config.Config) string
	DailyPath(cfg *config.Config) string
	Resolve(cfg *config.Config) Paths
	Validate(cfg *config.Config) error
}
