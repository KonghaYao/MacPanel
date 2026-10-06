package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/app/repo"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/utils/encrypt"
	"gorm.io/gorm"
)

type S3Service struct{}

type IS3Service interface {
	ListRustFS() (dto.S3RustFSList, error)
	ConnectRustFS(req dto.S3ConnectRustFS) (dto.S3ConnectionInfo, error)
	ListConnections() ([]dto.S3ConnectionInfo, error)
	CreateConnection(req dto.S3ConnectionCreate) error
	UpdateConnection(req dto.S3ConnectionUpdate) error
	DeleteConnection(id uint) error
	TestConnection(id uint) error
	ListBuckets(id uint) ([]dto.S3BucketInfo, error)
	CreateBucket(req dto.S3BucketCreate) error
	DeleteBucket(req dto.S3BucketReq) error
	ListObjects(req dto.S3ObjectListReq) (dto.S3ObjectList, error)
	CreateFolder(req dto.S3FolderCreate) error
	DeleteObjects(req dto.S3ObjectDelete) error
	Upload(req dto.S3UploadForm, filename string, reader io.Reader, size int64, contentType string) error
	Download(req dto.S3ObjectOp) (io.ReadCloser, string, string, int64, error)
	Preview(req dto.S3ObjectOp) (dto.S3Preview, error)
}

func NewIS3Service() IS3Service {
	return &S3Service{}
}

func (s *S3Service) ListRustFS() (dto.S3RustFSList, error) {
	var res dto.S3RustFSList
	res.Items = make([]dto.S3RustFSInstall, 0)
	app, err := appRepo.GetFirst(appRepo.WithKey(rustfsAppKey))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, nil
		}
		return res, err
	}
	installs, err := appInstallRepo.ListBy(context.Background(), appInstallRepo.WithAppId(app.ID))
	if err != nil {
		return res, err
	}
	res.Installed = len(installs) > 0
	conns, _ := s3ConnectionRepo.GetList(s3ConnectionRepo.WithSource(s3SourceRustFS))
	connByInstall := map[uint]model.S3Connection{}
	for _, conn := range conns {
		connByInstall[conn.AppInstallID] = conn
	}
	for _, install := range installs {
		item := dto.S3RustFSInstall{
			AppInstallID: install.ID,
			Name:         install.Name,
			Status:       install.Status,
			Version:      install.Version,
		}
		cfg, parseErr := parseRustFSInstall(install)
		if parseErr != nil {
			item.Error = parseErr.Error()
		} else {
			item.Endpoint = cfg.Endpoint
			item.AccessKey = cfg.AccessKey
			item.SecretKey = cfg.SecretKey
			item.APIPort = cfg.APIPort
			item.ConsolePort = cfg.ConsolePort
		}
		if conn, ok := connByInstall[install.ID]; ok {
			item.Connected = true
			item.ConnID = conn.ID
		}
		res.Items = append(res.Items, item)
	}
	return res, nil
}

