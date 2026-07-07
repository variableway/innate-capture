package idea

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/variableway/innate/capture/internal/model"
	"github.com/variableway/innate/capture/internal/workspace"
)

const SourceCLI = "capture-cli"
const SourceTUI = "capture-tui"

// Entry is a single inbox markdown file.
type Entry struct {
	Path     string
	Filename string
	Title    string
	Date     string
}

// Write creates an inbox markdown file for the given idea.
func Write(cfg *model.Config, title, description, context, source string) (string, error) {
	if err := workspace.Validate(cfg); err != nil {
		return "", err
	}

	inboxDir := workspace.InboxDir(cfg)
	if err := os.MkdirAll(inboxDir, 0755); err != nil {
		return "", fmt.Errorf("create inbox dir: %w", err)
	}

	date := time.Now().Format("2006-01-02")
	slug := Slug(title)
	path, err := resolvePath(inboxDir, date, slug)
	if err != nil {
		return "", err
	}

	if source == "" {
		source = SourceCLI
	}

	content := render(title, date, source, description, context)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write inbox file: %w", err)
	}

	return path, nil
}

func resolvePath(inboxDir, date, slug string) (string, error) {
	base := fmt.Sprintf("%s-%s.md", date, slug)
	path := filepath.Join(inboxDir, base)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path, nil
	}

	for i := 2; i < 1000; i++ {
		candidate := filepath.Join(inboxDir, fmt.Sprintf("%s-%s-%d.md", date, slug, i))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
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

// List returns inbox entries sorted by filename descending (newest date first).
func List(cfg *model.Config) ([]Entry, error) {
	if err := workspace.Validate(cfg); err != nil {
		return nil, err
	}

	inboxDir := workspace.InboxDir(cfg)
	entries, err := os.ReadDir(inboxDir)
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
		title, date := parseEntryMeta(path, e.Name())
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

func parseDateFromFilename(filename string) string {
	name := strings.TrimSuffix(filename, ".md")
	if len(name) >= 10 && name[4] == '-' && name[7] == '-' {
		return name[:10]
	}
	return ""
}

func parseEntryMeta(path, filename string) (title, date string) {
	date = parseDateFromFilename(filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return filename, date
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 && strings.HasPrefix(lines[0], "# ") {
		return strings.TrimPrefix(lines[0], "# "), date
	}
	return filename, date
}
