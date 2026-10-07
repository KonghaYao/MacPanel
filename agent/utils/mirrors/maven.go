package mirrors

import (
	"html"
	"os"
	"strings"
)

const macpanelMirrorID = "macpanel"

func readMaven(opts Options) (map[string]string, error) {
	path, err := mavenSettingsPath(opts)
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
	repository, err := mavenMirrorURL(string(data))
	if err != nil {
		return nil, err
	}
	return map[string]string{"repository": repository}, nil
}

func writeMaven(opts Options, values map[string]string) error {
	path, err := mavenSettingsPath(opts)
	if err != nil {
		return err
	}
	existing := ""
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		existing = string(data)
	}
	next, remove, err := updateMavenSettings(existing, values["repository"])
	if err != nil {
		return err
	}
	if remove {
		return removeFile(path)
	}
	return writeAtomic(path, []byte(next))
}

func mavenMirrorURL(content string) (string, error) {
	block, ok, err := macpanelMirrorBlock(content)
	if err != nil || !ok {
		return "", err
	}
	text, ok := xmlElementText(block, "url")
	if !ok {
		return "", invalidf("maven macpanel mirror is missing a url")
	}
	return html.UnescapeString(strings.TrimSpace(text)), nil
}

func updateMavenSettings(content, repository string) (string, bool, error) {
	if strings.TrimSpace(content) == "" {
		if repository == "" {
			return "", true, nil
		}
		return mavenDocument(repository), false, nil
	}
	if _, err := elementSpans(content, "settings"); err != nil {
		return "", false, err
	}
	spans, err := elementSpans(content, "settings")
	if err != nil {
		return "", false, err
	}
	if len(spans) == 0 {
		return "", false, invalidf("maven settings.xml is missing a settings element")
	}
	cleaned, err := removeMacpanelMirrors(content)
	if err != nil {
		return "", false, err
	}
	if repository == "" {
		if !mavenHasForeignElements(cleaned) {
			return "", true, nil
		}
		return cleaned, false, nil
	}
	updated, err := insertMacpanelMirror(cleaned, repository)
	if err != nil {
		return "", false, err
	}
	return updated, false, nil
}

func mavenDocument(repository string) string {
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<settings>\n  <mirrors>\n" + mavenMirrorXML(repository, 4) + "\n  </mirrors>\n</settings>\n"
}

func mavenMirrorXML(repository string, indent int) string {
	pad := strings.Repeat(" ", indent)
	return pad + "<mirror>\n" +
		pad + "  <id>" + macpanelMirrorID + "</id>\n" +
		pad + "  <url>" + html.EscapeString(repository) + "</url>\n" +
		pad + "  <mirrorOf>*</mirrorOf>\n" +
		pad + "</mirror>"
}

func removeMacpanelMirrors(content string) (string, error) {
	spans, err := elementSpans(content, "mirror")
	if err != nil {
		return "", err
	}
	var remove []xmlSpan
	for _, span := range spans {
		block := content[span.start:span.end]
		id, ok := xmlElementText(block, "id")
		if ok && strings.TrimSpace(id) == macpanelMirrorID {
			remove = append(remove, span)
		}
	}
	for i := len(remove) - 1; i >= 0; i-- {
		content = content[:remove[i].start] + content[remove[i].end:]
	}
	return content, nil
}

func insertMacpanelMirror(content, repository string) (string, error) {
	mirrors, err := elementSpans(content, "mirrors")
	if err != nil {
		return "", err
	}
	if len(mirrors) > 0 {
		insertAt := mirrors[0].innerStart
		snippet := "\n" + mavenMirrorXML(repository, 4)
		return content[:insertAt] + snippet + content[insertAt:], nil
	}
	settings, err := elementSpans(content, "settings")
	if err != nil {
		return "", err
	}
	if len(settings) == 0 {
		return "", invalidf("maven settings.xml is missing a settings element")
	}
	snippet := "\n  <mirrors>\n" + mavenMirrorXML(repository, 4) + "\n  </mirrors>"
	return content[:settings[0].innerStart] + snippet + content[settings[0].innerStart:], nil
}

func macpanelMirrorBlock(content string) (string, bool, error) {
	spans, err := elementSpans(content, "mirror")
	if err != nil {
		return "", false, err
	}
	for _, span := range spans {
		block := content[span.start:span.end]
		id, ok := xmlElementText(block, "id")
		if ok && strings.TrimSpace(id) == macpanelMirrorID {
			return block, true, nil
		}
	}
	return "", false, nil
}

func mavenHasForeignElements(content string) bool {
	for _, name := range xmlOpenTagNames(content) {
		if name != "settings" && name != "mirrors" {
			return true
		}
	}
	return false
}

func xmlElementText(block, name string) (string, bool) {
	spans, err := elementSpans(block, name)
	if err != nil || len(spans) == 0 {
		return "", false
	}
	return block[spans[0].innerStart:spans[0].innerEnd], true
}
