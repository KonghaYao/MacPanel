package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/common"
	"github.com/1Panel-dev/1Panel/agent/utils/docker"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
)

const (
	containerMonitorConcurrency  = 8
	containerMonitorMaxRawSample = 1200
	containerMonitorBucketStep   = 600
)

var containerMonitorMutex sync.Mutex

type containerMonitorCounters struct {
	mu   sync.Mutex
	data map[string]containerRawCounter
}

type containerRawCounter struct {
	read  uint64
	write uint64
	rx    uint64
	tx    uint64
	at    time.Time
}

var containerMonitorSamples = containerMonitorCounters{data: make(map[string]containerRawCounter)}

func sumBlockIOBytes(stats *container.StatsResponse) (read, write uint64) {
	for _, entry := range stats.BlkioStats.IoServiceBytesRecursive {
		switch strings.ToLower(entry.Op) {
		case "read":
			read += entry.Value
		case "write":
			write += entry.Value
		}
	}
	return read, write
}

func sumNetworkBytes(stats *container.StatsResponse) (rx, tx uint64) {
	for _, network := range stats.Networks {
		rx += network.RxBytes
		tx += network.TxBytes
	}
	return rx, tx
}

func containerMonitorRate(cur, prev uint64, seconds float64) float64 {
	if seconds <= 0 || cur < prev {
		return 0
	}
	return float64(cur-prev) / seconds
}

func containerMonitorBucketSeconds(start, end time.Time, sampleCount int64) int64 {
	if sampleCount <= containerMonitorMaxRawSample {
		return 0
	}
	return (end.Unix() - start.Unix() + containerMonitorBucketStep) / containerMonitorBucketStep
}

func containerMonitorElapsed(now time.Time, prev containerRawCounter) float64 {
	if prev.at.IsZero() {
		return 0
	}
	return now.Sub(prev.at).Seconds()
}

func StartContainerMonitor(interval string) error {
	seconds, err := strconv.Atoi(interval)
	if err != nil || seconds < 10 || seconds > 43200 {
		return fmt.Errorf("invalid container monitoring interval: %s", interval)
	}
	service := &MonitorService{}
	job := cron.NewChain(cron.Recover(cron.DefaultLogger)).Then(cron.FuncJob(service.saveContainerData))
	id, err := global.Cron.AddFunc(fmt.Sprintf("@every %ds", seconds), func() { go job.Run() })
	if err != nil {
		return err
	}
	if global.ContainerMonitorCronID != 0 {
		global.Cron.Remove(global.ContainerMonitorCronID)
	}
	global.ContainerMonitorCronID = id
	return nil
}

func (m *MonitorService) saveContainerData() {
	if !containerMonitorMutex.TryLock() {
		return
	}
	defer containerMonitorMutex.Unlock()
	status, err := settingRepo.GetValueByKey("ContainerMonitorStatus")
	if err != nil {
		global.LOG.Errorf("Load container monitoring status failed: %v", err)
		return
	}
	if status != constant.StatusEnable {
		return
	}
	retention, err := settingRepo.GetValueByKey("ContainerMonitorStoreDays")
	if err != nil {
		global.LOG.Errorf("Load container monitoring retention failed: %v", err)
		return
	}
	days, err := strconv.Atoi(retention)
	if err != nil || days < 1 {
		global.LOG.Errorf("Invalid container monitoring retention: %s", retention)
		return
	}
	if err := containerMonitorRepo.DeleteBefore(time.Now().AddDate(0, 0, -days)); err != nil {
		global.LOG.Errorf("Clean container monitoring data failed: %v", err)
	}

	cli, err := docker.NewDockerClient()
	if err != nil {
		global.LOG.Debugf("Connect docker for container monitoring failed: %v", err)
		return
	}
	defer cli.Close()

	listOptions := container.ListOptions{}
	listOptions.Filters = filters.NewArgs(filters.Arg("status", "running"))
	list, err := cli.ContainerList(context.Background(), listOptions)
	if err != nil {
		global.LOG.Debugf("Load containers for monitoring failed: %v", err)
		return
	}

	seen := make(map[string]struct{}, len(list))
	var counterMutex sync.Mutex
	var group errgroup.Group
	group.SetLimit(containerMonitorConcurrency)
	for _, item := range list {
		if len(item.Names) == 0 {
			continue
		}
		name := strings.TrimPrefix(item.Names[0], "/")
		seen[item.ID] = struct{}{}
		group.Go(func() error {
			point, err := collectContainerPoint(context.Background(), cli, item.ID, name)
			if err != nil {
				global.LOG.Debugf("Collect container metrics for %s failed: %v", name, err)
				return nil
			}
			counterMutex.Lock()
			defer counterMutex.Unlock()
			return containerMonitorRepo.Create(&point)
		})
	}
	if err := group.Wait(); err != nil {
		global.LOG.Errorf("Save container monitoring data failed: %v", err)
	}
	containerMonitorSamples.prune(seen)
}

