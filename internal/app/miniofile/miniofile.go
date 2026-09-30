package miniofile

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const bucketName = "drengo-age-img"

type Client struct {
	client   *minio.Client
	endpoint string
	public   string
}

func NewFromEnv() (*Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	access := os.Getenv("MINIO_ROOT_USER")
	if access == "" {
		access = "minioadmin"
	}
	secret := os.Getenv("MINIO_ROOT_PASSWORD")
	if secret == "" {
		secret = "minioadmin"
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(access, secret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucketName)
	_ = client.SetBucketPolicy(ctx, bucketName, policy)

	return &Client{
		client:   client,
		endpoint: endpoint,
		public:   fmt.Sprintf("http://%s/%s/", endpoint, bucketName),
	}, nil
}

func (c *Client) PublicURL(objectName string) string {
	return c.public + objectName
}

func (c *Client) Upload(fileHeader *multipart.FileHeader, kind string) (string, error) {
	if fileHeader == nil {
		return "", fmt.Errorf("файл не передан")
	}
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := strings.ToLower(path.Ext(fileHeader.Filename))
	if ext == "" {
		if kind == "video" {
			ext = ".mp4"
		} else {
			ext = ".jpg"
		}
	}
	objectName := fmt.Sprintf("%s_%s%s", kind, strings.ReplaceAll(uuid.NewString(), "-", ""), ext)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = c.client.PutObject(ctx, bucketName, objectName, src, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	return c.PublicURL(objectName), nil
}

func (c *Client) UploadReader(reader io.Reader, size int64, objectName, contentType string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := c.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	return c.PublicURL(objectName), nil
}

func ObjectNameFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return path.Base(u.Path)
}
