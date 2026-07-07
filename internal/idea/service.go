package idea

import "github.com/variableway/innate/capture/internal/config"

// Service defines the idea domain API.
type Service interface {
	Write(cfg *config.Config, in CreateInput) (string, error)
	List(cfg *config.Config) ([]Entry, error)
}
