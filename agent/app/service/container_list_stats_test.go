package service

import (
	"testing"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
)

func TestCloneContainerListStats(t *testing.T) {
	src := []dto.ContainerListStats{
		{ContainerID: "abc", CPUPercent: 12.5, MemoryUsage: 1024},
	}
	cloned := cloneContainerListStats(src)
	if len(cloned) != 1 {
		t.Fatalf("expected 1 item, got %d", len(cloned))
	}
	if cloned[0].ContainerID != src[0].ContainerID {
		t.Fatalf("unexpected container id: %s", cloned[0].ContainerID)
	}

	cloned[0].CPUPercent = 99
	if src[0].CPUPercent == 99 {
		t.Fatal("clone should not share underlying values")
	}
}

func TestContainerListStatsCacheStore(t *testing.T) {
	cache := &containerListStatsCacheStore{}
	data := []dto.ContainerListStats{{ContainerID: "abc", CPUPercent: 1.2}}

	if _, ok := cache.get(); ok {
		t.Fatal("empty cache should miss")
	}

	cache.set(data)
	cached, ok := cache.get()
	if !ok {
		t.Fatal("expected cache hit immediately after set")
	}
	if len(cached) != 1 || cached[0].ContainerID != "abc" {
		t.Fatalf("unexpected cached data: %+v", cached)
	}

	cache.mu.Lock()
	cache.fetchedAt = time.Now().Add(-containerListStatsCacheTTL - time.Millisecond)
	cache.mu.Unlock()

	if _, ok := cache.get(); ok {
		t.Fatal("expired cache should miss")
	}
}
