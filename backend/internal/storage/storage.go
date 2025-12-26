package storage

import (
	"io"
	"mime/multipart"
)

// StorageDriver 存储驱动接口
type StorageDriver interface {
	// Upload 上传文件
	// file: 文件内容
	// filename: 原始文件名
	// 返回: 文件访问URL, 文件key, 错误
	Upload(file multipart.File, header *multipart.FileHeader) (url string, key string, err error)

	// UploadReader 从 Reader 上传文件
	UploadReader(reader io.Reader, filename string, size int64) (url string, key string, err error)

	// Delete 删除文件
	Delete(key string) error

	// GetURL 获取文件访问URL
	GetURL(key string) string
}

// UploadResult 上传结果
type UploadResult struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
}
