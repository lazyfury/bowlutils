package files

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("file not found")
)

// LocalStorage 是简单的基于文件系统的 Storage 实现
// 它在 basePath 下保存二进制文件和对应的 metadata（JSON sidecar）。
type LocalStorage struct {
	basePath      string
	publicBaseURL string // 可选：用于生成公开访问 URL
}

func NewLocalStorage(basePath string, publicBaseURL string) (*LocalStorage, error) {
	if basePath == "" {
		return nil, fmt.Errorf("basePath required")
	}
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{basePath: basePath, publicBaseURL: publicBaseURL}, nil
}

func (s *LocalStorage) filePath(id string) string {
	return filepath.Join(s.basePath, id)
}

func (s *LocalStorage) metaPath(id string) string {
	return filepath.Join(s.basePath, id+".meta.json")
}

func (s *LocalStorage) Save(ctx context.Context, r io.Reader, meta Metadata) (string, error) {
	id := uuid.New().String()
	fp := s.filePath(id)
	f, err := os.Create(fp)
	if err != nil {
		return "", err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return "", err
	}
	meta.ID = id
	meta.Size = n
	meta.CreatedAt = time.Now()

	b, _ := json.MarshalIndent(meta, "", "  ")
	if err := os.WriteFile(s.metaPath(id), b, 0644); err != nil {
		// attempt cleanup
		_ = os.Remove(fp)
		return "", err
	}
	return id, nil
}

func (s *LocalStorage) Get(ctx context.Context, id string) (io.ReadCloser, Metadata, error) {
	fp := s.filePath(id)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		return nil, Metadata{}, ErrNotFound
	}
	f, err := os.Open(fp)
	if err != nil {
		return nil, Metadata{}, err
	}
	meta, err := s.loadMeta(id)
	if err != nil {
		f.Close()
		return nil, Metadata{}, err
	}
	return f, meta, nil
}

func (s *LocalStorage) Delete(ctx context.Context, id string) error {
	if err := os.Remove(s.filePath(id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove(s.metaPath(id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *LocalStorage) Stat(ctx context.Context, id string) (Metadata, error) {
	if _, err := os.Stat(s.filePath(id)); os.IsNotExist(err) {
		return Metadata{}, ErrNotFound
	}
	return s.loadMeta(id)
}

func (s *LocalStorage) URL(ctx context.Context, id string, opts URLOptions) (string, error) {
	// 简单实现：如果配置了 publicBaseURL，则使用它拼接路径；否则不支持
	if s.publicBaseURL == "" {
		return "", fmt.Errorf("public URL not configured")
	}
	// 不做签名，仅返回基础拼接
	return fmt.Sprintf("%s/%s", s.publicBaseURL, id), nil
}

func (s *LocalStorage) loadMeta(id string) (Metadata, error) {
	b, err := os.ReadFile(s.metaPath(id))
	if err != nil {
		if os.IsNotExist(err) {
			return Metadata{}, ErrNotFound
		}
		return Metadata{}, err
	}
	var m Metadata
	if err := json.Unmarshal(b, &m); err != nil {
		return Metadata{}, err
	}
	return m, nil
}
