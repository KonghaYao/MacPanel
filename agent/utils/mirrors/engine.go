package mirrors

func List(opts Options) ([]Ecosystem, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	defs := definitions()
	out := make([]Ecosystem, 0, len(defs))
	for _, def := range defs {
		path, err := def.configPath(opts)
		if err != nil {
			return nil, err
		}
		current, readErr := def.read(opts)
		out = append(out, BuildState(def, current, path, readErr))
	}
	return out, nil
}

func Apply(opts Options, req ApplyRequest) error {
	if err := opts.validate(); err != nil {
		return err
	}
	var def definition
	found := false
	for _, item := range definitions() {
		if item.id == req.Ecosystem {
			def = item
			found = true
			break
		}
	}
	if !found || def.write == nil {
		return invalidf("unknown ecosystem %q", req.Ecosystem)
	}
	values, err := ResolveValues(def, req)
	if err != nil {
		return err
	}
	if err := validateValues(def, values); err != nil {
		return err
	}
	return def.write(opts, values)
}
