package infra

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/huodaoshi/harness/backend/conf"
)

// ObjectStorage is the interface for uploading and downloading binary objects.
type ObjectStorage interface {
	Upload(ctx context.Context, key string, data []byte) error
	Download(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}

// LocalStorage stores objects as files under a local directory.
type LocalStorage struct {
	dir string
}

// Upload writes data to a file at <dir>/<key>, creating parent directories as needed.
func (s *LocalStorage) Upload(_ context.Context, key string, data []byte) error {
	fullPath := filepath.Join(s.dir, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("infra: local storage: mkdir: %w", err)
	}
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return fmt.Errorf("infra: local storage: write: %w", err)
	}
	return nil
}

// Download reads and returns the file at <dir>/<key>.
func (s *LocalStorage) Download(_ context.Context, key string) ([]byte, error) {
	fullPath := filepath.Join(s.dir, key)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("infra: local storage: read: %w", err)
	}
	return data, nil
}

// Delete removes the file at <dir>/<key>. Missing files are ignored.
func (s *LocalStorage) Delete(_ context.Context, key string) error {
	fullPath := filepath.Join(s.dir, key)
	if err := os.Remove(fullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("infra: local storage: delete: %w", err)
	}
	return nil
}

// COSStorage is a placeholder for the Tencent COS-backed object storage.
type COSStorage struct{}

// Upload is not yet implemented for COS.
func (s *COSStorage) Upload(_ context.Context, _ string, _ []byte) error {
	return errors.New("cos: not implemented")
}

// Download is not yet implemented for COS.
func (s *COSStorage) Download(_ context.Context, _ string) ([]byte, error) {
	return nil, errors.New("cos: not implemented")
}

// Delete is not yet implemented for COS.
func (s *COSStorage) Delete(_ context.Context, _ string) error {
	return errors.New("cos: delete not implemented")
}

// NewObjectStorage creates an ObjectStorage based on cfg.Provider.
// "local" returns a LocalStorage writing under tmp/cos/.
// "cos" returns a COSStorage skeleton.
func NewObjectStorage(cfg conf.COSConfig) (ObjectStorage, error) {
	switch cfg.Provider {
	case "local":
		return &LocalStorage{dir: "tmp/cos/"}, nil
	case "cos":
		return &COSStorage{}, nil
	default:
		return nil, fmt.Errorf("infra: cos: unknown provider %q", cfg.Provider)
	}
}
