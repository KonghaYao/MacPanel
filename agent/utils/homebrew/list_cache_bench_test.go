package homebrew

import (
	"fmt"
	"testing"
	"time"
)

func TestListPackagesCacheTiming(t *testing.T) {
	client, err := NewClient(MirrorConfig{})
	if err != nil {
		t.Skip(err)
	}

	InvalidatePackageListCache()
	start := time.Now()
	if _, err := client.ListPackages("all"); err != nil {
		t.Fatal(err)
	}
	cold := time.Since(start)

	start = time.Now()
	if _, err := client.ListPackages("all"); err != nil {
		t.Fatal(err)
	}
	warm := time.Since(start)

	start = time.Now()
	if _, err := client.ListPackages("formula"); err != nil {
		t.Fatal(err)
	}
	warmFormula := time.Since(start)

	t.Logf("cold list(all): %v warm list(all): %v warm list(formula): %v", cold, warm, warmFormula)
	fmt.Printf("cold list(all): %v\nwarm list(all): %v\nwarm list(formula): %v\n", cold, warm, warmFormula)
}
