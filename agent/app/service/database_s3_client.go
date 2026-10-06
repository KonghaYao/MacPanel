package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	s3DefaultTimeout  = 30 * time.Second
	s3TransferTimeout = 24 * time.Hour
)

type s3Runtime struct {
	client *minio.Client
}

func newS3Runtime(endpoint, accessKey, secretKey, region string, useSSL, pathStyle bool) (*s3Runtime, error) {
	host, ssl, err := splitS3Endpoint(endpoint, useSSL)
	if err != nil {
		return nil, err
	}
	lookup := minio.BucketLookupAuto
	if pathStyle {
		lookup = minio.BucketLookupPath
	}
	client, err := minio.New(host, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:       ssl,
		Region:       region,
		BucketLookup: lookup,
		Transport:    &http.Transport{},
	})
	if err != nil {
		return nil, err
	}
	return &s3Runtime{client: client}, nil
}

func (s *s3Runtime) ctx(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func (s *s3Runtime) listBuckets() ([]minio.BucketInfo, error) {
	ctx, cancel := s.ctx(s3DefaultTimeout)
	defer cancel()
	return s.client.ListBuckets(ctx)
}

func (s *s3Runtime) makeBucket(name, region string) error {
	ctx, cancel := s.ctx(s3DefaultTimeout)
	defer cancel()
	return s.client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: region})
}

func (s *s3Runtime) removeBucket(name string) error {
	ctx, cancel := s.ctx(s3DefaultTimeout)
	defer cancel()
	return s.client.RemoveBucket(ctx, name)
}

func (s *s3Runtime) listObjects(bucket, prefix string) ([]minio.ObjectInfo, error) {
	opts := minio.ListObjectsOptions{
		Recursive: false,
		Prefix:    normalizeS3Prefix(prefix),
	}
	ctx, cancel := s.ctx(s3DefaultTimeout)
	defer cancel()
	var result []minio.ObjectInfo
	for object := range s.client.ListObjects(ctx, bucket, opts) {
		if object.Err != nil {
			return nil, object.Err
		}
		result = append(result, object)
	}
	return result, nil
}

func (s *s3Runtime) listObjectsRecursive(bucket, prefix string) ([]minio.ObjectInfo, error) {
	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    prefix,
	}
	ctx, cancel := s.ctx(s3DefaultTimeout)
	defer cancel()
	var result []minio.ObjectInfo
	for object := range s.client.ListObjects(ctx, bucket, opts) {
		if object.Err != nil {
			return nil, object.Err
		}
		result = append(result, object)
	}
	return result, nil
}

func (s *s3Runtime) putObject(bucket, key string, reader io.Reader, size int64, contentType string) error {
	ctx, cancel := s.ctx(s3TransferTimeout)
	defer cancel()
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}
	_, err := s.client.PutObject(ctx, bucket, key, reader, size, opts)
	return err
}

func (s *s3Runtime) statObject(bucket, key string) (minio.ObjectInfo, error) {
	ctx, cancel := s.ctx(s3DefaultTimeout)
	defer cancel()
	return s.client.StatObject(ctx, bucket, key, minio.StatObjectOptions{})
}

func (s *s3Runtime) getObject(ctx context.Context, bucket, key string) (*minio.Object, error) {
	return s.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
}

func (s *s3Runtime) removeObject(bucket, key string) error {
	ctx, cancel := s.ctx(s3DefaultTimeout)
	defer cancel()
	return s.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

func (s *s3Runtime) removeObjects(bucket string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	objectsCh := make(chan minio.ObjectInfo, len(keys))
	for _, key := range keys {
		objectsCh <- minio.ObjectInfo{Key: key}
	}
	close(objectsCh)
	ctx, cancel := s.ctx(s3TransferTimeout)
	defer cancel()
	var errs []string
	for err := range s.client.RemoveObjects(ctx, bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		if err.Err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", err.ObjectName, err.Err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("delete objects failed: %s", strings.Join(errs, "; "))
	}
	return nil
}
