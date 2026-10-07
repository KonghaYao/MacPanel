package viper

import (
	"bytes"
	"fmt"
	"path"
	"strconv"

	"github.com/1Panel-dev/1Panel/agent/cmd/server/conf"
	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/files"
	"github.com/1Panel-dev/1Panel/agent/utils/xpack"
	"github.com/1Panel-dev/1Panel/pkg/platform/paths"
	"github.com/1Panel-dev/1Panel/pkg/platform/version"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func Init() {
	if err := paths.Bootstrap(version.Version); err != nil {
		panic(err)
	}
	mode := ""
	fileOp := files.NewFileOp()
	v := viper.NewWithOptions()
	v.SetConfigType("yaml")

	config := global.ServerConfig{}
	if err := yaml.Unmarshal(conf.AppYaml, &config); err != nil {
		panic(err)
	}
	if config.Base.Mode != "" {
		mode = config.Base.Mode
	}
	devConfFile := paths.LinuxDevConfFile()
	if paths.IsDarwin() {
		devConfFile = path.Join(paths.ConfDir(), "app.yaml")
	}
	useDevConfig := mode == "dev" && fileOp.Stat(devConfFile)
	if useDevConfig {
		v.SetConfigName("app")
		v.AddConfigPath(paths.ConfDir())
		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	} else {
		reader := bytes.NewReader(conf.AppYaml)
		if err := v.ReadConfig(reader); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&global.CONF); err != nil {
			panic(err)
		}
	})
	serverConfig := global.ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}

	global.CONF = serverConfig
	if paths.IsDarwin() && !useDevConfig {
		global.CONF.Base.Mode = "stable"
		if global.CONF.Log.Level == "" || global.CONF.Log.Level == "debug" {
			global.CONF.Log.Level = "info"
		}
	}

	initBaseInfo()
	global.Viper = v
}

func initBaseInfo() {
	nodeInfo, err := xpack.MultiNodeProvider.LoadNodeInfo(true)
	if err != nil {
		panic(err)
	}
	global.CONF.Base.InstallDir = nodeInfo.BaseDir
	if !global.IsMaster {
		global.CONF.Base.Port = strconv.FormatUint(uint64(nodeInfo.NodePort), 10)
		if nodeInfo.NodePort == 0 {
			global.CONF.Base.Port = "9999"
		}
	}
}
