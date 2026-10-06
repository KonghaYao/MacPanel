package v2

import (
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/1Panel-dev/1Panel/agent/app/api/v2/helper"
	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) ListS3RustFS(c *gin.Context) {
	data, err := s3Service.ListRustFS()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

func (b *BaseApi) ConnectS3RustFS(c *gin.Context) {
	var req dto.S3ConnectRustFS
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	data, err := s3Service.ConnectRustFS(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

func (b *BaseApi) ListS3Connections(c *gin.Context) {
	data, err := s3Service.ListConnections()
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

func decodeS3Secret(reqSecret string) (string, error) {
	if reqSecret == "" {
		return "", nil
	}
	secret, err := base64.StdEncoding.DecodeString(reqSecret)
	if err != nil {
		return "", err
	}
	return string(secret), nil
}

func (b *BaseApi) CreateS3Connection(c *gin.Context) {
	var req dto.S3ConnectionCreate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	secret, err := decodeS3Secret(req.SecretKey)
	if err != nil {
		helper.BadRequest(c, err)
		return
	}
	req.SecretKey = secret
	if err := s3Service.CreateConnection(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) UpdateS3Connection(c *gin.Context) {
	var req dto.S3ConnectionUpdate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	secret, err := decodeS3Secret(req.SecretKey)
	if err != nil {
		helper.BadRequest(c, err)
		return
	}
	req.SecretKey = secret
	if err := s3Service.UpdateConnection(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) DeleteS3Connection(c *gin.Context) {
	var req dto.S3ConnID
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := s3Service.DeleteConnection(req.ID); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) TestS3Connection(c *gin.Context) {
	var req dto.S3ConnID
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := s3Service.TestConnection(req.ID); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) ListS3Buckets(c *gin.Context) {
	var req dto.S3ConnID
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	data, err := s3Service.ListBuckets(req.ID)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

func (b *BaseApi) CreateS3Bucket(c *gin.Context) {
	var req dto.S3BucketCreate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := s3Service.CreateBucket(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) DeleteS3Bucket(c *gin.Context) {
	var req dto.S3BucketReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := s3Service.DeleteBucket(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) ListS3Objects(c *gin.Context) {
	var req dto.S3ObjectListReq
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	data, err := s3Service.ListObjects(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

func (b *BaseApi) CreateS3Folder(c *gin.Context) {
	var req dto.S3FolderCreate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := s3Service.CreateFolder(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) DeleteS3Objects(c *gin.Context) {
	var req dto.S3ObjectDelete
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := s3Service.DeleteObjects(req); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}

func (b *BaseApi) PreviewS3Object(c *gin.Context) {
	var req dto.S3ObjectOp
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	data, err := s3Service.Preview(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.SuccessWithData(c, data)
}

func (b *BaseApi) DownloadS3Object(c *gin.Context) {
	var req dto.S3ObjectOp
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	reader, fileName, contentType, size, err := s3Service.Download(req)
	if err != nil {
		helper.InternalServer(c, err)
		return
	}
	defer reader.Close()
	c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(fileName))
	c.DataFromReader(http.StatusOK, size, contentType, reader, nil)
}

func (b *BaseApi) UploadS3Object(c *gin.Context) {
	var req dto.S3UploadForm
	if err := c.ShouldBind(&req); err != nil {
		helper.BadRequest(c, err)
		return
	}
	if req.ID == 0 {
		idStr := c.PostForm("id")
		id, convErr := strconv.ParseUint(idStr, 10, 64)
		if convErr != nil || id == 0 {
			helper.BadRequest(c, errors.New("id is required"))
			return
		}
		req.ID = uint(id)
	}
	if req.Bucket == "" {
		helper.BadRequest(c, errors.New("bucket is required"))
		return
	}
	header, err := c.FormFile("file")
	if err != nil {
		helper.BadRequest(c, err)
		return
	}
	file, err := header.Open()
	if err != nil {
		helper.BadRequest(c, err)
		return
	}
	defer file.Close()
	if err := s3Service.Upload(req, header.Filename, file, header.Size, header.Header.Get("Content-Type")); err != nil {
		helper.InternalServer(c, err)
		return
	}
	helper.Success(c)
}
