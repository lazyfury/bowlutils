package files

import (
	"context"
	"io"
	"time"
)

// Storage 抽象所有二进制存储后端（本地、S3 等）
// 实现必须保证对同一 ID 的并发访问行为合理（例如读取与删除的顺序一致性由调用方保证）。
type Storage interface {
	// Save 将流保存到存储后端，返回生成的文件 ID。
	Save(ctx context.Context, r io.Reader, meta Metadata) (id string, err error)
	// Get 返回文件内容的 ReadCloser 与对应 Metadata。调用者负责关闭 ReadCloser。
	Get(ctx context.Context, id string) (rc io.ReadCloser, meta Metadata, err error)
	// Delete 删除资源
	Delete(ctx context.Context, id string) error
	// Stat 返回 Metadata
	Stat(ctx context.Context, id string) (Metadata, error)
	// URL 返回该文件的可访问链接（可以是签名 URL 或公开 URL）
	URL(ctx context.Context, id string, opts URLOptions) (string, error)
}

// Processor 提供文件处理能力（缩略、转码、扫描等）
// 实现可以是同步或异步（异步可返回 task ID 或将结果写入 Storage）。
type Processor interface {
	Process(ctx context.Context, in io.Reader, meta Metadata, task ProcessTask) (resultID string, resultMeta Metadata, err error)
}

// Metadata 为文件的元信息，建议同时持久化在 DB 中以便查询。
type Metadata struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Size        int64             `json:"size"`
	ContentType string            `json:"content_type"`
	OwnerID     string            `json:"owner_id"`
	CreatedAt   time.Time         `json:"created_at"`
	Extra       map[string]string `json:"extra,omitempty"`
}

// FileType 用于描述文件类别，便于选择处理器
type FileType string

const (
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
	FileTypeDocument FileType = "document"
	FileTypeAudio    FileType = "audio"
	FileTypeBlob     FileType = "blob"
)

// ProcessTask 描述处理任务的类型与参数
type ProcessTask struct {
	Name   string // e.g. "thumbnail", "transcode"
	Params map[string]string
}

// URLOptions 为生成 URL 的可选参数（例如签名过期时间）
type URLOptions struct {
	ExpiresInSeconds int64
}
