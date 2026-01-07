package files

import (
	"context"
	"io"
)

// FileService 提供业务层的便捷方法，组合 Storage 与 Processor
type FileService struct {
	Store Storage
}

func NewService(s Storage) *FileService {
	return &FileService{Store: s}
}

// Upload 保存文件并返回 id 与 metadata
func (s *FileService) Upload(ctx context.Context, r io.Reader, meta Metadata) (string, Metadata, error) {
	id, err := s.Store.Save(ctx, r, meta)
	if err != nil {
		return "", Metadata{}, err
	}
	m, err := s.Store.Stat(ctx, id)
	if err != nil {
		return id, Metadata{}, err
	}
	return id, m, nil
}

// Get 用于获取文件和元信息
func (s *FileService) Get(ctx context.Context, id string) (io.ReadCloser, Metadata, error) {
	return s.Store.Get(ctx, id)
}

// Delete 删除文件
func (s *FileService) Delete(ctx context.Context, id string) error {
	return s.Store.Delete(ctx, id)
}
