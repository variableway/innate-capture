package workspace

import (
	"github.com/variableway/innate/capture/internal/config"
)

func Root(cfg *config.Config) string {
	return DefaultService().Root(cfg)
}

func InboxDir(cfg *config.Config) string {
	return DefaultService().InboxDir(cfg)
}

func DailyPath(cfg *config.Config) string {
	return DefaultService().DailyPath(cfg)
}

func Validate(cfg *config.Config) error {
	return DefaultService().Validate(cfg)
}
