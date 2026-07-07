package idea

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/variableway/innate/capture/internal/config"
	"github.com/variableway/innate/capture/internal/platform/fsrepo"
	"github.com/variableway/innate/capture/internal/workspace"
)

type ideaService struct {
	fs fsrepo.FSRepo
	ws workspace.Service
}

var defaultService Service = NewService(fsrepo.NewOSFSRepo(), workspace.DefaultService())

func NewService(fs fsrepo.FSRepo, ws workspace.Service) Service {
	return ideaService{
		fs: fs,
		ws: ws,
	}
}

func DefaultService() Service {
	return defaultService
}

func (s ideaService) Write(cfg *config.Config, in CreateInput) (string, error) {
	if err := s.ws.Validate(cfg); err != nil {
		return "", err
	}

	inboxDir := s.ws.InboxDir(cfg)
	if err := s.fs.MkdirAll(inboxDir, 0o755); err != nil {
		return "", fmt.Errorf("create inbox dir: %w", err)
	}

	date := time.Now().Format("2006-01-02")
	slug := Slug(in.Title)
	path, err := s.resolvePath(inboxDir, date, slug)
	if err != nil {
		return "", err
	}

	if in.Source == "" {
		in.Source = SourceCLI
	}

	content := render(in.Title, date, in.Source, in.Description, in.Context)
	if err := s.fs.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write inbox file: %w", err)
	}

	return path, nil
}

func (s ideaService) List(cfg *config.Config) ([]Entry, error) {
	if err := s.ws.Validate(cfg); err != nil {
		return nil, err
	}

	inboxDir := s.ws.InboxDir(cfg)
	entries, err := s.fs.ReadDir(inboxDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read inbox dir: %w", err)
	}

	var result []Entry
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		path := filepath.Join(inboxDir, e.Name())
		title, date := s.parseEntryMeta(path, e.Name())
		result = append(result, Entry{
			Path:     path,
			Filename: e.Name(),
			Title:    title,
			Date:     date,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Filename > result[j].Filename
	})

	return result, nil
}

func (s ideaService) resolvePath(inboxDir, date, slug string) (string, error) {
	base := fmt.Sprintf("%s-%s.md", date, slug)
	path := filepath.Join(inboxDir, base)
	if _, err := s.fs.Stat(path); os.IsNotExist(err) {
		return path, nil
	}

	for i := 2; i < 1000; i++ {
		candidate := filepath.Join(inboxDir, fmt.Sprintf("%s-%s-%d.md", date, slug, i))
		if _, err := s.fs.Stat(candidate); os.IsNotExist(err) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("too many slug conflicts for %s-%s", date, slug)
}

func render(title, date, source, description, context string) string {
	if strings.TrimSpace(description) == "" {
		description = title
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", title)
	fmt.Fprintf(&b, "**Captured:** %s\n", date)
	fmt.Fprintf(&b, "**Source:** %s\n", source)
	fmt.Fprintln(&b, "**Stage:** inbox")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## 一句话")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, description)
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## 原始上下文")
	fmt.Fprintln(&b)
	if strings.TrimSpace(context) != "" {
		fmt.Fprintln(&b, context)
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## 初步问题")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "- [ ]")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## 下一步")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "- [ ] 留在 inbox 观察")
	fmt.Fprintln(&b, "- [ ] 晋升到 exploring/")

	return b.String()
}

func parseDateFromFilename(filename string) string {
	name := strings.TrimSuffix(filename, ".md")
	if len(name) >= 10 && name[4] == '-' && name[7] == '-' {
		return name[:10]
	}
	return ""
}

func (s ideaService) parseEntryMeta(path, filename string) (title, date string) {
	date = parseDateFromFilename(filename)

	data, err := s.fs.ReadFile(path)
	if err != nil {
		return filename, date
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "# ") {
		return strings.TrimPrefix(lines[0], "# "), date
	}
	return filename, date
}