func (s *containerMonitorCounters) load(containerID string) containerRawCounter {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[containerID]
}

func (s *containerMonitorCounters) store(containerID string, counter containerRawCounter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[containerID] = counter
}

func (s *containerMonitorCounters) prune(seen map[string]struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id := range s.data {
		if _, ok := seen[id]; !ok {
			delete(s.data, id)
		}
	}
}

func collectContainerPoint(ctx context.Context, cli *client.Client, containerID, name string) (model.MonitorContainer, error) {
	point := model.MonitorContainer{ContainerID: containerID, ContainerName: name, CreatedAt: time.Now()}
	res, err := cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		return point, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return point, err
	}
	var stats *container.StatsResponse
	if err := json.Unmarshal(body, &stats); err != nil {
		return point, err
	}

	point.CPUPercent = calculateCPUPercentUnix(stats)
	point.Memory = float64(stats.MemoryStats.Usage) / 1024 / 1024
	point.MemoryLimit = float64(stats.MemoryStats.Limit) / 1024 / 1024
	point.MemoryPercent = calculateMemPercentUnix(stats.MemoryStats)
	if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		point.Cache = float64(cache) / 1024 / 1024
	}

	read, write := sumBlockIOBytes(stats)
	rx, tx := sumNetworkBytes(stats)
	prev := containerMonitorSamples.load(containerID)
	elapsed := containerMonitorElapsed(point.CreatedAt, prev)
	if elapsed > 0 {
		point.IORead = containerMonitorRate(read, prev.read, elapsed) / 1024
		point.IOWrite = containerMonitorRate(write, prev.write, elapsed) / 1024
		point.NetworkRX = containerMonitorRate(rx, prev.rx, elapsed) / 1024
		point.NetworkTX = containerMonitorRate(tx, prev.tx, elapsed) / 1024
	}
	containerMonitorSamples.store(containerID, containerRawCounter{read: read, write: write, rx: rx, tx: tx, at: point.CreatedAt})
	return point, nil
}

func (m *MonitorService) LoadContainerMonitorData(req dto.MonitorContainerSearch) (dto.MonitorContainerData, error) {
	data := dto.MonitorContainerData{Points: []model.MonitorContainer{}}
	if !req.EndTime.After(req.StartTime) {
		return data, fmt.Errorf("invalid container history request")
	}
	loc, err := time.LoadLocation(common.LoadTimeZoneByCmd())
	if err != nil {
		return data, err
	}
	req.StartTime, req.EndTime = req.StartTime.In(loc), req.EndTime.In(loc)
	data.SampleCount, err = containerMonitorRepo.Count(req.ContainerName, req.StartTime, req.EndTime)
	if err != nil || data.SampleCount == 0 {
		return data, err
	}
	data.BucketSeconds = containerMonitorBucketSeconds(req.StartTime, req.EndTime, data.SampleCount)
	points, err := containerMonitorRepo.History(req.ContainerName, req.StartTime, req.EndTime, data.BucketSeconds, req.Aggregation)
	if err != nil {
		return data, err
	}
	if data.BucketSeconds > 0 {
		next := 0
		for bucket := int64(0); bucket <= (req.EndTime.Unix()-req.StartTime.Unix())/data.BucketSeconds; bucket++ {
			var point model.MonitorContainer
			if next < len(points) && points[next].Bucket == bucket {
				point = points[next].MonitorContainer
				next++
			}
			point.ContainerName = req.ContainerName
			point.CreatedAt = time.Unix(req.StartTime.Unix()+bucket*data.BucketSeconds, 0).In(loc)
			if bucket == 0 {
				point.CreatedAt = req.StartTime
			}
			data.Points = append(data.Points, point)
		}
		return data, nil
	}
	setting, err := settingRepo.GetValueByKey("ContainerMonitorInterval")
	if err != nil {
		return data, err
	}
	intervalSeconds, err := strconv.Atoi(setting)
	if err != nil || intervalSeconds <= 0 {
		return data, fmt.Errorf("invalid container monitoring interval: %s", setting)
	}
	interval := time.Duration(intervalSeconds) * time.Second
	for i, point := range points {
		if i > 0 && point.CreatedAt.Sub(points[i-1].CreatedAt) > 2*interval {
			data.Points = append(data.Points, model.MonitorContainer{
				CreatedAt:     points[i-1].CreatedAt.Add(interval),
				ContainerName: req.ContainerName,
			})
		}
		data.Points = append(data.Points, point.MonitorContainer)
	}
	return data, nil
}

func (m *MonitorService) CleanContainerMonitor(req dto.MonitorContainerClean) error {
	containerMonitorMutex.Lock()
	defer containerMonitorMutex.Unlock()
	if len(req.ContainerName) == 0 {
		return containerMonitorRepo.CleanAll()
	}
	return containerMonitorRepo.CleanTarget(req.ContainerName)
}
