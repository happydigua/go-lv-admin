package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"go-lv-vue-admin/internal/config"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// COSDriver 腾讯云COS存储驱动
type COSDriver struct {
	config config.COSStorage
	client *cos.Client
}

// NewCOSDriver 创建COS存储驱动
func NewCOSDriver(cfg config.COSStorage) (*COSDriver, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.Bucket, cfg.Region))
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	return &COSDriver{
		config: cfg,
		client: client,
	}, nil
}

// Upload 上传文件到COS
func (d *COSDriver) Upload(file multipart.File, header *multipart.FileHeader) (string, string, error) {
	return d.UploadReader(file, header.Filename, header.Size)
}

// UploadReader 从Reader上传文件到COS
func (d *COSDriver) UploadReader(reader io.Reader, filename string, size int64) (string, string, error) {
	ext := filepath.Ext(filename)
	key := d.generateKey(filename, ext)

	_, err := d.client.Object.Put(context.Background(), key, reader, nil)
	if err != nil {
		return "", "", fmt.Errorf("上传到COS失败: %w", err)
	}

	url := d.GetURL(key)
	return url, key, nil
}

// Delete 从COS删除文件
func (d *COSDriver) Delete(key string) error {
	_, err := d.client.Object.Delete(context.Background(), key)
	return err
}

// GetURL 获取COS文件访问URL
func (d *COSDriver) GetURL(key string) string {
	if d.config.Domain != "" {
		return fmt.Sprintf("%s/%s", d.config.Domain, key)
	}
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s", d.config.Bucket, d.config.Region, key)
}

// generateKey 生成唯一文件key
func (d *COSDriver) generateKey(filename, ext string) string {
	dateDir := time.Now().Format("2006/01/02")
	hash := md5.New()
	hash.Write([]byte(filename + time.Now().String()))
	md5Str := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("%s/%s%s", dateDir, md5Str[:16], ext)
}
