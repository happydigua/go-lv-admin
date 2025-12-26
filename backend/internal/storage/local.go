package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-lv-vue-admin/internal/config"
)

// LocalDriver 本地存储驱动
type LocalDriver struct {
	config config.LocalStorage
}

// NewLocalDriver 创建本地存储驱动
func NewLocalDriver(cfg config.LocalStorage) *LocalDriver {
	// 确保上传目录存在
	if cfg.Path == "" {
		cfg.Path = "./uploads"
	}
	os.MkdirAll(cfg.Path, 0755)

	return &LocalDriver{config: cfg}
}

// Upload 上传文件
func (d *LocalDriver) Upload(file multipart.File, header *multipart.FileHeader) (string, string, error) {
	return d.UploadReader(file, header.Filename, header.Size)
}

// UploadReader 从 Reader 上传文件
func (d *LocalDriver) UploadReader(reader io.Reader, filename string, size int64) (string, string, error) {
	// 生成唯一文件名
	ext := filepath.Ext(filename)
	key := d.generateKey(filename, ext)

	// 创建日期目录
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(d.config.Path, dateDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 完整文件路径
	fullPath := filepath.Join(fullDir, filepath.Base(key))

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, reader); err != nil {
		return "", "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建访问URL
	relativePath := filepath.Join(dateDir, filepath.Base(key))
	url := d.GetURL(relativePath)

	return url, relativePath, nil
}

// Delete 删除文件
func (d *LocalDriver) Delete(key string) error {
	fullPath := filepath.Join(d.config.Path, key)
	return os.Remove(fullPath)
}

// GetURL 获取文件访问URL
func (d *LocalDriver) GetURL(key string) string {
	domain := strings.TrimRight(d.config.Domain, "/")
	return fmt.Sprintf("%s/uploads/%s", domain, key)
}

// generateKey 生成唯一文件key
func (d *LocalDriver) generateKey(filename, ext string) string {
	hash := md5.New()
	hash.Write([]byte(filename + time.Now().String()))
	md5Str := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("%s%s", md5Str[:16], ext)
}
