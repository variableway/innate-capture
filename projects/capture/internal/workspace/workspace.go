package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/variableway/innate/capture/internal/model"
)

func Root(cfg *model.Config) string {
	if root := os.Getenv("CAPTURE_WORKSPACE_ROOT"); root != "" {
		return root
	}
	return cfg.Workspace.Root
}

func InboxDir(cfg *model.Config) string {
	return filepath.Join(Root(cfg), cfg.Workspace.IdeasInbox)
}

func DailyPath(cfg *model.Config) string {
	return filepath.Join(Root(cfg), cfg.Workspace.DailyToday)
}

func Validate(cfg *model.Config) error {
	root := Root(cfg)
	if root == "" {
		return fmt.Errorf("workspace root is not configured; set workspace.root in config or CAPTURE_WORKSPACE_ROOT env")
	}

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return fmt.Errorf("workspace root does not exist: %s", root)
	}

	ideasDir := filepath.Join(root, "ideas")
	if _, err := os.Stat(ideasDir); os.IsNotExist(err) {
		return fmt.Errorf("missing ideas/ directory under workspace root: %s", ideasDir)
	}

	dailyDir := filepath.Join(root, "daily")
	if _, err := os.Stat(dailyDir); os.IsNotExist(err) {
		return fmt.Errorf("missing daily/ directory under workspace root: %s", dailyDir)
	}

	return nil
}
