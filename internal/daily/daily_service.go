package daily

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/variableway/innate/capture/internal/config"
	"github.com/variableway/innate/capture/internal/platform/fsrepo"
	"github.com/variableway/innate/capture/internal/workspace"
)

type dailyService struct {
	fs       fsrepo.FSRepo
	ws       workspace.Service
	sections map[string]string
}

var defaultService Service = NewService(fsrepo.NewOSFSRepo(), workspace.DefaultService())

func NewService(fs fsrepo.FSRepo, ws workspace.Service) Service {
	return dailyService{
		fs: fs,
		ws: ws,
		sections: map[string]string{
			"input":  "输入",
			"output": "输出",
			"ideas":  "Ideas",
		},
	}
}

func DefaultService() Service {
	return defaultService
}

func (s dailyService) ValidSections() []string {
	return []string{"input", "output", "ideas"}
}

func (s dailyService) IsValidSection(name string) bool {
	_, ok := s.sections[name]
	return ok
}

func (s dailyService) Read(cfg *config.Config) (string, error) {
	return s.ReadForDate(cfg, "")
}

func (s dailyService) ReadForDate(cfg *config.Config, date string) (string, error) {
	if err := s.ws.Validate(cfg); err != nil {
		return "", err
	}

	path, err := s.DailyFilePath(cfg, date)
	if err != nil {
		return "", err
	}
	data, err := s.fs.ReadFile(path)
	if err == nil {
		return string(data), nil
	}
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("read today: %w", err)
	}

	if err := s.BootstrapFromTemplateForDate(cfg, date); err != nil {
		return "", err
	}

	data, err = s.fs.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read today after bootstrap: %w", err)
	}
	return string(data), nil
}

func (s dailyService) BootstrapFromTemplate(cfg *config.Config) error {
	return s.BootstrapFromTemplateForDate(cfg, "")
}

func (s dailyService) BootstrapFromTemplateForDate(cfg *config.Config, date string) error {
	if err := s.ws.Validate(cfg); err != nil {
		return err
	}

	tplPath := filepath.Join(s.ws.Root(cfg), "daily", "_template", "day.md")
	tpl, err := s.fs.ReadFile(tplPath)
	if err != nil {
		return fmt.Errorf("read daily template: %w", err)
	}

	today, err := resolveDate(date)
	if err != nil {
		return err
	}
	content := strings.ReplaceAll(string(tpl), "__DATE__", today)

	path, err := s.DailyFilePath(cfg, today)
	if err != nil {
		return err
	}
	if err := s.fs.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create daily dir: %w", err)
	}
	if err := s.fs.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write today: %w", err)
	}
	return nil
}

func (s dailyService) PrintSection(cfg *config.Config, section string) (string, error) {
	return s.PrintSectionForDate(cfg, section, "")
}

func (s dailyService) PrintSectionForDate(cfg *config.Config, section, date string) (string, error) {
	if !s.IsValidSection(section) {
		return "", fmt.Errorf("invalid section %q; valid: %s", section, strings.Join(s.ValidSections(), ", "))
	}

	content, err := s.ReadForDate(cfg, date)
	if err != nil {
		return "", err
	}

	prefix := s.sections[section]
	lines := strings.Split(content, "\n")
	var out []string
	inSection := false
	for _, line := range lines {
		isH2 := strings.HasPrefix(line, "## ")
		if isH2 {
			if inSection {
				break
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

func (s dailyService) DailyFilePath(cfg *config.Config, date string) (string, error) {
	d, err := resolveDate(date)
	if err != nil {
		return "", err
	}
	dailyDir := filepath.Dir(s.ws.DailyPath(cfg))
	return filepath.Join(dailyDir, d+".md"), nil
}

func resolveDate(date string) (string, error) {
	if strings.TrimSpace(date) == "" {
		return time.Now().Format("2006-01-02"), nil
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return "", fmt.Errorf("invalid date %q, expected YYYY-MM-DD", date)
	}
	return date, nil
}
