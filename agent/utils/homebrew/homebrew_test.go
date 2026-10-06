package homebrew

import (
	"testing"
	"time"
)

func TestPackageListCache(t *testing.T) {
	InvalidatePackageListCache()

	client := &Client{mirrorEnv: []string{"HOMEBREW_BOTTLE_DOMAIN=https://example.com"}}

	formulaItems := []PackageItem{{Name: "wget", Version: "1.0", Type: "formula"}}
	caskItems := []PackageItem{{Name: "firefox", Version: "1.0", Type: "cask"}}
	client.storeCachedPackageLists(formulaItems, caskItems)

	gotFormula, gotCask, ok := client.loadCachedPackageLists()
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(gotFormula) != 1 || gotFormula[0].Name != "wget" {
		t.Fatalf("unexpected formula cache: %#v", gotFormula)
	}
	if len(gotCask) != 1 || gotCask[0].Name != "firefox" {
		t.Fatalf("unexpected cask cache: %#v", gotCask)
	}

	packageListCacheMu.Lock()
	packageListCacheAt = time.Now().Add(-packageListCacheTTL - time.Second)
	packageListCacheMu.Unlock()

	if _, _, ok := client.loadCachedPackageLists(); ok {
		t.Fatal("expected cache miss after ttl")
	}

	InvalidatePackageListCache()
	if _, _, ok := client.loadCachedPackageLists(); ok {
		t.Fatal("expected cache miss after invalidation")
	}
}
