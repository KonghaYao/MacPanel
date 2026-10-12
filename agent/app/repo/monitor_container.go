package repo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/global"
	"gorm.io/gorm"
)

type ContainerMonitorRepo struct{}

type ContainerHistoryPoint struct {
	model.MonitorContainer
	Bucket int64
}

var containerMonitorMetrics = []string{"cpu_percent", "memory", "memory_limit", "memory_percent", "cache", "io_read", "io_write", "network_rx", "network_tx"}

func (r *ContainerMonitorRepo) Create(point *model.MonitorContainer) error {
	return global.ContainerMonitorDB.Create(point).Error
}

func (r *ContainerMonitorRepo) Latest(name string) (model.MonitorContainer, error) {
	var point model.MonitorContainer
	db := global.ContainerMonitorDB.Where("container_name = ?", name)
	err := db.Order("created_at DESC, id DESC").First(&point).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return point, nil
	}
	return point, err
}

func (r *ContainerMonitorRepo) CleanTarget(name string) error {
	return global.ContainerMonitorDB.Where("container_name = ?", name).Delete(&model.MonitorContainer{}).Error
}

func (r *ContainerMonitorRepo) CleanAll() error {
	return global.ContainerMonitorDB.Where("1 = 1").Delete(&model.MonitorContainer{}).Error
}

func (r *ContainerMonitorRepo) DeleteBefore(before time.Time) error {
	return global.ContainerMonitorDB.Where("created_at < ?", before).Delete(&model.MonitorContainer{}).Error
}

func (r *ContainerMonitorRepo) Count(name string, start, end time.Time) (int64, error) {
	var count int64
	db := global.ContainerMonitorDB.Model(&model.MonitorContainer{}).Where("container_name = ? AND created_at >= ? AND created_at <= ?", name, start, end)
	err := db.Count(&count).Error
	return count, err
}

func (r *ContainerMonitorRepo) History(name string, start, end time.Time, seconds int64, aggregation string) ([]ContainerHistoryPoint, error) {
	db := global.ContainerMonitorDB.Model(&model.MonitorContainer{}).Where("container_name = ? AND created_at >= ? AND created_at <= ?", name, start, end)
	var columns []string
	if seconds > 0 {
		operation := "AVG"
		if aggregation == "max" {
			operation = "MAX"
		}
		columns = []string{fmt.Sprintf("(CAST(strftime('%%s', created_at) AS INTEGER) - %d) / %d AS bucket", start.Unix(), seconds)}
		for _, column := range containerMonitorMetrics {
			columns = append(columns, operation+"("+column+") AS "+column)
		}
		db = db.Group("bucket").Order("bucket ASC")
	} else {
		columns = append([]string{"id", "created_at", "container_id", "container_name"}, containerMonitorMetrics...)
		db = db.Order("created_at ASC, id ASC")
	}
	var points []ContainerHistoryPoint
	err := db.Select(strings.Join(columns, ", ")).Scan(&points).Error
	return points, err
}
