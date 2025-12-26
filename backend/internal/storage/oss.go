package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"go-lv-vue-admin/internal/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSDriver 阿里云OSS存储驱动
type OSSDriver struct {
	config config.OSSStorage
	client *oss.Client
	bucket *oss.Bucket
}

// NewOSSDriver 创建OSS存储驱动
func NewOSSDriver(cfg config.OSSStorage) (*OSSDriver, error) {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("初始化OSS客户端失败: %w", err)
	}

	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取OSS Bucket失败: %w", err)
	}

	return &OSSDriver{
		config: cfg,
		client: client,
		bucket: bucket,
	}, nil
}

// Upload 上传文件到OSS
func (d *OSSDriver) Upload(file multipart.File, header *multipart.FileHeader) (string, string, error) {
	return d.UploadReader(file, header.Filename, header.Size)
}

// UploadReader 从Reader上传文件到OSS
func (d *OSSDriver) UploadReader(reader io.Reader, filename string, size int64) (string, string, error) {
	ext := filepath.Ext(filename)
	key := d.generateKey(filename, ext)

	err := d.bucket.PutObject(key, reader)
	if err != nil {
		return "", "", fmt.Errorf("上传到OSS失败: %w", err)
	}

	url := d.GetURL(key)
	return url, key, nil
}

// Delete 从OSS删除文件
func (d *OSSDriver) Delete(key string) error {
	return d.bucket.DeleteObject(key)
}

// GetURL 获取OSS文件访问URL
func (d *OSSDriver) GetURL(key string) string {
	if d.config.Domain != "" {
		return fmt.Sprintf("%s/%s", d.config.Domain, key)
	}
	return fmt.Sprintf("https://%s.%s/%s", d.config.Bucket, d.config.Endpoint, key)
}

// generateKey 生成唯一文件key
func (d *OSSDriver) generateKey(filename, ext string) string {
	dateDir := time.Now().Format("2006/01/02")
	hash := md5.New()
	hash.Write([]byte(filename + time.Now().String()))
	md5Str := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("%s/%s%s", dateDir, md5Str[:16], ext)
}
