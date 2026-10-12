package service

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
)

func TestSumBlockIOBytes(t *testing.T) {
	stats := &container.StatsResponse{
		BlkioStats: container.BlkioStats{
			IoServiceBytesRecursive: []container.BlkioStatEntry{
				{Op: "Read", Value: 1024},
				{Op: "read", Value: 512},
				{Op: "write", Value: 2048},
				{Op: "sync", Value: 4096},
			},
		},
	}
	read, write := sumBlockIOBytes(stats)
	if read != 1536 {
		t.Fatalf("expected read 1536, got %d", read)
	}
	if write != 2048 {
		t.Fatalf("expected write 2048, got %d", write)
	}

	empty := &container.StatsResponse{}
	if read, write := sumBlockIOBytes(empty); read != 0 || write != 0 {
		t.Fatalf("expected zero counters, got %d/%d", read, write)
	}
}

func TestSumNetworkBytes(t *testing.T) {
	stats := &container.StatsResponse{
		Networks: map[string]container.NetworkStats{
			"eth0": {RxBytes: 1024, TxBytes: 2048},
			"eth1": {RxBytes: 512, TxBytes: 256},
		},
	}
	rx, tx := sumNetworkBytes(stats)
	if rx != 1536 {
		t.Fatalf("expected rx 1536, got %d", rx)
	}
	if tx != 2304 {
		t.Fatalf("expected tx 2304, got %d", tx)
	}
}

func TestContainerMonitorRate(t *testing.T) {
	if rate := containerMonitorRate(2048, 1024, 10); rate != 102.4 {
		t.Fatalf("expected rate 102.4, got %f", rate)
	}
	if rate := containerMonitorRate(512, 1024, 10); rate != 0 {
		t.Fatalf("counter reset should report 0, got %f", rate)
	}
	if rate := containerMonitorRate(2048, 1024, 0); rate != 0 {
		t.Fatalf("missing elapsed should report 0, got %f", rate)
	}
}

func TestContainerMonitorElapsed(t *testing.T) {
	now := time.Now()
	if elapsed := containerMonitorElapsed(now, containerRawCounter{}); elapsed != 0 {
		t.Fatalf("expected 0 for empty counter, got %f", elapsed)
	}
	if elapsed := containerMonitorElapsed(now, containerRawCounter{at: now.Add(-10 * time.Second)}); elapsed != 10 {
		t.Fatalf("expected 10, got %f", elapsed)
	}
}

func TestContainerMonitorBucketSeconds(t *testing.T) {
	start := time.Unix(1700000000, 0)
	end := start.Add(time.Hour)
	if seconds := containerMonitorBucketSeconds(start, end, containerMonitorMaxRawSample); seconds != 0 {
		t.Fatalf("raw sample range should not be bucketed, got %d", seconds)
	}
	if seconds := containerMonitorBucketSeconds(start, end, containerMonitorMaxRawSample+1); seconds != 7 {
		t.Fatalf("expected 7 seconds per bucket, got %d", seconds)
	}
}
