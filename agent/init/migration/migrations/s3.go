package migrations

import (
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddS3ConnectionTable = &gormigrate.Migration{
	ID: "20261006-add-s3-connection-table",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&model.S3Connection{})
	},
}
