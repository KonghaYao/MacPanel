package mirrors

import (
	"os"
)

func readNpm(opts Options) (map[string]string, error) {
	path, err := npmrcPath(opts)
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
	found := readAssignments(string(data), []string{"registry"})
	return map[string]string{"registry": found["registry"]}, nil
}

func writeNpm(opts Options, values map[string]string) error {
	return writeAssignmentFile(npmrcPath, opts, map[string]string{"registry": values["registry"]})
}

func readGo(opts Options) (map[string]string, error) {
	path, err := goEnvPath(opts)
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
	found := readAssignments(string(data), []string{"GOPROXY", "GOSUMDB"})
	return map[string]string{
		"goproxy": found["GOPROXY"],
		"gosumdb": found["GOSUMDB"],
	}, nil
}

func writeGo(opts Options, values map[string]string) error {
	return writeAssignmentFile(goEnvPath, opts, map[string]string{
		"GOPROXY": values["goproxy"],
		"GOSUMDB": values["gosumdb"],
	})
}

func writeAssignmentFile(pathFn func(Options) (string, error), opts Options, updates map[string]string) error {
	path, err := pathFn(opts)
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
	next := updateAssignments(existing, updates)
	if !assignmentFileMeaningful(next) {
		return removeFile(path)
	}
	return writeAtomic(path, []byte(next))
}
