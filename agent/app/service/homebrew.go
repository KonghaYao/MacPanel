package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/dto/request"
	"github.com/1Panel-dev/1Panel/agent/app/dto/response"
	"github.com/1Panel-dev/1Panel/agent/app/task"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"github.com/1Panel-dev/1Panel/agent/utils/homebrew"
	"github.com/1Panel-dev/1Panel/pkg/platform/capabilities"
)

var errHomebrewUnsupported = errors.New("homebrew is only supported on macOS")

type HomebrewService struct{}

type IHomebrewService interface {
	GetStatus() (*response.HomebrewStatus, error)
	Search(req request.HomebrewSearchReq) ([]response.HomebrewSearchResult, error)
	List(req request.HomebrewListReq) (int64, []response.HomebrewPackage, error)
	Install(req request.HomebrewPackageReq) error
	Uninstall(req request.HomebrewPackageReq) error
	Upgrade(req request.HomebrewPackageReq) error
	Update(req request.HomebrewTaskReq) error
	GetMirror() (*response.HomebrewMirrorConfig, error)
	UpdateMirror(req request.HomebrewMirrorUpdateReq) error
	Doctor() (*response.HomebrewDoctorResult, error)
}

func NewIHomebrewService() IHomebrewService {
	return &HomebrewService{}
}

func (h *HomebrewService) requireDarwin() error {
	if !capabilities.IsDarwin() {
		return errHomebrewUnsupported
	}
	return nil
}

func (h *HomebrewService) loadMirrorConfig() homebrew.MirrorConfig {
	raw, err := settingRepo.GetValueByKey(constant.HomebrewMirrorConfigKey)
	if err != nil {
		return homebrew.MirrorConfig{}
	}
	return homebrew.ParseMirrorConfig(raw)
}

func (h *HomebrewService) saveMirrorConfig(config homebrew.MirrorConfig) error {
	value, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return settingRepo.UpdateOrCreate(constant.HomebrewMirrorConfigKey, string(value))
}

func (h *HomebrewService) newClient() (*homebrew.Client, error) {
	if err := h.requireDarwin(); err != nil {
		return nil, err
	}
	return homebrew.NewClient(h.loadMirrorConfig())
}

func (h *HomebrewService) GetStatus() (*response.HomebrewStatus, error) {
	if err := h.requireDarwin(); err != nil {
		return nil, err
	}
	res := &response.HomebrewStatus{}
	brewPath, err := homebrew.ResolveBrewPath()
	if err != nil {
		if errors.Is(err, homebrew.ErrNotInstalled) {
			res.IsExist = false
			return res, nil
		}
		return nil, err
	}
	res.IsExist = true
	res.BrewPath = brewPath

	client, err := homebrew.NewClient(h.loadMirrorConfig())
	if err != nil {
		return nil, err
	}
	version, err := client.Version()
	if err != nil {
		return nil, err
	}
	res.Version = version

	prefix, err := client.Prefix()
	if err != nil {
		return nil, err
	}
	res.Prefix = prefix
	res.CellarPath = prefix + "/Cellar"

	formulaItems, err := client.ListPackages("formula")
	if err != nil {
		return nil, err
	}
	caskItems, err := client.ListPackages("cask")
	if err != nil {
		return nil, err
	}
	res.FormulaCount = len(formulaItems)
	res.CaskCount = len(caskItems)
	return res, nil
}

func (h *HomebrewService) Search(req request.HomebrewSearchReq) ([]response.HomebrewSearchResult, error) {
	client, err := h.newClient()
	if err != nil {
		return nil, err
	}
	pkgType := req.Type
	if pkgType == "" {
		pkgType = "all"
	}
	var results []response.HomebrewSearchResult
	if pkgType == "all" || pkgType == "formula" {
		items, err := client.Search(req.Query, "formula")
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			results = append(results, response.HomebrewSearchResult{Name: item.Name, Type: "formula"})
		}
	}
	if pkgType == "all" || pkgType == "cask" {
		items, err := client.Search(req.Query, "cask")
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			results = append(results, response.HomebrewSearchResult{Name: item.Name, Type: "cask"})
		}
	}
	return results, nil
}

func (h *HomebrewService) List(req request.HomebrewListReq) (int64, []response.HomebrewPackage, error) {
	client, err := h.newClient()
	if err != nil {
		return 0, nil, err
	}
	pkgType := req.Type
	if pkgType == "" {
		pkgType = "all"
	}
	items, err := client.ListPackages(pkgType)
	if err != nil {
		return 0, nil, err
	}
	items = homebrew.FilterPackages(items, req.Info)
	pageItems, total := homebrew.PaginatePackages(items, req.Page, req.PageSize)
	records := make([]response.HomebrewPackage, 0, len(pageItems))
	for _, item := range pageItems {
		records = append(records, response.HomebrewPackage{
			Name:    item.Name,
			Version: item.Version,
			Type:    item.Type,
		})
	}
	return int64(total), records, nil
}

