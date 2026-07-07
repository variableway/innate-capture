package idea

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/variableway/innate/capture/internal/config"
)

func TestSlug(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Hello World", "hello-world"},
		{"  Test_Idea  ", "test-idea"},
		{"My Cool Idea!", "my-cool-idea"},
		{"a--b___c", "a-b-c"},
		{"...", "idea"},
		{"v1.0 release", "v1.0-release"},
	}
	for _, tt := range tests {
		if got := Slug(tt.in); got != tt.want {
			t.Errorf("Slug(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestWriteAndList(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "ideas"), 0755)
	os.MkdirAll(filepath.Join(root, "daily"), 0755)
	inbox := filepath.Join(root, "ideas", "inbox")
	os.MkdirAll(inbox, 0755)

	cfg := config.DefaultConfig()
	cfg.Workspace.Root = root

	path, err := Write(cfg, "测试想法", "一句话描述", "some context", SourceCLI)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	if !strings.HasSuffix(path, ".md") {
		t.Fatalf("unexpected path: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)
	for _, want := range []string{
		"# 测试想法",
		"**Source:** capture-cli",
		"**Stage:** inbox",
		"## 一句话",
		"一句话描述",
		"some context",
		"- [ ] 留在 inbox 观察",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("content missing %q:\n%s", want, content)
		}
	}

	// slug conflict → -2 suffix
	path2, err := Write(cfg, "测试想法", "第二篇", "", SourceCLI)
	if err != nil {
		t.Fatalf("Write conflict: %v", err)
	}
	if !strings.Contains(filepath.Base(path2), "-2.md") {
		t.Errorf("expected -2 suffix, got %s", path2)
	}

	entries, err := List(cfg)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestWrite_DefaultDescription(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "ideas"), 0755)
	os.MkdirAll(filepath.Join(root, "daily"), 0755)
	os.MkdirAll(filepath.Join(root, "ideas", "inbox"), 0755)

	cfg := config.DefaultConfig()
	cfg.Workspace.Root = root

	path, err := Write(cfg, "Only Title", "", "", SourceCLI)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(data), "Only Title\n") {
		t.Errorf("expected title as description fallback")
	}
}

func TestParseDateFromFilename(t *testing.T) {
	if got := parseDateFromFilename("2026-07-06-test-idea.md"); got != "2026-07-06" {
		t.Errorf("got %q", got)
	}
	if got := parseDateFromFilename("bad.md"); got != "" {
		t.Errorf("got %q", got)
	}
}
