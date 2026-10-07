package repo

import (
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/global"
	"gorm.io/gorm"
)

type SSHHostRepo struct{}

type ISSHHostRepo interface {
	Get(opts ...DBOption) (model.SSHHost, error)
	List(opts ...DBOption) ([]model.SSHHost, error)
	Page(page, size int, opts ...DBOption) (int64, []model.SSHHost, error)
	Create(host *model.SSHHost) error
	Delete(opts ...DBOption) error

	WithByAlias(alias string) DBOption
	WithByInfo(info string) DBOption
}

func NewISSHHostRepo() ISSHHostRepo {
	return &SSHHostRepo{}
}

func (r *SSHHostRepo) Get(opts ...DBOption) (model.SSHHost, error) {
	var host model.SSHHost
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&host).Error
	return host, err
}

func (r *SSHHostRepo) List(opts ...DBOption) ([]model.SSHHost, error) {
	var hosts []model.SSHHost
	db := global.DB.Model(&model.SSHHost{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Order("id desc").Find(&hosts).Error
	return hosts, err
}

func (r *SSHHostRepo) Page(page, size int, opts ...DBOption) (int64, []model.SSHHost, error) {
	var hosts []model.SSHHost
	db := global.DB.Model(&model.SSHHost{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Order("id desc").Limit(size).Offset(size * (page - 1)).Find(&hosts).Error
	return count, hosts, err
}

func (r *SSHHostRepo) Create(host *model.SSHHost) error {
	return global.DB.Create(host).Error
}

func (r *SSHHostRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&model.SSHHost{}).Error
}

func (r *SSHHostRepo) WithByAlias(alias string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("alias = ?", alias)
	}
}

func (r *SSHHostRepo) WithByInfo(info string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(info) == 0 {
			return g
		}
		infoStr := "%" + info + "%"
		return g.Where("alias LIKE ? OR host_name LIKE ? OR user LIKE ?", infoStr, infoStr, infoStr)
	}
}
