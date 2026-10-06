package repo

import (
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/global"
	"gorm.io/gorm"
)

type S3ConnectionRepo struct{}

type IS3ConnectionRepo interface {
	Get(opts ...DBOption) (model.S3Connection, error)
	GetList(opts ...DBOption) ([]model.S3Connection, error)
	Create(conn *model.S3Connection) error
	Save(conn *model.S3Connection) error
	Delete(opts ...DBOption) error
	WithSource(source string) DBOption
}

func NewIS3ConnectionRepo() IS3ConnectionRepo {
	return &S3ConnectionRepo{}
}

func (s *S3ConnectionRepo) WithSource(source string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("source = ?", source)
	}
}

func (s *S3ConnectionRepo) Get(opts ...DBOption) (model.S3Connection, error) {
	var conn model.S3Connection
	err := getDb(opts...).Model(&model.S3Connection{}).First(&conn).Error
	return conn, err
}

func (s *S3ConnectionRepo) GetList(opts ...DBOption) ([]model.S3Connection, error) {
	var list []model.S3Connection
	err := getDb(opts...).Model(&model.S3Connection{}).Order("id desc").Find(&list).Error
	return list, err
}

func (s *S3ConnectionRepo) Create(conn *model.S3Connection) error {
	return global.DB.Create(conn).Error
}

func (s *S3ConnectionRepo) Save(conn *model.S3Connection) error {
	return global.DB.Save(conn).Error
}

func (s *S3ConnectionRepo) Delete(opts ...DBOption) error {
	return getDb(opts...).Delete(&model.S3Connection{}).Error
}