func (h *HomebrewService) Install(req request.HomebrewPackageReq) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("package name is required")
	}
	client, err := h.newClient()
	if err != nil {
		return err
	}
	pkgType := req.Type
	if pkgType == "" {
		pkgType = "formula"
	}
	taskItem, err := task.NewTaskWithOps(
		homebrew.FormatTaskName("install", req.Name),
		task.TaskInstall,
		task.TaskScopeSystem,
		req.TaskID,
		1,
	)
	if err != nil {
		return fmt.Errorf("new task for brew install failed, err: %v", err)
	}
	go func() {
		taskItem.AddSubTask("Install "+req.Name, func(t *task.Task) error {
			if err := client.Install(req.Name, pkgType, t); err != nil {
				return err
			}
			homebrew.InvalidatePackageListCache()
			return nil
		}, nil)
		_ = taskItem.Execute()
	}()
	return nil
}

func (h *HomebrewService) Uninstall(req request.HomebrewPackageReq) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("package name is required")
	}
	client, err := h.newClient()
	if err != nil {
		return err
	}
	pkgType := req.Type
	if pkgType == "" {
		pkgType = "formula"
	}
	taskItem, err := task.NewTaskWithOps(
		homebrew.FormatTaskName("uninstall", req.Name),
		task.TaskUninstall,
		task.TaskScopeSystem,
		req.TaskID,
		1,
	)
	if err != nil {
		return fmt.Errorf("new task for brew uninstall failed, err: %v", err)
	}
	go func() {
		taskItem.AddSubTask("Uninstall "+req.Name, func(t *task.Task) error {
			if err := client.Uninstall(req.Name, pkgType, t); err != nil {
				return err
			}
			homebrew.InvalidatePackageListCache()
			return nil
		}, nil)
		_ = taskItem.Execute()
	}()
	return nil
}

func (h *HomebrewService) Upgrade(req request.HomebrewPackageReq) error {
	client, err := h.newClient()
	if err != nil {
		return err
	}
	pkgType := req.Type
	if pkgType == "" {
		pkgType = "formula"
	}
	taskName := homebrew.FormatTaskName("upgrade", req.Name)
	if req.Name == "" {
		taskName = homebrew.FormatTaskName("upgrade", "all")
	}
	taskItem, err := task.NewTaskWithOps(taskName, task.TaskUpgrade, task.TaskScopeSystem, req.TaskID, 1)
	if err != nil {
		return fmt.Errorf("new task for brew upgrade failed, err: %v", err)
	}
	go func() {
		subTaskName := "Upgrade packages"
		if req.Name != "" {
			subTaskName = "Upgrade " + req.Name
		}
		taskItem.AddSubTask(subTaskName, func(t *task.Task) error {
			if err := client.Upgrade(req.Name, pkgType, t); err != nil {
				return err
			}
			homebrew.InvalidatePackageListCache()
			return nil
		}, nil)
		_ = taskItem.Execute()
	}()
	return nil
}

func (h *HomebrewService) Update(req request.HomebrewTaskReq) error {
	client, err := h.newClient()
	if err != nil {
		return err
	}
	taskItem, err := task.NewTaskWithOps("brew update", task.TaskUpdate, task.TaskScopeSystem, req.TaskID, 1)
	if err != nil {
		return fmt.Errorf("new task for brew update failed, err: %v", err)
	}
	go func() {
		taskItem.AddSubTask("Update Homebrew", func(t *task.Task) error {
			return client.Update(t)
		}, nil)
		_ = taskItem.Execute()
	}()
	return nil
}

func (h *HomebrewService) GetMirror() (*response.HomebrewMirrorConfig, error) {
	if err := h.requireDarwin(); err != nil {
		return nil, err
	}
	config := h.loadMirrorConfig()
	return &response.HomebrewMirrorConfig{
		BottleDomain:  config.BottleDomain,
		APIDomain:     config.APIDomain,
		BrewGitRemote: config.BrewGitRemote,
		CoreGitRemote: config.CoreGitRemote,
		CaskGitRemote: config.CaskGitRemote,
	}, nil
}

func (h *HomebrewService) UpdateMirror(req request.HomebrewMirrorUpdateReq) error {
	if err := h.requireDarwin(); err != nil {
		return err
	}
	config := homebrew.MirrorConfig{
		BottleDomain:  strings.TrimSpace(req.BottleDomain),
		APIDomain:     strings.TrimSpace(req.APIDomain),
		BrewGitRemote: strings.TrimSpace(req.BrewGitRemote),
		CoreGitRemote: strings.TrimSpace(req.CoreGitRemote),
		CaskGitRemote: strings.TrimSpace(req.CaskGitRemote),
	}
	if err := homebrew.ValidateMirrorConfig(config); err != nil {
		return err
	}
	if err := h.saveMirrorConfig(config); err != nil {
		return err
	}
	homebrew.InvalidatePackageListCache()
	return nil
}

func (h *HomebrewService) Doctor() (*response.HomebrewDoctorResult, error) {
	client, err := h.newClient()
	if err != nil {
		return nil, err
	}
	output, err := client.Doctor()
	if err != nil {
		return nil, err
	}
	return &response.HomebrewDoctorResult{Output: output}, nil
}
