package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/variableway/innate/capture/internal/model"
)

func TestRoot_EnvVar(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Workspace.Root = "/default/root"

	os.Setenv("CAPTURE_WORKSPACE_ROOT", "/env/root")
	defer os.Unsetenv("CAPTURE_WORKSPACE_ROOT")

	got := Root(cfg)
	if got != "/env/root" {
		t.Errorf("expected env root /env/root, got %s", got)
	}
}

func TestRoot_ConfigFallback(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Workspace.Root = "/fallback/root"

	got := Root(cfg)
	if got != "/fallback/root" {
		t.Errorf("expected fallback root /fallback/root, got %s", got)
	}
}

func TestInboxDir(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Workspace.Root = "/workspace"
	cfg.Workspace.IdeasInbox = "ideas/inbox"

	got := InboxDir(cfg)
	expected := filepath.Join("/workspace", "ideas", "inbox")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestDailyPath(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Workspace.Root = "/workspace"
	cfg.Workspace.DailyToday = "daily/today.md"

	got := DailyPath(cfg)
	expected := filepath.Join("/workspace", "daily", "today.md")
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestValidate_EmptyRoot(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Workspace.Root = ""

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for empty root")
	}
}

func TestValidate_NonexistentRoot(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Workspace.Root = "/nonexistent/path/xyz"

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for nonexistent root")
	}
}

func TestValidate_MissingIdeas(t *testing.T) {
	dir := t.TempDir()
	dailyDir := filepath.Join(dir, "daily")
	os.MkdirAll(dailyDir, 0755)

	cfg := model.DefaultConfig()
	cfg.Workspace.Root = dir

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for missing ideas/")
	}
}

func TestValidate_MissingDaily(t *testing.T) {
	dir := t.TempDir()
	ideasDir := filepath.Join(dir, "ideas")
	os.MkdirAll(ideasDir, 0755)

	cfg := model.DefaultConfig()
	cfg.Workspace.Root = dir

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for missing daily/")
	}
}

func TestValidate_Success(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "ideas"), 0755)
	os.MkdirAll(filepath.Join(dir, "daily"), 0755)

	cfg := model.DefaultConfig()
	cfg.Workspace.Root = dir

	err := Validate(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInboxDir_EnvOverride(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Workspace.Root = "/default/root"

	os.Setenv("CAPTURE_WORKSPACE_ROOT", "/env/root")
	defer os.Unsetenv("CAPTURE_WORKSPACE_ROOT")

	got := InboxDir(cfg)
	expected := filepath.Join("/env/root", cfg.Workspace.IdeasInbox)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}

func TestDailyPath_EnvOverride(t *testing.T) {
	cfg := model.DefaultConfig()
	cfg.Workspace.Root = "/default/root"

	os.Setenv("CAPTURE_WORKSPACE_ROOT", "/env/root")
	defer os.Unsetenv("CAPTURE_WORKSPACE_ROOT")

	got := DailyPath(cfg)
	expected := filepath.Join("/env/root", cfg.Workspace.DailyToday)
	if got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}
}
