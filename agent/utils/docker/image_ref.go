package docker

import (
	"strings"

	"github.com/distribution/reference"
)

// ParseImageReference splits a Docker image reference into repository name and tag.
func ParseImageReference(ref string) (name, tag string) {
	ref = strings.TrimSpace(ref)
	if ref == "" || strings.Contains(ref, "<none>") {
		return ref, ""
	}
	named, err := reference.ParseNormalizedNamed(ref)
	if err != nil {
		return splitImageReferenceFallback(ref)
	}
	name = reference.FamiliarName(named)
	if tagged, ok := named.(reference.Tagged); ok {
		tag = tagged.Tag()
	}
	return name, tag
}

func splitImageReferenceFallback(ref string) (name, tag string) {
	tagSep := findTagSeparator(ref)
	if tagSep < 0 {
		return ref, ""
	}
	return ref[:tagSep], ref[tagSep+1:]
}

func findTagSeparator(ref string) int {
	colon := strings.LastIndex(ref, ":")
	if colon <= 0 {
		return -1
	}
	slash := strings.LastIndex(ref, "/")
	if colon <= slash {
		return -1
	}
	firstSlash := strings.Index(ref, "/")
	if firstSlash > 0 && colon < firstSlash {
		nextColon := strings.Index(ref[firstSlash:], ":")
		if nextColon >= 0 {
			return firstSlash + nextColon
		}
		return -1
	}
	return colon
}

// ImageRepositoryMatchesFilter reports whether any tag matches the filter by repository or short name.
func ImageRepositoryMatchesFilter(tags []string, filter string) bool {
	filter = strings.ToLower(strings.TrimSpace(filter))
	if filter == "" {
		return true
	}
	for _, item := range tags {
		name, tagPart := ParseImageReference(item)
		nameLower := strings.ToLower(name)
		shortName := name
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			shortName = name[idx+1:]
		}
		if strings.Contains(nameLower, filter) || strings.Contains(strings.ToLower(shortName), filter) {
			return true
		}
		if tagPart != "" && strings.Contains(strings.ToLower(tagPart), filter) {
			return true
		}
		if strings.Contains(strings.ToLower(item), filter) {
			return true
		}
	}
	return false
}

// PrimaryImageTag returns the first tag used for sorting.
func PrimaryImageTag(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return tags[0]
}
