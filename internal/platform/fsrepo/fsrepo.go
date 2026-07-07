package fsrepo

import (
	"io/fs"
	"os"
)

// FSRepo is the shared filesystem abstraction for domain services.
type FSRepo interface {
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
	ReadDir(name string) ([]os.DirEntry, error)
	MkdirAll(path string, perm fs.FileMode) error
}

type osFSRepo struct{}

func NewOSFSRepo() FSRepo {
	return osFSRepo{}
}

func (osFSRepo) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (osFSRepo) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (osFSRepo) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (osFSRepo) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (osFSRepo) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}
