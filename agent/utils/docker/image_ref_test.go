package docker

import "testing"

func TestParseImageReference(t *testing.T) {
	tests := []struct {
		ref      string
		wantName string
		wantTag  string
	}{
		{"nginx:latest", "nginx", "latest"},
		{"registry.io/foo/bar:v1", "registry.io/foo/bar", "v1"},
		{"localhost:5000/myapp:1.0", "localhost:5000/myapp", "1.0"},
		{"<none>:<none>", "<none>:<none>", ""},
		{"library/nginx", "nginx", ""},
	}
	for _, tt := range tests {
		name, tag := ParseImageReference(tt.ref)
		if name != tt.wantName || tag != tt.wantTag {
			t.Fatalf("ParseImageReference(%q) = (%q, %q), want (%q, %q)", tt.ref, name, tag, tt.wantName, tt.wantTag)
		}
	}
}

func TestImageRepositoryMatchesFilter(t *testing.T) {
	tags := []string{"registry.io/foo/bar:v1", "nginx:latest"}
	if !ImageRepositoryMatchesFilter(tags, "nginx") {
		t.Fatal("expected nginx to match")
	}
	if !ImageRepositoryMatchesFilter(tags, "foo/bar") {
		t.Fatal("expected foo/bar to match")
	}
	if ImageRepositoryMatchesFilter(tags, "redis") {
		t.Fatal("expected redis not to match")
	}
}
