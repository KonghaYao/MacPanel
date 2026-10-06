package homebrew

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/task"
	"github.com/1Panel-dev/1Panel/agent/buserr"
	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
)

const packageListCacheTTL = 15 * time.Second

var (
	packageListCacheMu     sync.RWMutex
	packageListCacheAt     time.Time
	packageListCacheMirror string
	packageListCacheFormula []PackageItem
	packageListCacheCask    []PackageItem
)

var packageNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9@./_-]*$`)

var ErrNotInstalled = errors.New("homebrew is not installed")

type MirrorConfig struct {
	BottleDomain  string `json:"bottleDomain"`
	APIDomain     string `json:"apiDomain"`
	BrewGitRemote string `json:"brewGitRemote"`
	CoreGitRemote string `json:"coreGitRemote"`
	CaskGitRemote string `json:"caskGitRemote"`
}

func (m MirrorConfig) Env() []string {
	var env []string
	if m.BottleDomain != "" {
		env = append(env, "HOMEBREW_BOTTLE_DOMAIN="+m.BottleDomain)
	}
	if m.APIDomain != "" {
		env = append(env, "HOMEBREW_API_DOMAIN="+m.APIDomain)
	}
	if m.BrewGitRemote != "" {
		env = append(env, "HOMEBREW_BREW_GIT_REMOTE="+m.BrewGitRemote)
	}
	if m.CoreGitRemote != "" {
		env = append(env, "HOMEBREW_CORE_GIT_REMOTE="+m.CoreGitRemote)
	}
	if m.CaskGitRemote != "" {
		env = append(env, "HOMEBREW_CASK_GIT_REMOTE="+m.CaskGitRemote)
	}
	return env
}

type Client struct {
	brewPath  string
	mirrorEnv []string
}

func ResolveBrewPath() (string, error) {
	candidates := []string{
		"/opt/homebrew/bin/brew",
		"/usr/local/bin/brew",
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, nil
		}
	}
	path, err := exec.LookPath("brew")
	if err != nil {
		return "", ErrNotInstalled
	}
	return path, nil
}

func NewClient(mirror MirrorConfig) (*Client, error) {
	brewPath, err := ResolveBrewPath()
	if err != nil {
		return nil, err
	}
	return &Client{
		brewPath:  brewPath,
		mirrorEnv: mirror.Env(),
	}, nil
}

func ValidatePackageName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return buserr.New("ErrCmdIllegal")
	}
	if cmd.CheckIllegal(name) || !packageNamePattern.MatchString(name) {
		return buserr.New("ErrCmdIllegal")
	}
	return nil
}

func ValidateMirrorURL(raw string) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return buserr.New("ErrCmdIllegal")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return buserr.New("ErrCmdIllegal")
	}
	if parsed.Host == "" {
		return buserr.New("ErrCmdIllegal")
	}
	return nil
}

func ValidateMirrorConfig(config MirrorConfig) error {
	fields := []string{
		config.BottleDomain,
		config.APIDomain,
		config.BrewGitRemote,
		config.CoreGitRemote,
		config.CaskGitRemote,
	}
	for _, field := range fields {
		if err := ValidateMirrorURL(field); err != nil {
			return err
		}
	}
	return nil
}

func ParseMirrorConfig(raw string) MirrorConfig {
	if strings.TrimSpace(raw) == "" {
		return MirrorConfig{}
	}
	var config MirrorConfig
	_ = json.Unmarshal([]byte(raw), &config)
	return config
}

func (c *Client) Path() string {
	return c.brewPath
}

func (c *Client) runWithStdout(timeout time.Duration, args ...string) (string, error) {
	if err := validateArgs(args...); err != nil {
		return "", err
	}
	opts := []cmd.Option{cmd.WithTimeout(timeout)}
	if len(c.mirrorEnv) > 0 {
		opts = append(opts, cmd.WithEnv(c.mirrorEnv...))
	}
	return cmd.NewCommandMgr(opts...).RunWithStdout(c.brewPath, args...)
}

func (c *Client) run(timeout time.Duration, taskItem *task.Task, args ...string) error {
	if err := validateArgs(args...); err != nil {
		return err
	}
	opts := []cmd.Option{cmd.WithTimeout(timeout)}
	if len(c.mirrorEnv) > 0 {
		opts = append(opts, cmd.WithEnv(c.mirrorEnv...))
	}
	if taskItem != nil {
		opts = append(opts, cmd.WithTask(*taskItem))
	}
	return cmd.NewCommandMgr(opts...).Run(c.brewPath, args...)
}

func validateArgs(args ...string) error {
	if len(args) == 0 {
		return buserr.New("ErrCmdIllegal")
	}
	for _, arg := range args {
		if cmd.CheckIllegal(arg) {
			return buserr.New("ErrCmdIllegal")
		}
	}
	return nil
}

func (c *Client) Version() (string, error) {
	output, err := c.runWithStdout(20*time.Second, "--version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (c *Client) Prefix() (string, error) {
	output, err := c.runWithStdout(20*time.Second, "--prefix")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func InvalidatePackageListCache() {
	packageListCacheMu.Lock()
	defer packageListCacheMu.Unlock()
	packageListCacheAt = time.Time{}
	packageListCacheFormula = nil
	packageListCacheCask = nil
	packageListCacheMirror = ""
}

func (c *Client) mirrorCacheKey() string {
	return strings.Join(c.mirrorEnv, "|")
}

func (c *Client) loadCachedPackageLists() (formulaItems, caskItems []PackageItem, ok bool) {
	mirrorKey := c.mirrorCacheKey()
	packageListCacheMu.RLock()
	defer packageListCacheMu.RUnlock()
	if packageListCacheAt.IsZero() ||
		time.Since(packageListCacheAt) >= packageListCacheTTL ||
		packageListCacheMirror != mirrorKey {
		return nil, nil, false
	}
	formulaCopy := make([]PackageItem, len(packageListCacheFormula))
	copy(formulaCopy, packageListCacheFormula)
	caskCopy := make([]PackageItem, len(packageListCacheCask))
	copy(caskCopy, packageListCacheCask)
	return formulaCopy, caskCopy, true
}

func (c *Client) storeCachedPackageLists(formulaItems, caskItems []PackageItem) {
	formulaCopy := make([]PackageItem, len(formulaItems))
	copy(formulaCopy, formulaItems)
	caskCopy := make([]PackageItem, len(caskItems))
	copy(caskCopy, caskItems)

	packageListCacheMu.Lock()
	defer packageListCacheMu.Unlock()
	packageListCacheFormula = formulaCopy
	packageListCacheCask = caskCopy
	packageListCacheAt = time.Now()
	packageListCacheMirror = c.mirrorCacheKey()
}

func (c *Client) loadAllPackageLists() (formulaItems, caskItems []PackageItem, err error) {
	if formulaItems, caskItems, ok := c.loadCachedPackageLists(); ok {
		return formulaItems, caskItems, nil
	}

	var formulaErr, caskErr error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		formulaItems, formulaErr = c.listByType("formula")
	}()
	go func() {
		defer wg.Done()
		caskItems, caskErr = c.listByType("cask")
	}()
	wg.Wait()
	if formulaErr != nil {
		return nil, nil, formulaErr
	}
	if caskErr != nil {
		return nil, nil, caskErr
	}
	c.storeCachedPackageLists(formulaItems, caskItems)
	return formulaItems, caskItems, nil
}

func (c *Client) ListPackages(pkgType string) ([]PackageItem, error) {
	formulaItems, caskItems, err := c.loadAllPackageLists()
	if err != nil {
		return nil, err
	}
	switch pkgType {
	case "cask":
		return caskItems, nil
	case "formula", "":
		return formulaItems, nil
	default:
		items := make([]PackageItem, 0, len(formulaItems)+len(caskItems))
		items = append(items, formulaItems...)
		items = append(items, caskItems...)
		return items, nil
	}
}

type PackageItem struct {
	Name    string
	Version string
	Type    string
}

func (c *Client) listByType(pkgType string) ([]PackageItem, error) {
	args := []string{"list", "--versions"}
	if pkgType == "cask" {
		args = append(args, "--cask")
	} else {
		args = append(args, "--formula")
	}
	output, err := c.runWithStdout(120*time.Second, args...)
	if err != nil {
		return nil, err
	}
	return parsePackageList(output, pkgType), nil
}

func parsePackageList(output, pkgType string) []PackageItem {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var items []PackageItem
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		item := PackageItem{
			Name: fields[0],
			Type: pkgType,
		}
		if len(fields) > 1 {
			item.Version = strings.Join(fields[1:], " ")
		}
		items = append(items, item)
	}
	return items
}

func (c *Client) Search(query, pkgType string) ([]PackageItem, error) {
	query = strings.TrimSpace(query)
	if err := ValidatePackageName(query); err != nil {
		return nil, err
	}
	args := []string{"search"}
	switch pkgType {
	case "formula":
		args = append(args, "--formula")
	case "cask":
		args = append(args, "--cask")
	}
	args = append(args, query)
	output, err := c.runWithStdout(60*time.Second, args...)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no formulae or casks found") {
			return nil, nil
		}
		return nil, err
	}
	return parseSearchOutput(output, pkgType), nil
}

func parseSearchOutput(output, pkgType string) []PackageItem {
	output = strings.TrimSpace(output)
	if output == "" {
		return nil
	}
	lines := strings.Split(output, "\n")
	var items []PackageItem
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "==>") {
			continue
		}
		for _, token := range strings.Fields(line) {
			if token == "" {
				continue
			}
			itemType := pkgType
			if itemType == "" || itemType == "all" {
				itemType = "formula"
			}
			items = append(items, PackageItem{
				Name: token,
				Type: itemType,
			})
		}
	}
	return items
}

func (c *Client) Install(name, pkgType string, taskItem *task.Task) error {
	if err := ValidatePackageName(name); err != nil {
		return err
	}
	args := []string{"install"}
	if pkgType == "cask" {
		args = append(args, "--cask")
	}
	args = append(args, name)
	return c.run(120*time.Minute, taskItem, args...)
}

func (c *Client) Uninstall(name, pkgType string, taskItem *task.Task) error {
	if err := ValidatePackageName(name); err != nil {
		return err
	}
	args := []string{"uninstall"}
	if pkgType == "cask" {
		args = append(args, "--cask")
	}
	args = append(args, name)
	return c.run(60*time.Minute, taskItem, args...)
}

func (c *Client) Upgrade(name, pkgType string, taskItem *task.Task) error {
	args := []string{"upgrade"}
	if name != "" {
		if err := ValidatePackageName(name); err != nil {
			return err
		}
		if pkgType == "cask" {
			args = append(args, "--cask")
		}
		args = append(args, name)
	} else if pkgType == "cask" {
		args = append(args, "--cask")
	}
	return c.run(120*time.Minute, taskItem, args...)
}

func (c *Client) Update(taskItem *task.Task) error {
	return c.run(120*time.Minute, taskItem, "update")
}

func (c *Client) Doctor() (string, error) {
	output, err := c.runWithStdout(10*time.Minute, "doctor")
	if err != nil {
		if output != "" {
			return output, nil
		}
		return "", err
	}
	return output, nil
}

func (c *Client) Info(name, pkgType string) (string, error) {
	if err := ValidatePackageName(name); err != nil {
		return "", err
	}
	args := []string{"info"}
	if pkgType == "cask" {
		args = append(args, "--cask")
	}
	args = append(args, name)
	return c.runWithStdout(60*time.Second, args...)
}

func FilterPackages(items []PackageItem, keyword string) []PackageItem {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return items
	}
	var filtered []PackageItem
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), keyword) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func PaginatePackages(items []PackageItem, page, pageSize int) ([]PackageItem, int) {
	total := len(items)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	start := (page - 1) * pageSize
	if start >= total {
		return []PackageItem{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return items[start:end], total
}

func MirrorConfigFromResponse(config responseMirrorConfig) MirrorConfig {
	return MirrorConfig{
		BottleDomain:  config.BottleDomain,
		APIDomain:     config.APIDomain,
		BrewGitRemote: config.BrewGitRemote,
		CoreGitRemote: config.CoreGitRemote,
		CaskGitRemote: config.CaskGitRemote,
	}
}

type responseMirrorConfig struct {
	BottleDomain  string
	APIDomain     string
	BrewGitRemote string
	CoreGitRemote string
	CaskGitRemote string
}

func FormatTaskName(action, name string) string {
	if name == "" {
		return fmt.Sprintf("brew %s", action)
	}
	return fmt.Sprintf("brew %s %s", action, name)
}
