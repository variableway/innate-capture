package idea

import "github.com/variableway/innate/capture/internal/config"

const SourceCLI = "capture-cli"
const SourceTUI = "capture-tui"

// Write creates an inbox markdown file for the given idea.
func Write(cfg *config.Config, title, description, context, source string) (string, error) {
	return DefaultService().Write(cfg, CreateInput{
		Title:       title,
		Description: description,
		Context:     context,
		Source:      source,
	})
}

// List returns inbox entries sorted by filename descending (newest date first).
func List(cfg *config.Config) ([]Entry, error) {
	return DefaultService().List(cfg)
}
