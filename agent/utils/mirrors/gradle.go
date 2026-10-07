package mirrors

import (
	"os"
	"regexp"
	"strings"
)

var gradleMirrorPattern = regexp.MustCompile(`(?m)^def macpanelMirror = '([^']*)'$`)

func readGradle(opts Options) (map[string]string, error) {
	path, err := gradleInitPath(opts)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	match := gradleMirrorPattern.FindSubmatch(data)
	if match == nil {
		return nil, invalidf("gradle init script %s has no macpanelMirror assignment", path)
	}
	return map[string]string{"repository": string(match[1])}, nil
}

func writeGradle(opts Options, values map[string]string) error {
	path, err := gradleInitPath(opts)
	if err != nil {
		return err
	}
	if values["repository"] == "" {
		return removeFile(path)
	}
	script := strings.Replace(gradleScriptTemplate, "__MIRROR__", values["repository"], 1)
	return writeAtomic(path, []byte(script))
}

const gradleScriptTemplate = `// Managed by MacPanel.
def macpanelMirror = '__MIRROR__'

allprojects { project ->
    buildscript {
        repositories {
            all { ArtifactRepository repo ->
                if (repo instanceof MavenArtifactRepository) {
                    def repoUrl = repo.url.toString()
                    if (repoUrl.startsWith('https://repo.maven.apache.org/maven2') || repoUrl.startsWith('https://repo1.maven.org/maven2')) {
                        remove repo
                    }
                }
            }
            maven { url macpanelMirror }
        }
    }
    repositories {
        all { ArtifactRepository repo ->
            if (repo instanceof MavenArtifactRepository) {
                def repoUrl = repo.url.toString()
                if (repoUrl.startsWith('https://repo.maven.apache.org/maven2') || repoUrl.startsWith('https://repo1.maven.org/maven2')) {
                    remove repo
                }
            }
        }
        maven { url macpanelMirror }
    }
}
`
