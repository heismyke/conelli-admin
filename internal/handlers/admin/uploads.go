package admin

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/conelli/admin-backend/config"
	"github.com/gin-gonic/gin"
)

type presignUploadRequest struct {
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
	Folder      string `json:"folder"`
}

type presignUploadResponse struct {
	Key       string `json:"key"`
	Method    string `json:"method"`
	UploadURL string `json:"uploadUrl"`
	FileURL   string `json:"fileUrl"`
	ExpiresIn int64  `json:"expiresIn"`
}

func (h *Handler) PresignUpload(c *gin.Context) {
	var payload presignUploadRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		badRequest(c, err.Error())
		return
	}
	if strings.TrimSpace(payload.FileName) == "" {
		badRequest(c, "fileName is required")
		return
	}
	if strings.TrimSpace(config.Envs.AWS_S3_BUCKET) == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AWS_S3_BUCKET is not configured"})
		return
	}

	contentType := strings.TrimSpace(payload.ContentType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	key := uploadKey(payload.Folder, payload.FileName)
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(config.Envs.AWS_REGION))
	if err != nil {
		serverError(c, err)
		return
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint := strings.TrimSpace(config.Envs.AWS_S3_ENDPOINT); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}
	})
	presigner := s3.NewPresignClient(client)
	expires := 15 * time.Minute
	result, err := presigner.PresignPutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(config.Envs.AWS_S3_BUCKET),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		serverError(c, err)
		return
	}

	c.JSON(http.StatusOK, presignUploadResponse{
		Key:       key,
		Method:    http.MethodPut,
		UploadURL: result.URL,
		FileURL:   publicFileURL(key),
		ExpiresIn: int64(expires.Seconds()),
	})
}

func uploadKey(folder string, fileName string) string {
	cleanName := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '.' || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, path.Base(fileName))
	if cleanName == "." || cleanName == "/" || cleanName == "" {
		cleanName = "upload"
	}

	parts := []string{strings.Trim(config.Envs.AWS_S3_PREFIX, "/")}
	if folder = strings.Trim(folder, "/"); folder != "" {
		parts = append(parts, folder)
	}
	parts = append(parts, time.Now().UTC().Format("2006/01/02"), fmt.Sprintf("%d-%s", time.Now().UTC().UnixNano(), cleanName))
	return path.Join(parts...)
}

func publicFileURL(key string) string {
	if base := strings.TrimRight(config.Envs.AWS_S3_PUBLIC_URL, "/"); base != "" {
		return base + "/" + key
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", config.Envs.AWS_S3_BUCKET, config.Envs.AWS_REGION, key)
}
