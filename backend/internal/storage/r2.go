package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"go-lv-vue-admin/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// R2Driver Cloudflare R2存储驱动 (兼容S3 API)
type R2Driver struct {
	config config.R2Storage
	client *s3.S3
}

// NewR2Driver 创建R2存储驱动
func NewR2Driver(cfg config.R2Storage) (*R2Driver, error) {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)

	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("auto"),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.AccessKeySecret, ""),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("初始化R2客户端失败: %w", err)
	}

	return &R2Driver{
		config: cfg,
		client: s3.New(sess),
	}, nil
}

// Upload 上传文件到R2
func (d *R2Driver) Upload(file multipart.File, header *multipart.FileHeader) (string, string, error) {
	return d.UploadReader(file, header.Filename, header.Size)
}

// UploadReader 从Reader上传文件到R2
func (d *R2Driver) UploadReader(reader io.Reader, filename string, size int64) (string, string, error) {
	ext := filepath.Ext(filename)
	key := d.generateKey(filename, ext)

	// 读取全部内容
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", "", fmt.Errorf("读取文件内容失败: %w", err)
	}

	_, err = d.client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(d.config.Bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(content),
		ContentLength: aws.Int64(int64(len(content))),
	})
	if err != nil {
		return "", "", fmt.Errorf("上传到R2失败: %w", err)
	}

	url := d.GetURL(key)
	return url, key, nil
}

// Delete 从R2删除文件
func (d *R2Driver) Delete(key string) error {
	_, err := d.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(d.config.Bucket),
		Key:    aws.String(key),
	})
	return err
}

// GetURL 获取R2文件访问URL
func (d *R2Driver) GetURL(key string) string {
	if d.config.Domain != "" {
		return fmt.Sprintf("%s/%s", d.config.Domain, key)
	}
	// R2 公共访问需要自定义域名或开启 R2.dev
	return fmt.Sprintf("https://%s.r2.dev/%s", d.config.Bucket, key)
}

// generateKey 生成唯一文件key
func (d *R2Driver) generateKey(filename, ext string) string {
	dateDir := time.Now().Format("2006/01/02")
	hash := md5.New()
	hash.Write([]byte(filename + time.Now().String()))
	md5Str := hex.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("%s/%s%s", dateDir, md5Str[:16], ext)
}
