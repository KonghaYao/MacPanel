package service

import (
	"reflect"
	"testing"
)

func TestMiseUpgradeArgs(t *testing.T) {
	got := miseUpgradeArgs()
	want := [][]string{
		{"mise", "use", "-g", "github:KonghaYao/MacPanel[bin=macpanel]@latest"},
		{"mise", "reshim"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mise upgrade args = %#v, want %#v", got, want)
	}
}
