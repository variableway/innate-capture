package daily

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/variableway/innate/capture/internal/model"
	"github.com/variableway/innate/capture/internal/workspace"
)

// sectionHeadings maps CLI --section keys to the heading prefix they
// match in today.md (the text immediately after "## ").
var sectionHeadings = map[string]string{
	"input":  "输入",
	"output": "输出",
	"ideas":  "Ideas",
}

// ValidSections returns the accepted --section values.
func ValidSections() []string {
	return []string{"input", "output", "ideas"}
}

// IsValidSection reports whether s is a recognized section name.
func IsValidSection(s string) bool {
	_, ok := sectionHeadings[s]
	return ok
}

// templatePath returns the daily template path under the workspace root:
// {root}/daily/_template/day.md.
func templatePath(cfg *model.Config) string {
	return filepath.Join(workspace.Root(cfg), "daily", "_template", "day.md")
}

// Read returns the full contents of today.md. If today.md does not yet
// exist, it is bootstrapped from daily/_template/day.md first.
func Read(cfg *model.Config) (string, error) {
	if err := workspace.Validate(cfg); err != nil {
		return "", err
	}

	path := workspace.DailyPath(cfg)
	data, err := os.ReadFile(path)
	if err == nil {
		return string(data), nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("read today: %w", err)
	}

	if err := BootstrapFromTemplate(cfg); err != nil {
		return "", err
	}

	data, err = os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read today after bootstrap: %w", err)
	}
	return string(data), nil
}

// BootstrapFromTemplate copies daily/_template/day.md to today.md,
// substituting __DATE__ with today's date. It does not touch checkbox
// state — daily stays a read-only view in capture (--open hands off to
// the user's editor).
func BootstrapFromTemplate(cfg *model.Config) error {
	if err := workspace.Validate(cfg); err != nil {
		return err
	}

	tpl, err := os.ReadFile(templatePath(cfg))
	if err != nil {
		return fmt.Errorf("read daily template: %w", err)
	}

	today := time.Now().Format("2006-01-02")
	content := strings.ReplaceAll(string(tpl), "__DATE__", today)

	path := workspace.DailyPath(cfg)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create daily dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write today: %w", err)
	}
	return nil
}

// PrintSection returns a single section of today.md. A section begins at
// a "## " heading whose title starts with the prefix mapped from section,
// and ends at the next "## " heading or end of file. The heading line is
// included in the output.
func PrintSection(cfg *model.Config, section string) (string, error) {
	if !IsValidSection(section) {
		return "", fmt.Errorf("invalid section %q; valid: %s", section, strings.Join(ValidSections(), ", "))
	}

	content, err := Read(cfg)
	if err != nil {
		return "", err
	}

	prefix := sectionHeadings[section]
	lines := strings.Split(content, "\n")
	var out []string
	inSection := false
	for _, line := range lines {
		isH2 := strings.HasPrefix(line, "## ")
		if isH2 {
			if inSection {
				break // reached the next section
			}
			if strings.HasPrefix(strings.TrimPrefix(line, "## "), prefix) {
				inSection = true
				out = append(out, line)
			}
			continue
		}
		if inSection {
			out = append(out, line)
		}
	}

	if !inSection {
		return "", fmt.Errorf("section %q not found in today.md", section)
	}
	return strings.TrimRight(strings.Join(out, "\n"), "\n"), nil
}