func (s *S3Service) ConnectRustFS(req dto.S3ConnectRustFS) (dto.S3ConnectionInfo, error) {
	install, err := appInstallRepo.GetFirst(repo.WithByID(req.AppInstallID))
	if err != nil {
		return dto.S3ConnectionInfo{}, err
	}
	if install.App.Key != rustfsAppKey {
		app, appErr := appRepo.GetFirst(repo.WithByID(install.AppId))
		if appErr != nil {
			return dto.S3ConnectionInfo{}, appErr
		}
		if app.Key != rustfsAppKey {
			return dto.S3ConnectionInfo{}, fmt.Errorf("app install %s is not rustfs", install.Name)
		}
		install.App = app
	}
	cfg, err := parseRustFSInstall(install)
	if err != nil {
		return dto.S3ConnectionInfo{}, err
	}
	secretEnc, err := encrypt.StringEncrypt(cfg.SecretKey)
	if err != nil {
		return dto.S3ConnectionInfo{}, err
	}
	conn, err := s3ConnectionRepo.Get(s3ConnectionRepo.WithSource(s3SourceRustFS), repo.WithByAppInstallID(req.AppInstallID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.S3ConnectionInfo{}, err
	}
	if conn.ID == 0 {
		conn = model.S3Connection{
			Name:         cfg.Name,
			Source:       s3SourceRustFS,
			AppInstallID: cfg.AppInstallID,
			Endpoint:     cfg.Endpoint,
			AccessKey:    cfg.AccessKey,
			SecretKey:    secretEnc,
			PathStyle:    true,
			UseSSL:       false,
		}
		if err := s3ConnectionRepo.Create(&conn); err != nil {
			return dto.S3ConnectionInfo{}, err
		}
	} else {
		conn.Endpoint = cfg.Endpoint
		conn.AccessKey = cfg.AccessKey
		conn.SecretKey = secretEnc
		conn.PathStyle = true
		conn.UseSSL = false
		if err := s3ConnectionRepo.Save(&conn); err != nil {
			return dto.S3ConnectionInfo{}, err
		}
	}
	return s.toConnectionInfo(conn)
}

func (s *S3Service) ListConnections() ([]dto.S3ConnectionInfo, error) {
	list, err := s3ConnectionRepo.GetList()
	if err != nil {
		return nil, err
	}
	items := make([]dto.S3ConnectionInfo, 0, len(list))
	for _, conn := range list {
		info, err := s.toConnectionInfo(conn)
		if err != nil {
			return nil, err
		}
		items = append(items, info)
	}
	return items, nil
}

func (s *S3Service) CreateConnection(req dto.S3ConnectionCreate) error {
	if req.Source == s3SourceRustFS {
		if req.AppInstallID == 0 {
			return fmt.Errorf("appInstallId is required")
		}
		_, err := s.ConnectRustFS(dto.S3ConnectRustFS{AppInstallID: req.AppInstallID})
		return err
	}
	if req.Endpoint == "" || req.AccessKey == "" || req.SecretKey == "" {
		return fmt.Errorf("endpoint, accessKey and secretKey are required")
	}
	host, ssl, err := splitS3Endpoint(req.Endpoint, req.UseSSL)
	if err != nil {
		return err
	}
	secretEnc, err := encrypt.StringEncrypt(req.SecretKey)
	if err != nil {
		return err
	}
	conn := model.S3Connection{
		Name:        req.Name,
		Source:      s3SourceManual,
		Endpoint:    host,
		AccessKey:   req.AccessKey,
		SecretKey:   secretEnc,
		Region:      req.Region,
		Bucket:      req.Bucket,
		UseSSL:      ssl,
		PathStyle:   req.PathStyle,
		Description: req.Description,
	}
	return s3ConnectionRepo.Create(&conn)
}

func (s *S3Service) UpdateConnection(req dto.S3ConnectionUpdate) error {
	conn, err := s3ConnectionRepo.Get(repo.WithByID(req.ID))
	if err != nil {
		return err
	}
	if conn.Source == s3SourceRustFS {
		conn.Description = req.Description
		if req.Name != "" {
			conn.Name = req.Name
		}
		return s3ConnectionRepo.Save(&conn)
	}
	if req.Endpoint == "" || req.AccessKey == "" {
		return fmt.Errorf("endpoint and accessKey are required")
	}
	host, ssl, err := splitS3Endpoint(req.Endpoint, req.UseSSL)
	if err != nil {
		return err
	}
	conn.Name = req.Name
	conn.Endpoint = host
	conn.AccessKey = req.AccessKey
	conn.Region = req.Region
	conn.Bucket = req.Bucket
	conn.UseSSL = ssl
	conn.PathStyle = req.PathStyle
	conn.Description = req.Description
	if req.SecretKey != "" {
		secretEnc, err := encrypt.StringEncrypt(req.SecretKey)
		if err != nil {
			return err
		}
		conn.SecretKey = secretEnc
	}
	return s3ConnectionRepo.Save(&conn)
}

func (s *S3Service) DeleteConnection(id uint) error {
	return s3ConnectionRepo.Delete(repo.WithByID(id))
}

func (s *S3Service) TestConnection(id uint) error {
	rt, _, err := s.runtime(id)
	if err != nil {
		return err
	}
	_, err = rt.listBuckets()
	return err
}

func (s *S3Service) ListBuckets(id uint) ([]dto.S3BucketInfo, error) {
	rt, _, err := s.runtime(id)
	if err != nil {
		return nil, err
	}
	buckets, err := rt.listBuckets()
	if err != nil {
		return nil, err
	}
	items := make([]dto.S3BucketInfo, 0, len(buckets))
	for _, b := range buckets {
		items = append(items, dto.S3BucketInfo{Name: b.Name, CreationDate: b.CreationDate})
	}
	return items, nil
}

func (s *S3Service) CreateBucket(req dto.S3BucketCreate) error {
	rt, conn, err := s.runtime(req.ID)
	if err != nil {
		return err
	}
	region := req.Region
	if region == "" {
		region = conn.Region
	}
	return rt.makeBucket(req.Bucket, region)
}

func (s *S3Service) DeleteBucket(req dto.S3BucketReq) error {
	rt, _, err := s.runtime(req.ID)
	if err != nil {
		return err
	}
	return rt.removeBucket(req.Bucket)
}

func (s *S3Service) ListObjects(req dto.S3ObjectListReq) (dto.S3ObjectList, error) {
	rt, _, err := s.runtime(req.ID)
	if err != nil {
		return dto.S3ObjectList{}, err
	}
	prefix := normalizeS3Prefix(req.Prefix)
	objects, err := rt.listObjects(req.Bucket, prefix)
	if err != nil {
		return dto.S3ObjectList{}, err
	}
	items := make([]dto.S3ObjectItem, 0, len(objects))
	seen := make(map[string]struct{}, len(objects))
	for _, object := range objects {
		if prefix != "" && object.Key == prefix {
			continue
		}
		if _, ok := seen[object.Key]; ok {
			continue
		}
		seen[object.Key] = struct{}{}
		isDir := isS3FolderKey(object.Key)
		items = append(items, dto.S3ObjectItem{
			Key:          object.Key,
			Name:         s3ObjectName(object.Key, prefix),
			Prefix:       isDir,
			Size:         object.Size,
			LastModified: object.LastModified,
			ETag:         strings.Trim(object.ETag, `"`),
		})
	}
	return dto.S3ObjectList{Prefix: prefix, Objects: items}, nil
}

func (s *S3Service) CreateFolder(req dto.S3FolderCreate) error {
	if err := validateS3FolderName(req.Name); err != nil {
		return err
	}
	rt, _, err := s.runtime(req.ID)
	if err != nil {
		return err
	}
	key := joinS3Key(req.Prefix, req.Name) + "/"
	return rt.putObject(req.Bucket, key, bytes.NewReader(nil), 0, "application/x-directory")
}

func (s *S3Service) DeleteObjects(req dto.S3ObjectDelete) error {
	rt, _, err := s.runtime(req.ID)
	if err != nil {
		return err
	}
	var keys []string
	for _, key := range req.Keys {
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("object key is required")
		}
		if isS3FolderKey(key) {
			children, listErr := rt.listObjectsRecursive(req.Bucket, key)
			if listErr != nil {
				return listErr
			}
			for _, child := range children {
				keys = append(keys, child.Key)
			}
			keys = append(keys, key)
			continue
		}
		keys = append(keys, key)
	}
	return rt.removeObjects(req.Bucket, uniqueStrings(keys))
}

