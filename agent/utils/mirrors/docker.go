package mirrors

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
)

func readDocker(opts Options) (map[string]string, error) {
	path, err := dockerDaemonPath(opts)
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
	if strings.TrimSpace(string(data)) == "" {
		return map[string]string{}, nil
	}
	_, values, err := decodeJSONObject(data)
	if err != nil {
		return nil, err
	}
	raw, ok := values["registry-mirrors"]
	if !ok {
		return map[string]string{}, nil
	}
	var mirrors []string
	if err := json.Unmarshal(raw, &mirrors); err != nil {
		return nil, invalidf("docker registry-mirrors is not an array of strings")
	}
	return map[string]string{"mirrors": strings.Join(mirrors, "\n")}, nil
}

func writeDocker(opts Options, values map[string]string) error {
	path, err := dockerDaemonPath(opts)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if os.IsNotExist(err) {
		data = nil
	}
	if len(bytes.TrimSpace(data)) == 0 {
		data = []byte("{}")
	}
	keys, object, err := decodeJSONObject(data)
	if err != nil {
		return err
	}
	mirrors := splitURLList(values["mirrors"])
	if len(mirrors) == 0 {
		keys, object = deleteJSONKey(keys, object, "registry-mirrors")
	} else {
		raw, err := json.Marshal(mirrors)
		if err != nil {
			return err
		}
		keys, object = setJSONKey(keys, object, "registry-mirrors", raw)
	}
	if len(keys) == 0 {
		if _, statErr := os.Stat(path); statErr != nil && os.IsNotExist(statErr) {
			return nil
		}
	}
	encoded, err := encodeJSONObject(keys, object)
	if err != nil {
		return err
	}
	return writeAtomic(path, encoded)
}

func decodeJSONObject(data []byte) ([]string, map[string]json.RawMessage, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	token, err := dec.Token()
	if err != nil {
		return nil, nil, invalidf("docker daemon.json is not a JSON object")
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '{' {
		return nil, nil, invalidf("docker daemon.json is not a JSON object")
	}
	keys := make([]string, 0)
	values := map[string]json.RawMessage{}
	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return nil, nil, invalidf("docker daemon.json is not a JSON object")
		}
		key, ok := keyToken.(string)
		if !ok {
			return nil, nil, invalidf("docker daemon.json is not a JSON object")
		}
		if _, exists := values[key]; exists {
			return nil, nil, invalidf("docker daemon.json contains duplicate key %s", key)
		}
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return nil, nil, invalidf("docker daemon.json is not a JSON object")
		}
		keys = append(keys, key)
		values[key] = raw
	}
	token, err = dec.Token()
	if err != nil {
		return nil, nil, invalidf("docker daemon.json is not a JSON object")
	}
	delim, ok = token.(json.Delim)
	if !ok || delim != '}' {
		return nil, nil, invalidf("docker daemon.json is not a JSON object")
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, nil, invalidf("docker daemon.json has trailing content")
	}
	return keys, values, nil
}

func setJSONKey(keys []string, values map[string]json.RawMessage, key string, raw json.RawMessage) ([]string, map[string]json.RawMessage) {
	if _, ok := values[key]; !ok {
		keys = append(keys, key)
	}
	values[key] = raw
	return keys, values
}

func deleteJSONKey(keys []string, values map[string]json.RawMessage, key string) ([]string, map[string]json.RawMessage) {
	delete(values, key)
	next := make([]string, 0, len(keys))
	for _, item := range keys {
		if item != key {
			next = append(next, item)
		}
	}
	return next, values
}

func encodeJSONObject(keys []string, values map[string]json.RawMessage) ([]byte, error) {
	if len(keys) == 0 {
		return []byte("{}\n"), nil
	}
	var b strings.Builder
	b.WriteString("{\n")
	for i, key := range keys {
		pretty, err := prettyJSON(values[key])
		if err != nil {
			return nil, err
		}
		encodedKey, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		b.WriteString("  ")
		b.Write(encodedKey)
		b.WriteString(": ")
		b.WriteString(pretty)
		if i != len(keys)-1 {
			b.WriteString(",")
		}
		b.WriteByte('\n')
	}
	b.WriteString("}\n")
	return []byte(b.String()), nil
}

func prettyJSON(raw json.RawMessage) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "  ", "  "); err != nil {
		return "", invalidf("docker daemon.json contains invalid JSON")
	}
	return buf.String(), nil
}
