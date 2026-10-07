package mirrors

import (
	"os"
	"strings"
)

const (
	cargoSourceName = "macpanel"
	cargoSourceKey  = "source." + cargoSourceName
	cargoRegistry   = "registries." + cargoSourceName
)

func readCargo(opts Options) (map[string]string, error) {
	path, err := cargoConfigPath(opts)
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
	file, err := parseTOMLTables(string(data))
	if err != nil {
		return nil, err
	}
	replaceWith := unquoteTOML(file.key("source.crates-io", "replace-with"))
	if replaceWith == "" {
		return map[string]string{}, nil
	}
	registry := unquoteTOML(file.key("source."+replaceWith, "registry"))
	return map[string]string{"registry": registry}, nil
}

func writeCargo(opts Options, values map[string]string) error {
	path, err := cargoConfigPath(opts)
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
	file, err := parseTOMLTables(existing)
	if err != nil {
		return err
	}
	registry := values["registry"]
	if registry == "" {
		file.deleteKey("source.crates-io", "replace-with")
		file.deleteTableIfEmpty("source.crates-io")
		file.deleteTable(cargoSourceKey)
		file.deleteTable(cargoRegistry)
	} else {
		crates := file.ensure("source.crates-io")
		crates.set("replace-with", quoteTOML(cargoSourceName))
		source := file.ensure(cargoSourceKey)
		source.set("registry", quoteTOML(registry))
		reg := file.ensure(cargoRegistry)
		reg.set("index", quoteTOML(registry))
	}
	rendered := file.render()
	if strings.TrimSpace(rendered) == "" {
		return removeFile(path)
	}
	return writeAtomic(path, []byte(rendered))
}

type tomlTable struct {
	name  string
	lines []string
}

type tomlFile struct {
	preamble string
	tables   []tomlTable
}

func parseTOMLTables(content string) (*tomlFile, error) {
	file := &tomlFile{}
	if content == "" {
		return file, nil
	}
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	current := -1
	var preamble []string
	for _, line := range lines {
		if strings.Contains(line, `"""`) || strings.Contains(line, "'''") {
			return nil, invalidf("cargo config.toml uses a multiline string, which this editor does not rewrite")
		}
		if strings.HasPrefix(strings.TrimSpace(line), "[[") {
			return nil, invalidf("cargo config.toml uses an array of tables, which this editor does not rewrite")
		}
		if name, ok := tomlTableName(line); ok {
			file.tables = append(file.tables, tomlTable{name: name})
			current = len(file.tables) - 1
			continue
		}
		if current < 0 {
			preamble = append(preamble, line)
			continue
		}
		file.tables[current].lines = append(file.tables[current].lines, line)
	}
	if len(preamble) > 0 {
		file.preamble = strings.Join(preamble, "\n")
		if !strings.HasSuffix(file.preamble, "\n") {
			file.preamble += "\n"
		}
	}
	return file, nil
}

func tomlTableName(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "[[") {
		return "", false
	}
	end := strings.Index(trimmed, "]")
	if end <= 1 {
		return "", false
	}
	name := trimmed[1:end]
	rest := strings.TrimSpace(trimmed[end+1:])
	if rest != "" && !strings.HasPrefix(rest, "#") {
		return "", false
	}
	if strings.ContainsAny(name, "[]\"' ") {
		return "", false
	}
	return name, true
}

func (f *tomlFile) key(table, key string) string {
	for _, item := range f.tables {
		if item.name != table {
			continue
		}
		for _, line := range item.lines {
			name, value, ok := splitAssign(line)
			if ok && name == key {
				return value
			}
		}
	}
	return ""
}

func (f *tomlFile) ensure(name string) *tomlTable {
	for i := range f.tables {
		if f.tables[i].name == name {
			return &f.tables[i]
		}
	}
	f.tables = append(f.tables, tomlTable{name: name})
	return &f.tables[len(f.tables)-1]
}

func (t *tomlTable) set(key, rawValue string) {
	assignment := key + " = " + rawValue
	for i, line := range t.lines {
		name, _, ok := splitAssign(line)
		if ok && name == key {
			t.lines[i] = assignment
			return
		}
	}
	t.lines = append(t.lines, assignment)
}

func (f *tomlFile) deleteKey(table, key string) {
	for i := range f.tables {
		if f.tables[i].name != table {
			continue
		}
		var lines []string
		for _, line := range f.tables[i].lines {
			name, _, ok := splitAssign(line)
			if ok && name == key {
				continue
			}
			lines = append(lines, line)
		}
		f.tables[i].lines = lines
	}
}

func (f *tomlFile) deleteTable(name string) {
	var tables []tomlTable
	for _, table := range f.tables {
		if table.name == name {
			continue
		}
		tables = append(tables, table)
	}
	f.tables = tables
}

func (f *tomlFile) deleteTableIfEmpty(name string) {
	for _, table := range f.tables {
		if table.name != name {
			continue
		}
		for _, line := range table.lines {
			if _, _, ok := splitAssign(line); ok {
				return
			}
		}
		f.deleteTable(name)
		return
	}
}

func (f *tomlFile) render() string {
	var b strings.Builder
	b.WriteString(f.preamble)
	for _, table := range f.tables {
		if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n") {
			b.WriteByte('\n')
		}
		if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n\n") {
			b.WriteByte('\n')
		}
		b.WriteString("[")
		b.WriteString(table.name)
		b.WriteString("]\n")
		for _, line := range table.lines {
			b.WriteString(line)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func quoteTOML(value string) string {
	return `"` + value + `"`
}

func unquoteTOML(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		if value[0] == '"' && strings.HasSuffix(value, `"`) && !strings.Contains(value[1:len(value)-1], `"`) {
			return value[1 : len(value)-1]
		}
		if value[0] == '\'' && strings.HasSuffix(value, "'") && !strings.Contains(value[1:len(value)-1], "'") {
			return value[1 : len(value)-1]
		}
	}
	if idx := strings.Index(value, " #"); idx >= 0 {
		value = strings.TrimSpace(value[:idx])
	}
	return value
}