func (s *S3Service) Upload(req dto.S3UploadForm, filename string, reader io.Reader, size int64, contentType string) error {
	rt, _, err := s.runtime(req.ID)
	if err != nil {
		return err
	}
	key := strings.TrimSpace(req.Key)
	if key == "" {
		key = joinS3Key(req.Prefix, path.Base(filename))
	}
	if key == "" || isS3FolderKey(key) {
		return fmt.Errorf("object key is required")
	}
	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(key))
	}
	return rt.putObject(req.Bucket, key, reader, size, contentType)
}

func (s *S3Service) Download(req dto.S3ObjectOp) (io.ReadCloser, string, string, int64, error) {
	rt, _, err := s.runtime(req.ID)
	if err != nil {
		return nil, "", "", 0, err
	}
	info, err := rt.statObject(req.Bucket, req.Key)
	if err != nil {
		return nil, "", "", 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), s3TransferTimeout)
	obj, err := rt.getObject(ctx, req.Bucket, req.Key)
	if err != nil {
		cancel()
		return nil, "", "", 0, err
	}
	name := path.Base(strings.TrimSuffix(req.Key, "/"))
	if name == "" || name == "." {
		name = "object"
	}
	contentType := info.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return &s3ReadCloser{ReadCloser: obj, cancel: cancel}, name, contentType, info.Size, nil
}

