package service

import (
	"os"
	"runtime"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/dto/request"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/utils/mirrors"
	"github.com/1Panel-dev/1Panel/pkg/platform/capabilities"
)

type IMirrorService interface {
	List() ([]mirrors.Ecosystem, error)
	Apply(req request.MirrorApplyReq) error
}

type MirrorService struct{}

func NewIMirrorService() IMirrorService {
	return &MirrorService{}
}

func (s *MirrorService) List() ([]mirrors.Ecosystem, error) {
	opts, err := s.options()
	if err != nil {
		return nil, err
	}
	items, err := mirrors.List(opts)
	if err != nil {
		return nil, err
	}
	if !capabilities.IsDarwin() {
		return items, nil
	}
	config := (&HomebrewService{}).loadMirrorConfig()
	items = append(items, mirrors.HomebrewState(map[string]string{
		"bottleDomain":  config.BottleDomain,
		"apiDomain":     config.APIDomain,
		"brewGitRemote": config.BrewGitRemote,
		"coreGitRemote": config.CoreGitRemote,
		"caskGitRemote": config.CaskGitRemote,
	}, constant.HomebrewMirrorConfigKey))
	return items, nil
}

func (s *MirrorService) Apply(req request.MirrorApplyReq) error {
	applyReq := mirrors.ApplyRequest{
		Ecosystem: strings.TrimSpace(req.Ecosystem),
		PresetID:  strings.TrimSpace(req.PresetID),
		Values:    req.Values,
		ValuesSet: req.Values != nil,
	}
	if applyReq.Ecosystem == mirrors.EcosystemHomebrew {
		return s.applyHomebrew(applyReq)
	}
	opts, err := s.options()
	if err != nil {
		return err
	}
	return mirrors.Apply(opts, applyReq)
}

func (s *MirrorService) applyHomebrew(req mirrors.ApplyRequest) error {
	if !capabilities.IsDarwin() {
		return &mirrors.ValidationError{Reason: "homebrew mirrors are only available on macOS"}
	}
	values, err := mirrors.ResolveHomebrew(req)
	if err != nil {
		return err
	}
	return (&HomebrewService{}).UpdateMirror(request.HomebrewMirrorUpdateReq{
		BottleDomain:  values["bottleDomain"],
		APIDomain:     values["apiDomain"],
		BrewGitRemote: values["brewGitRemote"],
		CoreGitRemote: values["coreGitRemote"],
		CaskGitRemote: values["caskGitRemote"],
	})
}

func (s *MirrorService) options() (mirrors.Options, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return mirrors.Options{}, err
	}
	return mirrors.Options{
		Home:             home,
		GOOS:             runtime.GOOS,
		GoEnvFile:        os.Getenv("GOENV"),
		XDGConfigHome:    os.Getenv("XDG_CONFIG_HOME"),
		CargoHome:        os.Getenv("CARGO_HOME"),
		PipConfigFile:    os.Getenv("PIP_CONFIG_FILE"),
		NpmrcFile:        os.Getenv("NPM_CONFIG_USERCONFIG"),
		GradleUserHome:   os.Getenv("GRADLE_USER_HOME"),
		DockerDaemonPath: constant.DaemonJsonPath,
	}, nil
}
