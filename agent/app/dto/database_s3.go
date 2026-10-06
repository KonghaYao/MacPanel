package dto

import "time"

type S3ConnectionCreate struct {
	Name         string `json:"name" validate:"required"`
	Source       string `json:"source" validate:"required,oneof=rustfs manual"`
	AppInstallID uint   `json:"appInstallId"`
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Region       string `json:"region"`
	Bucket       string `json:"bucket"`
	UseSSL       bool   `json:"useSSL"`
	PathStyle    bool   `json:"pathStyle"`
	Description  string `json:"description"`
}

type S3ConnectionUpdate struct {
	ID uint `json:"id" validate:"required"`
	S3ConnectionCreate
}

type S3ConnectionInfo struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Name         string    `json:"name"`
	Source       string    `json:"source"`
	AppInstallID uint      `json:"appInstallId"`
	AppName      string    `json:"appName"`
	Status       string    `json:"status"`
	Endpoint     string    `json:"endpoint"`
	AccessKey    string    `json:"accessKey"`
	SecretKey    string    `json:"secretKey"`
	Region       string    `json:"region"`
	Bucket       string    `json:"bucket"`
	UseSSL       bool      `json:"useSSL"`
	PathStyle    bool      `json:"pathStyle"`
	Description  string    `json:"description"`
	APIPort      int       `json:"apiPort"`
	ConsolePort  int       `json:"consolePort"`
}

type S3RustFSInstall struct {
	AppInstallID uint   `json:"appInstallId"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Version      string `json:"version"`
	Endpoint     string `json:"endpoint"`
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	APIPort      int    `json:"apiPort"`
	ConsolePort  int    `json:"consolePort"`
	Connected    bool   `json:"connected"`
	ConnID       uint   `json:"connId"`
	Error        string `json:"error"`
}

type S3RustFSList struct {
	Installed bool              `json:"installed"`
	Items     []S3RustFSInstall `json:"items"`
}

type S3ConnectRustFS struct {
	AppInstallID uint `json:"appInstallId" validate:"required"`
}

type S3ConnID struct {
	ID uint `json:"id" validate:"required"`
}

type S3BucketCreate struct {
	ID     uint   `json:"id" validate:"required"`
	Bucket string `json:"bucket" validate:"required"`
	Region string `json:"region"`
}

type S3BucketReq struct {
	ID     uint   `json:"id" validate:"required"`
	Bucket string `json:"bucket" validate:"required"`
}

type S3ObjectListReq struct {
	ID     uint   `json:"id" validate:"required"`
	Bucket string `json:"bucket" validate:"required"`
	Prefix string `json:"prefix"`
}

type S3FolderCreate struct {
	ID     uint   `json:"id" validate:"required"`
	Bucket string `json:"bucket" validate:"required"`
	Prefix string `json:"prefix"`
	Name   string `json:"name" validate:"required"`
}

type S3ObjectDelete struct {
	ID     uint     `json:"id" validate:"required"`
	Bucket string   `json:"bucket" validate:"required"`
	Keys   []string `json:"keys" validate:"required,min=1"`
}

type S3ObjectOp struct {
	ID     uint   `json:"id" validate:"required"`
	Bucket string `json:"bucket" validate:"required"`
	Key    string `json:"key" validate:"required"`
}

type S3BucketInfo struct {
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
}

type S3ObjectItem struct {
	Key          string    `json:"key"`
	Name         string    `json:"name"`
	Prefix       bool      `json:"prefix"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	ETag         string    `json:"etag"`
}

type S3ObjectList struct {
	Prefix  string         `json:"prefix"`
	Objects []S3ObjectItem `json:"objects"`
}

type S3Preview struct {
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
	Kind        string `json:"kind"`
	TooLarge    bool   `json:"tooLarge"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
}

type S3UploadForm struct {
	ID     uint   `form:"id" validate:"required"`
	Bucket string `form:"bucket" validate:"required"`
	Prefix string `form:"prefix"`
	Key    string `form:"key"`
}
