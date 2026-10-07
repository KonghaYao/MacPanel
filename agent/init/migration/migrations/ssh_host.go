package migrations

import (
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddSSHHostTable = &gormigrate.Migration{
	ID: "20261007-add-ssh-host-table",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&model.SSHHost{})
	},
}
