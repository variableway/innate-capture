package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/variableway/innate/capture/internal/config"
	"github.com/variableway/innate/capture/internal/platform/fsrepo"
)

type workspaceService struct {
	fs fsrepo.FSRepo
}

var defaultService Service = NewService(fsrepo.NewOSFSRepo())

func NewService(fs fsrepo.FSRepo) Service {
	return workspaceService{fs: fs}
}

func DefaultService() Service {
	return defaultService
}

func (s workspaceService) Root(cfg *config.Config) string {
	if root := os.Getenv("CAPTURE_WORKSPACE_ROOT"); root != "" {
		return root
	}
	return cfg.Workspace.Root
}

func (s workspaceService) InboxDir(cfg *config.Config) string {
	return filepath.Join(s.Root(cfg), cfg.Workspace.IdeasInbox)
}

func (s workspaceService) DailyPath(cfg *config.Config) string {
	return filepath.Join(s.Root(cfg), cfg.Workspace.DailyToday)
}

func (s workspaceService) Resolve(cfg *config.Config) Paths {
	root := s.Root(cfg)
	return Paths{
		Root:     root,
		IdeasDir: filepath.Join(root, "ideas"),
		DailyDir: filepath.Join(root, "daily"),
		InboxDir: filepath.Join(root, cfg.Workspace.IdeasInbox),
		TodayMD:  filepath.Join(root, cfg.Workspace.DailyToday),
	}
}

func (s workspaceService) Validate(cfg *config.Config) error {
	paths := s.Resolve(cfg)
	if paths.Root == "" {
		return fmt.Errorf("workspace root is not configured; set workspace.root in config or CAPTURE_WORKSPACE_ROOT env")
	}

	if _, err := s.fs.Stat(paths.Root); os.IsNotExist(err) {
		return fmt.Errorf("workspace root does not exist: %s", paths.Root)
	}

	if _, err := s.fs.Stat(paths.IdeasDir); os.IsNotExist(err) {
		return fmt.Errorf("missing ideas/ directory under workspace root: %s", paths.IdeasDir)
	}

	if _, err := s.fs.Stat(paths.DailyDir); os.IsNotExist(err) {
		return fmt.Errorf("missing daily/ directory under workspace root: %s", paths.DailyDir)
	}

	return nil
}