func (s *S3Service) Preview(req dto.S3ObjectOp) (dto.S3Preview, error) {
	rt, _, err := s.runtime(req.ID)
	if err != nil {
		return dto.S3Preview{}, err
	}
	info, err := rt.statObject(req.Bucket, req.Key)
	if err != nil {
		return dto.S3Preview{}, err
	}
	kind := s3PreviewKind(info.ContentType, req.Key)
	res := dto.S3Preview{
		Key:         req.Key,
		Size:        info.Size,
		ContentType: info.ContentType,
		Kind:        kind,
	}
	if kind == s3PreviewNone {
		return res, nil
	}
	limit := s3PreviewLimit(kind)
	if info.Size > limit {
		res.TooLarge = true
		return res, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), s3DefaultTimeout)
	defer cancel()
	obj, err := rt.getObject(ctx, req.Bucket, req.Key)
	if err != nil {
		return dto.S3Preview{}, err
	}
	defer obj.Close()
	body, err := readLimited(obj, limit)
	if err != nil {
		return dto.S3Preview{}, err
	}
	if kind == s3PreviewText {
		if !utf8.Valid(body) {
			res.Kind = s3PreviewNone
			return res, nil
		}
		res.Content = string(body)
		res.Encoding = "utf8"
		return res, nil
	}
	res.Content = base64.StdEncoding.EncodeToString(body)
	res.Encoding = "base64"
	return res, nil
}

func (s *S3Service) runtime(id uint) (*s3Runtime, model.S3Connection, error) {
	conn, err := s3ConnectionRepo.Get(repo.WithByID(id))
	if err != nil {
		return nil, model.S3Connection{}, err
	}
	endpoint := conn.Endpoint
	accessKey := conn.AccessKey
	secretKey, err := encrypt.StringDecrypt(conn.SecretKey)
	if err != nil {
		return nil, model.S3Connection{}, err
	}
	useSSL := conn.UseSSL
	pathStyle := conn.PathStyle
	region := conn.Region
	if conn.Source == s3SourceRustFS {
		cfg, liveErr := s.liveRustFS(conn.AppInstallID)
		if liveErr != nil {
			return nil, model.S3Connection{}, liveErr
		}
		endpoint = cfg.Endpoint
		accessKey = cfg.AccessKey
		secretKey = cfg.SecretKey
		useSSL = cfg.UseSSL
		pathStyle = cfg.PathStyle
	}
	rt, err := newS3Runtime(endpoint, accessKey, secretKey, region, useSSL, pathStyle)
	if err != nil {
		return nil, model.S3Connection{}, err
	}
	return rt, conn, nil
}

func (s *S3Service) liveRustFS(appInstallID uint) (rustFSConnConfig, error) {
	install, err := appInstallRepo.GetFirst(repo.WithByID(appInstallID))
	if err != nil {
		return rustFSConnConfig{}, err
	}
	if install.Status != constant.StatusRunning {
		return rustFSConnConfig{}, fmt.Errorf("rustfs install %s is not running", install.Name)
	}
	return parseRustFSInstall(install)
}

func (s *S3Service) toConnectionInfo(conn model.S3Connection) (dto.S3ConnectionInfo, error) {
	secret, err := encrypt.StringDecrypt(conn.SecretKey)
	if err != nil {
		return dto.S3ConnectionInfo{}, err
	}
	info := dto.S3ConnectionInfo{
		ID:           conn.ID,
		CreatedAt:    conn.CreatedAt,
		Name:         conn.Name,
		Source:       conn.Source,
		AppInstallID: conn.AppInstallID,
		Endpoint:     conn.Endpoint,
		AccessKey:    conn.AccessKey,
		SecretKey:    secret,
		Region:       conn.Region,
		Bucket:       conn.Bucket,
		UseSSL:       conn.UseSSL,
		PathStyle:    conn.PathStyle,
		Description:  conn.Description,
	}
	if conn.Source == s3SourceRustFS && conn.AppInstallID > 0 {
		install, err := appInstallRepo.GetFirst(repo.WithByID(conn.AppInstallID))
		if err == nil {
			info.AppName = install.Name
			info.Status = install.Status
			if cfg, parseErr := parseRustFSInstall(install); parseErr == nil {
				info.Endpoint = cfg.Endpoint
				info.AccessKey = cfg.AccessKey
				info.SecretKey = cfg.SecretKey
				info.APIPort = cfg.APIPort
				info.ConsolePort = cfg.ConsolePort
				info.UseSSL = cfg.UseSSL
				info.PathStyle = cfg.PathStyle
			}
		}
	}
	return info, nil
}

type s3ReadCloser struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (r *s3ReadCloser) Close() error {
	if r.cancel != nil {
		r.cancel()
	}
	if r.ReadCloser == nil {
		return nil
	}
	return r.ReadCloser.Close()
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}
