package daily

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/variableway/innate/capture/internal/config"
)

const testTemplate = `# Daily: __DATE__

## 输入（今天读 / 收）

- [ ] read CONTEXT

## 输出（今天必须交付）

- [ ] write daily-read

## Ideas 焦点

| Idea | 阶段 | 今日一步 |
|---|---|---|
| test-idea | inbox | step |

## 备注

-
`

// setupWorkspace creates a temp innate-works workspace with ideas/ and
// daily/_template/day.md, returning a config pointing at it.
func setupWorkspace(t *testing.T) *config.Config {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "ideas"), 0755)
	os.MkdirAll(filepath.Join(root, "daily", "_template"), 0755)
	if err := os.WriteFile(filepath.Join(root, "daily", "_template", "day.md"), []byte(testTemplate), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	cfg := config.DefaultConfig()
	cfg.Workspace.Root = root
	return cfg
}

func TestRead_ExistingFile(t *testing.T) {
	cfg := setupWorkspace(t)
	path := filepath.Join(cfg.Workspace.Root, "daily", "today.md")
	os.WriteFile(path, []byte("# Daily: 2026-07-06\n\nbody\n"), 0644)

	got, err := Read(cfg)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !strings.Contains(got, "# Daily: 2026-07-06") || !strings.Contains(got, "body") {
		t.Errorf("unexpected content: %q", got)
	}
}

func TestRead_BootstrapWhenMissing(t *testing.T) {
	cfg := setupWorkspace(t)
	today := filepath.Join(cfg.Workspace.Root, "daily", "today.md")
	if _, err := os.Stat(today); !os.IsNotExist(err) {
		t.Fatalf("today.md should not exist yet")
	}

	got, err := Read(cfg)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	todayStr := time.Now().Format("2006-01-02")
	if !strings.Contains(got, "# Daily: "+todayStr) {
		t.Errorf("expected bootstrapped date %q in:\n%s", todayStr, got)
	}
	if strings.Contains(got, "__DATE__") {
		t.Errorf("__DATE__ not substituted:\n%s", got)
	}
	if _, err := os.Stat(today); err != nil {
		t.Errorf("today.md not created: %v", err)
	}
}

func TestBootstrapFromTemplate(t *testing.T) {
	cfg := setupWorkspace(t)
	if err := BootstrapFromTemplate(cfg); err != nil {
		t.Fatalf("BootstrapFromTemplate: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(cfg.Workspace.Root, "daily", "today.md"))
	if err != nil {
		t.Fatalf("read today: %v", err)
	}
	s := string(data)
	todayStr := time.Now().Format("2006-01-02")
	if !strings.Contains(s, "# Daily: "+todayStr) {
		t.Errorf("expected date %q, got:\n%s", todayStr, s)
	}
	if strings.Contains(s, "__DATE__") {
		t.Errorf("__DATE__ not substituted")
	}
}

func TestBootstrapFromTemplate_NoTemplate(t *testing.T) {
	cfg := setupWorkspace(t)
	os.Remove(filepath.Join(cfg.Workspace.Root, "daily", "_template", "day.md"))

	err := BootstrapFromTemplate(cfg)
	if err == nil {
		t.Fatal("expected error when template missing")
	}
}

func TestPrintSection(t *testing.T) {
	cfg := setupWorkspace(t)
	// Pre-seed today.md so PrintSection reads exactly what we wrote.
	today := []byte(strings.ReplaceAll(testTemplate, "__DATE__", "2026-07-06"))
	os.WriteFile(filepath.Join(cfg.Workspace.Root, "daily", "today.md"), today, 0644)

	tests := []struct {
		section string
		wantHas []string
		wantNot []string
	}{
		{
			section: "input",
			wantHas: []string{"## 输入", "read CONTEXT"},
			wantNot: []string{"## 输出", "## Ideas", "## 备注"},
		},
		{
			section: "output",
			wantHas: []string{"## 输出", "write daily-read"},
			wantNot: []string{"## 输入", "## Ideas"},
		},
		{
			section: "ideas",
			wantHas: []string{"## Ideas 焦点", "test-idea"},
			wantNot: []string{"## 输入", "## 输出", "## 备注"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.section, func(t *testing.T) {
			got, err := PrintSection(cfg, tt.section)
			if err != nil {
				t.Fatalf("PrintSection(%q): %v", tt.section, err)
			}
			for _, w := range tt.wantHas {
				if !strings.Contains(got, w) {
					t.Errorf("section %q missing %q in:\n%s", tt.section, w, got)
				}
			}
			for _, w := range tt.wantNot {
				if strings.Contains(got, w) {
					t.Errorf("section %q should not contain %q in:\n%s", tt.section, w, got)
				}
			}
		})
	}
}

func TestPrintSection_Invalid(t *testing.T) {
	cfg := setupWorkspace(t)
	if _, err := PrintSection(cfg, "bogus"); err == nil {
		t.Fatal("expected error for invalid section")
	}
}

func TestPrintSection_NotFound(t *testing.T) {
	cfg := setupWorkspace(t)
	// today.md whose only H2 does not match any known section.
	os.WriteFile(filepath.Join(cfg.Workspace.Root, "daily", "today.md"), []byte("# Daily: x\n\n## 其他\n\nbody\n"), 0644)

	if _, err := PrintSection(cfg, "ideas"); err == nil {
		t.Fatal("expected error when section heading absent")
	}
}

func TestIsValidSection(t *testing.T) {
	for _, s := range ValidSections() {
		if !IsValidSection(s) {
			t.Errorf("IsValidSection(%q) = false, want true", s)
		}
	}
	if IsValidSection("nope") {
		t.Errorf("IsValidSection(nope) = true, want false")
	}
}
