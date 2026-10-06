package model

type S3Connection struct {
	BaseModel
	Name         string `json:"name" gorm:"not null;uniqueIndex"`
	Source       string `json:"source" gorm:"not null"`
	AppInstallID uint   `json:"appInstallId"`
	Endpoint     string `json:"endpoint" gorm:"not null"`
	AccessKey    string `json:"accessKey" gorm:"not null"`
	SecretKey    string `json:"secretKey" gorm:"not null"`
	Region       string `json:"region"`
	Bucket       string `json:"bucket"`
	UseSSL       bool   `json:"useSSL"`
	PathStyle    bool   `json:"pathStyle"`
	Description  string `json:"description"`
}
