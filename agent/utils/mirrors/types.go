package mirrors

import (
	"fmt"
	"strings"
)

const (
	KindURL     = "url"
	KindURLs    = "urls"
	KindHost    = "host"
	KindGoproxy = "goproxy"
	KindGosumdb = "gosumdb"
	KindCargo   = "cargo"

	EcosystemHomebrew = "homebrew"
)

type Field struct {
	Key  string `json:"key"`
	Kind string `json:"kind"`
}

type Preset struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Values  map[string]string `json:"values"`
	Snippet string            `json:"snippet"`
}

type Ecosystem struct {
	ID           string            `json:"id"`
	ConfigPath   string            `json:"configPath"`
	ActivePreset string            `json:"activePreset"`
	Fields       []Field           `json:"fields"`
	Current      map[string]string `json:"current"`
	Snippet      string            `json:"snippet"`
	Presets      []Preset          `json:"presets"`
	ReadError    string            `json:"readError,omitempty"`
}

type ApplyRequest struct {
	Ecosystem string
	PresetID  string
	Values    map[string]string
	ValuesSet bool
}

type ValidationError struct {
	Reason string
}

func (e *ValidationError) Error() string {
	return e.Reason
}

func invalidf(format string, args ...any) error {
	return &ValidationError{Reason: fmt.Sprintf(format, args...)}
}

type presetDef struct {
	id     string
	name   string
	values map[string]string
}

type definition struct {
	id         string
	fields     []Field
	presets    []presetDef
	configPath func(Options) (string, error)
	read       func(Options) (map[string]string, error)
	write      func(Options, map[string]string) error
	snippet    func(map[string]string) string
}

func cloneValues(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func completeValues(def definition, current map[string]string) map[string]string {
	out := make(map[string]string, len(def.fields))
	for _, field := range def.fields {
		out[field.Key] = normalizeValue(field.Kind, current[field.Key])
	}
	return out
}

func normalizeValue(kind, value string) string {
	value = strings.TrimSpace(value)
	if kind != KindURLs {
		return value
	}
	return strings.Join(splitURLList(value), "\n")
}

func splitURLList(value string) []string {
	value = strings.ReplaceAll(value, ",", "\n")
	var out []string
	for _, part := range strings.Split(value, "\n") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func sameValues(def definition, left, right map[string]string) bool {
	for _, field := range def.fields {
		if normalizeValue(field.Kind, left[field.Key]) != normalizeValue(field.Kind, right[field.Key]) {
			return false
		}
	}
	return true
}

func matchPreset(def definition, current map[string]string) string {
	for _, preset := range def.presets {
		if sameValues(def, preset.values, current) {
			return preset.id
		}
	}
	return ""
}

func BuildState(def definition, current map[string]string, configPath string, readErr error) Ecosystem {
	fields := append([]Field(nil), def.fields...)
	state := Ecosystem{
		ID:         def.id,
		ConfigPath: configPath,
		Fields:     fields,
		Current:    completeValues(def, nil),
		Presets:    exportPresets(def),
	}
	if readErr != nil {
		state.ReadError = readErr.Error()
		state.Snippet = def.snippet(state.Current)
		return state
	}
	state.Current = completeValues(def, current)
	state.ActivePreset = matchPreset(def, state.Current)
	state.Snippet = def.snippet(state.Current)
	return state
}

func exportPresets(def definition) []Preset {
	presets := make([]Preset, 0, len(def.presets))
	for _, preset := range def.presets {
		values := completeValues(def, preset.values)
		presets = append(presets, Preset{
			ID:      preset.id,
			Name:    preset.name,
			Values:  values,
			Snippet: def.snippet(values),
		})
	}
	return presets
}

func ResolveValues(def definition, req ApplyRequest) (map[string]string, error) {
	if req.PresetID != "" && req.ValuesSet {
		return nil, invalidf("presetId and values are mutually exclusive")
	}
	if req.PresetID == "" && !req.ValuesSet {
		return nil, invalidf("presetId or values is required")
	}
	if req.PresetID != "" {
		for _, preset := range def.presets {
			if preset.id == req.PresetID {
				return completeValues(def, preset.values), nil
			}
		}
		return nil, invalidf("unknown preset %q for %s", req.PresetID, def.id)
	}
	for key := range req.Values {
		if !knownField(def, key) {
			return nil, invalidf("unknown field %q for %s", key, def.id)
		}
	}
	values := completeValues(def, req.Values)
	if err := validateValues(def, values); err != nil {
		return nil, err
	}
	return values, nil
}

func knownField(def definition, key string) bool {
	for _, field := range def.fields {
		if field.Key == key {
			return true
		}
	}
	return false
}

func validateValues(def definition, values map[string]string) error {
	for _, field := range def.fields {
		if err := validateField(def.id, field, values[field.Key]); err != nil {
			return err
		}
	}
	return nil
}
