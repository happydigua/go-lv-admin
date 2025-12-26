package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/storage"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadApi struct{}

// 允许的文件类型
var (
	imageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".ico"}
	fileExtensions  = []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".zip", ".rar"}
)

// UploadImage 上传图片
func (u *UploadApi) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "请选择文件"})
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !isAllowedExt(ext, imageExtensions) {
		c.JSON(400, gin.H{"code": 7, "msg": "文件类型不支持，仅支持图片: " + strings.Join(imageExtensions, ", ")})
		return
	}

	// 验证文件大小 (5MB)
	if header.Size > 5*1024*1024 {
		c.JSON(400, gin.H{"code": 7, "msg": "文件大小超过限制(5MB)"})
		return
	}

	result, err := upload(file, header)
	if err != nil {
		global.LV_LOG.Error("上传图片失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "上传失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": result,
		"msg":  "上传成功",
	})
}

// UploadFile 上传文件
func (u *UploadApi) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "请选择文件"})
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allExtensions := append(imageExtensions, fileExtensions...)
	if !isAllowedExt(ext, allExtensions) {
		c.JSON(400, gin.H{"code": 7, "msg": "文件类型不支持"})
		return
	}

	// 验证文件大小 (20MB)
	if header.Size > 20*1024*1024 {
		c.JSON(400, gin.H{"code": 7, "msg": "文件大小超过限制(20MB)"})
		return
	}

	result, err := upload(file, header)
	if err != nil {
		global.LV_LOG.Error("上传文件失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "上传失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": result,
		"msg":  "上传成功",
	})
}

// DeleteFile 删除文件
func (u *UploadApi) DeleteFile(c *gin.Context) {
	var req struct {
		Key string `json:"key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": "参数错误"})
		return
	}

	driver := storage.GetDriver()
	if err := driver.Delete(req.Key); err != nil {
		global.LV_LOG.Error("删除文件失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "删除失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
}

// 执行上传
func upload(file multipart.File, header *multipart.FileHeader) (*storage.UploadResult, error) {
	driver := storage.GetDriver()
	url, key, err := driver.Upload(file, header)
	if err != nil {
		return nil, err
	}

	return &storage.UploadResult{
		URL:      url,
		Key:      key,
		Filename: header.Filename,
		Size:     header.Size,
		MimeType: header.Header.Get("Content-Type"),
	}, nil
}

// 检查文件扩展名是否允许
func isAllowedExt(ext string, allowed []string) bool {
	for _, a := range allowed {
		if ext == a {
			return true
		}
	}
	return false
}
