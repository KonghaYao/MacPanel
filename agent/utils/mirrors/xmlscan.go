package mirrors

import "strings"

type xmlSpan struct {
	start      int
	end        int
	innerStart int
	innerEnd   int
}

func elementSpans(content, name string) ([]xmlSpan, error) {
	var spans []xmlSpan
	i := 0
	for i < len(content) {
		next, err := skipXMLNoise(content, i)
		if err != nil {
			return nil, err
		}
		if next != i {
			i = next
			continue
		}
		if content[i] != '<' {
			i++
			continue
		}
		if strings.HasPrefix(content[i:], "</") {
			i++
			continue
		}
		if !isOpenTag(content, i, name) {
			i++
			continue
		}
		gt := strings.IndexByte(content[i:], '>')
		if gt < 0 {
			return nil, invalidf("unclosed xml tag %s", name)
		}
		startTagEnd := i + gt + 1
		selfClose := gt > 0 && content[i+gt-1] == '/'
		if selfClose {
			spans = append(spans, xmlSpan{start: i, end: startTagEnd, innerStart: startTagEnd, innerEnd: startTagEnd})
			i = startTagEnd
			continue
		}
		closeAt, innerEnd, err := findCloseTag(content, startTagEnd, name)
		if err != nil {
			return nil, err
		}
		spans = append(spans, xmlSpan{start: i, end: closeAt, innerStart: startTagEnd, innerEnd: innerEnd})
		i = closeAt
	}
	return spans, nil
}

func findCloseTag(content string, from int, name string) (int, int, error) {
	depth := 1
	i := from
	for i < len(content) {
		next, err := skipXMLNoise(content, i)
		if err != nil {
			return 0, 0, err
		}
		if next != i {
			i = next
			continue
		}
		if isCloseTag(content, i, name) {
			depth--
			if depth == 0 {
				gt := strings.IndexByte(content[i:], '>')
				if gt < 0 {
					return 0, 0, invalidf("unclosed xml tag %s", name)
				}
				return i + gt + 1, i, nil
			}
		} else if isOpenTag(content, i, name) {
			gt := strings.IndexByte(content[i:], '>')
			if gt < 0 {
				return 0, 0, invalidf("unclosed xml tag %s", name)
			}
			if !(gt > 0 && content[i+gt-1] == '/') {
				depth++
			}
			i = i + gt + 1
			continue
		}
		i++
	}
	return 0, 0, invalidf("unclosed xml tag %s", name)
}

func skipXMLNoise(content string, i int) (int, error) {
	if strings.HasPrefix(content[i:], "<!--") {
		end := strings.Index(content[i+4:], "-->")
		if end < 0 {
			return 0, invalidf("unclosed xml comment")
		}
		return i + 4 + end + 3, nil
	}
	if strings.HasPrefix(content[i:], "<![CDATA[") {
		end := strings.Index(content[i:], "]]>")
		if end < 0 {
			return 0, invalidf("unclosed xml cdata")
		}
		return i + end + 3, nil
	}
	if strings.HasPrefix(content[i:], "<?") {
		end := strings.Index(content[i:], "?>")
		if end < 0 {
			return 0, invalidf("unclosed xml processing instruction")
		}
		return i + end + 2, nil
	}
	return i, nil
}

func isOpenTag(content string, i int, name string) bool {
	if i >= len(content) || content[i] != '<' {
		return false
	}
	rest := content[i+1:]
	if !strings.HasPrefix(rest, name) {
		return false
	}
	if len(rest) == len(name) {
		return false
	}
	return isTagBoundary(rest[len(name)])
}

func isCloseTag(content string, i int, name string) bool {
	prefix := "</" + name
	if !strings.HasPrefix(content[i:], prefix) {
		return false
	}
	if len(content) == i+len(prefix) {
		return false
	}
	return isTagBoundary(content[i+len(prefix)])
}

func isTagBoundary(char byte) bool {
	return char == '>' || char == '/' || char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

func xmlOpenTagNames(content string) []string {
	var names []string
	i := 0
	for i < len(content) {
		next, err := skipXMLNoise(content, i)
		if err != nil || next != i {
			if err != nil {
				return names
			}
			i = next
			continue
		}
		if content[i] != '<' || strings.HasPrefix(content[i:], "</") {
			i++
			continue
		}
		name, ok := tagNameAt(content, i+1)
		if ok {
			names = append(names, name)
		}
		i++
	}
	return names
}

func tagNameAt(content string, i int) (string, bool) {
	if i >= len(content) {
		return "", false
	}
	j := i
	for j < len(content) {
		char := content[j]
		if isTagBoundary(char) {
			break
		}
		if !isNameChar(char) {
			return "", false
		}
		j++
	}
	if j == i {
		return "", false
	}
	return content[i:j], true
}

func isNameChar(char byte) bool {
	return char == '_' || char == ':' || char == '-' || char == '.' ||
		(char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}
