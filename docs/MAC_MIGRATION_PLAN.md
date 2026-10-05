# MacPanel macOS 迁移计划

> **文档版本**: 1.0  
> **最后更新**: 2026-10-06  
> **项目路径**: `/Users/mino/code/MacPanel`  
> **状态**: 进行中（统一二进制、Makefile、CI 已落地，路径/Capability 等 Phase 1 项持续完善）

---

## 目录

1. [概述](#1-概述)
2. [已确认的用户决策](#2-已确认的用户决策)
3. [架构背景](#3-架构背景)
4. [迁移问题全量清单](#4-迁移问题全量清单)
5. [macOS 路径设计](#5-macos-路径设计)
6. [Docker 适配规则](#6-docker-适配规则)
7. [darwin 平台禁用能力清单](#7-darwin-平台禁用能力清单)
8. [分阶段路线图](#8-分阶段路线图)
9. [纯 Web 决策下的排除项](#9-纯-web-决策下的排除项)
10. [迁移路线对比（A vs B）](#10-迁移路线对比a-vs-b)
11. [Makefile 构建目标设计](#11-makefile-构建目标设计)
12. [各阶段验收标准](#12-各阶段验收标准)
13. [优先级总览](#13-优先级总览)
14. [关键受影响文件索引](#14-关键受影响文件索引)

---

## 1. 概述

MacPanel 是 1Panel 的分支项目，目标是在 **macOS 本地环境** 以 **纯 Web 面板** 形式运行（`localhost:9999`），保留 Core + Agent 双进程架构与 Vue3 嵌入式前端，同时 **彻底去除 Linux 专属依赖**（systemd、FHS 系统路径、iptables、sudo 等）。

本计划合并了四轮子代理扫描结果、架构分析与用户确认决策，作为后续开发的唯一权威参考。

**核心原则**:

- 不使用 `sudo`，所有路径必须位于用户可写目录
- Linux 专属功能在 darwin 上 **显式禁用**（非仅隐藏 UI），API 返回统一错误
- Docker 仅 **使用已有安装**，MacPanel 不触发 Docker Desktop 自动安装
- 产品形态为纯 Web，不引入 Electron/Tauri/Swift 原生壳

---

## 2. 已确认的用户决策

| 决策项 | 确认内容 |
|--------|----------|
| **产品形态** | 纯 Web，无原生 App（无 Electron/Tauri/Swift，无托盘/Dock/公证） |
| **部署场景** | Mac 本地运行，默认 `localhost:9999` |
| **Docker** | Mac 用户自行安装 Docker Desktop；MacPanel **不得** 触发自动安装，仅连接已有 Docker |
| **权限模型** | 整个项目 **不得** 在 macOS 上使用 `sudo`；路径设计必须正确，不得遗留 hack 代码 |
| **功能裁剪** | Linux 专属功能必须 **显式 DISABLE**（非仅隐藏菜单），API 返回统一错误响应 |
| **迁移路线** | 选择 **路线 A**（见第 10 节） |

---

## 3. 架构背景

### 3.1 现有架构

```
┌─────────────────────────────────────────────────────────┐
│                     浏览器 (Vue3 SPA)                    │
└──────────────────────────┬──────────────────────────────┘
                           │ HTTP/WebSocket
┌──────────────────────────▼──────────────────────────────┐
│              1panel-core (Go, 嵌入前端静态资源)           │
│  - 用户认证、权限、设置、升级                             │
│  - 通过 Unix Socket 代理请求到 Agent                      │
└──────────────────────────┬──────────────────────────────┘
                           │ /etc/1panel/agent.sock (当前)
┌──────────────────────────▼──────────────────────────────┐
│              1panel-agent (Go)                           │
│  - 容器、网站、防火墙、主机管理、计划任务等                │
│  - 直接操作 OS：systemd、iptables、/proc、fstab 等       │
└──────────────────────────┬──────────────────────────────┘
                           │ Docker API
┌──────────────────────────▼──────────────────────────────┐
│              Docker Engine / Docker Desktop              │
└─────────────────────────────────────────────────────────┘
```

### 3.2 当前状态

- **定位**: Linux 服务器面板，深度绑定 Linux 生态
- **路径约定**: FHS 标准（`/opt/1panel`、`/usr/local/bin/1pctl`、`/etc/1panel`）
- **服务管理**: 仅支持 systemd / OpenRC / SysVinit
- **防火墙**: iptables / nftables / firewalld 全栈
- **构建**: Makefile 中 `build_on_darwin` 交叉编译为 `GOOS=linux`，无法在 Mac 本地原生运行
- **CI**: 无 macOS 构建矩阵

### 3.3 目标状态（MacPanel on macOS）

- Core + Agent 均以 `GOOS=darwin` 原生编译运行
- 所有数据与配置位于 `~/Library/Application Support/MacPanel/`
- 通过 Capability API 向前端声明平台能力边界
- Docker 容器管理可用（依赖用户已安装的 Docker Desktop）
- Linux 专属模块 API 返回 `ErrNotSupportedOnMac`

---

## 4. 迁移问题全量清单

以下合并四轮子代理扫描结果，按类别去重整理。

### 4.1 硬编码 Linux 路径

| 问题 | 当前值 | 影响范围 |
|------|--------|----------|
| 安装根目录 | `/opt/1panel` | `core/init/viper/viper.go`、`agent/init/viper/viper.go`、迁移脚本 |
| 控制脚本 | `/usr/local/bin/1pctl` | `core/utils/ctl_conf/ctl_conf.go`、`agent/utils/ctl_conf/ctl_conf.go` |
| 敏感目录标记 | `/.1panel_clash` | 前端文件管理多处校验 |
| Swap 文件 | `/.1panel_swap` | `frontend/src/views/toolbox/device/swap/` |
| Tokenizer 示例 | `/opt/1panel/tokenizers/` | 多语言 i18n 占位符 |
| OpenResty 路径 | `/usr/local/openresty/nginx/` | `scripts/openresty-modules/diagnose-install.sh` |
| 节点代理 ID | `/etc/1panel/.nodeProxyID` | `agent/middleware/certificate.go` |

**迁移动作**: 引入 `internal/platform/paths`（或 `pkg/paths`）统一路径解析，按 `runtime.GOOS` 分支。

### 4.2 IPC Unix Socket

| 问题 | 当前值 | 文件 |
|------|--------|------|
| Agent 监听路径 | `/etc/1panel/agent.sock` | `agent/server/server.go` |
| Core 代理路径 | `/etc/1panel/agent.sock` | `core/utils/req_helper/proxy_local/req_to_local.go` |
| 常量定义 | `SockPath = "/etc/1panel/agent.sock"` | `core/init/proxy/proxy.go` |

**问题本质**: `/etc/1panel` 需要 root 权限创建，与「无 sudo」原则冲突。

**迁移动作**: Socket 迁移至 `{Base}/run/agent.sock`，Core 与 Agent 均从 paths 包读取。

### 4.3 systemd 专属服务管理

| 问题 | 说明 |
|------|------|
| Controller 初始化 | `core/utils/controller/controller.go`、`agent/utils/controller/controller.go` 仅探测 systemctl/rc-service/service |
| systemd 实现 | `core/utils/controller/manager/systemd.go`、`agent/utils/controller/manager/systemd.go` |
| Docker 服务控制 | `agent/app/service/docker.go` 中 `OperateDocker` 调用 systemctl |
| 升级流程 | `core/app/service/upgrade.go` 依赖 systemd 服务路径与 init 脚本 |
| CLI 命令 | `core/cmd/server/cmd/` 下 restore、reset、update 等均调用 systemctl |

**迁移动作**:

- Phase 1: darwin 上跳过 systemd 调用，改用手动脚本启停
- Phase 2: 可选 LaunchAgent plist（用户手动 `launchctl load`）

### 4.4 Linux 防火墙栈

| 组件 | 路径 |
|------|------|
| iptables 管理 | `agent/utils/firewall/iptables_helper/` |
| nftables 转发 | `agent/utils/firewall/forwarding/nftables.go` |
| Docker 防火墙守护 | `agent/utils/firewall/docker_guard/` |
| sysctl 转发 | `agent/utils/firewall/forwarding/sysctl.go` |
| ICMP ping 控制 | `agent/utils/firewall/ping.go`（读写 `/proc/sys/net/...`） |
| Fail2ban | `agent/utils/toolbox/fail2ban.go`、前端 `toolbox/fail2ban/` |
| 防火墙初始化 | `agent/init/firewall/firewall.go` |

**迁移动作**: darwin 上整体禁用，API 返回 `ErrNotSupportedOnMac`。

### 4.5 构建与升级（Linux-only 二进制）

| 问题 | 说明 |
|------|------|
| 当前 Makefile | `build_core_on_darwin` / `build_agent_on_darwin` 强制 `GOOS=linux GOARCH=amd64` |
| 在线升级 | `core/app/service/upgrade.go` 下载 Linux 包、替换 systemd 服务 |
| UPX 压缩 | `upx_bin` target 对 darwin 二进制可能不可用 |

**迁移动作**: 新增 `build_for_darwin` 原生编译；Phase 2 支持 darwin 升级包或禁用 Mac 在线升级。

### 4.6 Docker 依赖与 Socket 路径

| 问题 | 说明 |
|------|------|
| 默认 Socket | `unix:///var/run/docker.sock`（`agent/utils/docker/docker.go`） |
| Mac 实际路径 | Docker Desktop 使用 `~/.docker/run/docker.sock` |
| Docker 安装 | 当前可能尝试通过脚本安装 Docker 引擎 |
| systemd 控制 Docker | `OperateDocker` 使用 systemctl start/stop docker |

**迁移动作**: Socket 探测顺序 `~/.docker/run/docker.sock` → `/var/run/docker.sock`；禁用安装/卸载/systemd 控制。

### 4.7 主机/设备 Linux 工具

| 功能 | Linux 依赖 | 相关区域 |
|------|-----------|----------|
| 磁盘管理 | `/etc/fstab`、`lsblk`、`swapon` | agent 主机模块、前端 toolbox/device |
| SSH 服务配置 | `sshd_config`、`systemctl restart sshd` | `agent/app/service/ssh.go` |
| NTP/时区 | `timedatectl`、`chronyd` | agent 设置模块 |
| 进程监控 | 直接读 `/proc/<pid>/comm` | `agent/utils/websocket/process_data.go` |
| CPU 统计 | 读 `/proc/stat` | `agent/utils/psutil/cpu.go` |

**迁移动作**: darwin 禁用或改用 gopsutil 跨平台 API（Phase 1 监控类可先禁用）。

### 4.8 UID/权限（1000 vs 501）

| 问题 | 说明 |
|------|------|
| 网站默认 UID | 前端 `site-folder/index.vue` 默认 user/group = `1000` |
| Linux 惯例 | UID 1000 为首个普通用户 |
| macOS 惯例 | UID 501 为首个普通用户 |
| 容器卷权限 | Docker 挂载时 UID 不匹配导致读写失败 |

**迁移动作**: Phase 2 根据 `os.Getuid()` 或平台常量设置默认 UID（501 on darwin）。

### 4.9 前端 Linux 路径校验

| 问题 | 文件示例 |
|------|----------|
| `.1panel_clash` 禁止 | `frontend/src/views/host/file-management/` 多处 |
| `/opt/1panel` 示例路径 | 各语言 `lang/modules/*.ts` |
| 防火墙 UI | `frontend/src/views/host/firewall/` 全套 |
| 容器设置 Docker 安装提示 | `frontend/src/views/container/setting/index.vue` |

**迁移动作**: Phase 3 通过 Capability API 禁用菜单；Phase 1 起后端已拦截则前端可渐进适配。

### 4.10 Shell/Bash 默认值

| 问题 | 说明 |
|------|------|
| 脚本假设 | 多处 shell 脚本使用 bash 语法、Linux 工具链 |
| docker.sock 组权限 | Linux 将用户加入 docker 组；macOS 无此概念 |

**迁移动作**: Go 侧用 `runtime.GOOS` 分支；脚本类功能在 darwin 禁用或改写。

### 4.11–4.14 UI 项（纯 Web 决策下 mostly 排除或降级）

| 原扫描项 | 处理方式 |
|----------|----------|
| 原生系统托盘 | **排除**（纯 Web） |
| Dock 图标 / 菜单栏 | **排除** |
| macOS 系统通知 | **排除**（可用浏览器 Notification API 替代，非 P0） |
| _codesign / 公证 / Hardened Runtime_ | **排除** |
| 原生文件选择器 | **降级**为 Web `<input type="file">` |
| 全局快捷键 | **排除**或文档说明浏览器限制 |
| 深色模式跟随系统 | 可选 P3，非阻塞 |

### 4.15 Core Embed 构建顺序依赖

```
正确构建顺序:
1. frontend: npm install && npm run build:pro
2. 输出复制到 core/cmd/server/web/assets/
3. core: go build (embed 指令打包 assets)
4. agent: go build
```

**问题**: `build_on_local` 已有 `clean_assets → build_frontend → build_core_on_darwin` 顺序，但 cross-compile 为 linux 导致 Mac 无法本地验证。

**迁移动作**: 文档化 embed 顺序；`build_for_darwin` 改为 `GOOS=darwin`。

### 4.16 Go 1.26.6 + 无 macOS CI

| 问题 | 说明 |
|------|------|
| Go 版本 | `core/go.mod`、`agent/go.mod` 均为 `go 1.26.6` |
| CI | `.github/workflows/` 无 macos-latest 矩阵 |
| 风险 | darwin 编译错误无法自动发现 |

**迁移动作**: Phase 4 添加 GitHub Actions macOS 矩阵。

### 4.17 OpenResty Linux-only 诊断

| 问题 | 文件 |
|------|------|
| 诊断脚本 | `scripts/openresty-modules/diagnose-install.sh` |
| 硬编码路径 | `/usr/local/openresty/nginx/modules/1panel/` |
| 依赖 docker exec + nginx -t | 仅适用于 Linux 容器内 OpenResty |

**迁移动作**: darwin 禁用该诊断 API；OpenResty 模块编译仍可在 Docker 容器内进行（若网站功能保留）。

### 4.18 /proc 直接读取监控

| 路径 | 用途 |
|------|------|
| `/proc/<pid>/comm` | WebSocket 进程数据 |
| `/proc/stat` | CPU 统计 |
| `/proc/sys/net/ipv4/...` | 防火墙 ICMP |
| `/proc/uptime` | psutil 注释中提及的 LXC 问题 |

**迁移动作**: Phase 1 禁用依赖 `/proc` 的监控路径；长期可改用 `gopsutil` 跨平台实现。

---

## 5. macOS 路径设计

### 5.1 基础目录

```
~/Library/Application Support/MacPanel/     ← Base（用户可写，无需 sudo）
├── config/
│   └── 1pctl                               ← 原 /usr/local/bin/1pctl 配置
├── run/
│   └── agent.sock                          ← 原 /etc/1panel/agent.sock
└── 1panel/
    ├── conf/
    │   └── app.yaml                        ← 原 /opt/1panel/conf/app.yaml
    ├── db/                                 ← SQLite 等
    ├── log/
    ├── apps/
    ├── runtime/
    ├── backup/
    ├── cache/
    └── geo/                                ← GeoIP 等
```

### 5.2 路径映射表

| 用途 | Linux (1Panel) | macOS (MacPanel) |
|------|----------------|------------------|
| 安装根 | `/opt/1panel` | `{Base}/1panel` |
| 控制配置 | `/usr/local/bin/1pctl` | `{Base}/config/1pctl` |
| IPC Socket | `/etc/1panel/agent.sock` | `{Base}/run/agent.sock` |
| 节点 ID | `/etc/1panel/.nodeProxyID` | `{Base}/run/.nodeProxyID` |
| 系统服务 | `/etc/systemd/system/1panel-*.service` | `~/Library/LaunchAgents/com.macpanel.*.plist`（Phase 2 可选） |

### 5.3 路径包设计

建议新建 **`internal/platform/paths`**（或 **`pkg/paths`**）:

```go
// 伪代码示意
package paths

func BaseDir() string      // ~/Library/Application Support/MacPanel
func ConfigFile() string   // {Base}/config/1pctl
func SocketPath() string   // {Base}/run/agent.sock
func DataDir() string       // {Base}/1panel
func ConfDir() string       // {Base}/1panel/conf
// ... 其他子目录
```

**约束**:

- 所有硬编码路径必须替换为 paths 包调用
- Linux 路径作为 `GOOS=linux` 分支保留，不得删除
- 禁止在 darwin 代码路径中出现 `sudo`、`/opt`、`/etc/1panel`

### 5.4 1pctl 配置示例（macOS）

```ini
# ~/Library/Application Support/MacPanel/config/1pctl
BASE_DIR=/Users/<username>/Library/Application Support/MacPanel/1panel
ORIGINAL_PORT=9999
ORIGINAL_VERSION=v2.x.x
ORIGINAL_USERNAME=admin
ORIGINAL_PASSWORD=<generated>
ORIGINAL_ENTRANCE=
LANGUAGE=zh
PANEL_EDITION=standard
```

---

## 6. Docker 适配规则

### 6.1 Socket 探测顺序

```
1. ~/.docker/run/docker.sock          ← Docker Desktop for Mac (新版默认)
2. /var/run/docker.sock               ← 旧版或 Colima 等兼容路径
3. 环境变量 DOCKER_HOST（若已设置则优先）
```

实现位置: `agent/utils/docker/docker.go` 的 `NewDockerClient()` / `NewClient()`。

### 6.2 darwin 禁止的操作

| 操作 | 处理方式 |
|------|----------|
| 安装 Docker Engine | API 返回 `ErrNotSupportedOnMac` |
| 卸载 Docker | 同上 |
| systemctl start/stop docker | 跳过，提示用户通过 Docker Desktop 管理 |
| 修改 daemon.json 并重启 dockerd | 受限；可读取，写入需提示手动重启 Docker Desktop |
| 将用户加入 docker 组 | 不适用，跳过 |

### 6.3 前端表现

- 容器模块: Docker 未运行时显示「请确保 Docker Desktop 已启动」
- 设置页: 安装/卸载按钮 **置灰**，文案「请自行安装 Docker Desktop」
- Docker Socket 路径: 可选手动配置，默认自动探测

### 6.4 功能保留（Docker 可用时）

- 容器列表、创建、日志、终端
- 镜像管理、Compose
- 应用商店（基于 Docker 的部分）
- 网站运行时（OpenResty/Nginx 容器）

---

## 7. darwin 平台禁用能力清单

### 7.1 必须显式禁用的模块

| 模块 | 禁用原因 | Phase |
|------|----------|-------|
| 防火墙（iptables/nftables/firewalld） | Linux netfilter 栈 | P0 |
| Fail2ban | 依赖 iptables + systemd | P0 |
| 磁盘分区 / fstab / swap | macOS 无 fstab | P0 |
| NTP / 时区同步（timedatectl） | 无 timedatectl | P0 |
| SSH 服务配置（sshd） | macOS 使用 system sshd，配置路径不同 | P0 |
| OpenResty 诊断脚本 | Linux 容器内路径 | P0 |
| /proc 监控（进程 comm、CPU stat） | macOS 无 /proc | P0（禁用）/ P2（gopsutil） |
| Docker 安装/卸载 | 用户自行管理 | P0 |
| systemd 服务注册 | 改用 LaunchAgent（Phase 2）或手动脚本 | P0 |
| 在线升级（初期） | 无 darwin 升级包 | P1 |
| ICMP sysctl 控制 | /proc/sys 不可用 | P0 |
| IP 转发 / sysctl 转发 | 同上 | P0 |

### 7.2 API 设计

**统一错误**:

```go
var ErrNotSupportedOnMac = errors.New("this feature is not supported on macOS")
```

或使用 `buserr` 包装，HTTP 状态码建议 **501 Not Implemented** 或 **403 Forbidden**，响应体:

```json
{
  "code": 501,
  "message": "此功能在 macOS 上不可用",
  "reason": "NOT_SUPPORTED_ON_DARWIN"
}
```

**Capability 端点**（建议 Phase 1 引入）:

```
GET /api/v2/platform/capabilities
```

响应示例:

```json
{
  "os": "darwin",
  "arch": "arm64",
  "features": {
    "firewall": false,
    "fail2ban": false,
    "disk_management": false,
    "fstab": false,
    "swap": false,
    "ntp_sync": false,
    "sshd_config": false,
    "docker_install": false,
    "docker_manage": true,
    "online_upgrade": false,
    "proc_monitoring": false,
    "openresty_diagnose": false
  }
}
```

### 7.3 前端联动

- 应用启动时拉取 `/platform/capabilities`
- `features.<name> === false` 的菜单项: **disabled + tooltip** 说明原因
- 禁止仅 CSS 隐藏（`display: none`），须使用路由守卫或菜单配置过滤

---

## 8. 分阶段路线图

### Phase 1（P0）— 可在 Mac 上无 sudo 启动

**目标**: 开发者能在 macOS 上编译、启动 Core + Agent，访问 `localhost:9999`，管理 Docker 容器。

| 任务 | 详情 |
|------|------|
| 路径抽象层 | 新建 `paths` 包，替换所有硬编码 Linux 路径 |
| IPC 迁移 | Socket 迁至 `{Base}/run/agent.sock`，Core/Agent 同步 |
| Makefile | 新增 `build_for_darwin`（`GOOS=darwin`），保留 `build_for_linux GOARCH=amd64\|arm64` |
| Embed 构建文档 | 在本文档及 Makefile 注释中明确 frontend → assets → go build 顺序 |
| Docker use-only | Socket 自动探测，禁用安装/卸载/systemd |
| Capability 注册表 | 后端 capability 包 + API 端点，各禁用模块注册 |
| 手动启动脚本 | `scripts/mac/start.sh`：检查 Docker、创建目录、启动 Agent + Core |
| 清理 sudo 引用 | 审计并移除 darwin 代码路径中所有 sudo 调用 |

**不在 Phase 1 范围**:

- LaunchAgent 自动安装
- 在线升级
- 前端菜单全面适配（可显示后端错误）

### Phase 2（P1）— 运维体验与权限

| 任务 | 详情 |
|------|------|
| LaunchAgent plist | `scripts/mac/com.macpanel.core.plist`、`com.macpanel.agent.plist`，文档说明手动 `launchctl load` |
| 升级包 darwin 支持 | 构建 darwin 发布包，或明确禁用 Mac 在线升级并在 UI 提示 |
| UID 501 适配 | 网站/运行时默认 user/group 改为平台感知（501 on darOS） |
| gopsutil 监控 | 可选：进程/CPU 监控改用 gopsutil 替代 /proc |

### Phase 3（P2）— 前端与浏览器体验

| 任务 | 详情 |
|------|------|
| 前端 Capability 联动 | 菜单 disable、路由守卫、统一「macOS 不可用」组件 |
| Safari 兼容性 | 文件上传、剪贴板、Passkey/HTTPS 相关说明与测试 |
| 终端键盘 | 文档记录 Option/Command 键与 xterm.js 差异 |
| macOS 敏感路径校验 | 前端文件管理增加 `.Trashes`、`Library/Keychains` 等警告 |

### Phase 4（P2，可选）— CI 与文档

| 任务 | 详情 |
|------|------|
| GitHub Actions | `macos-latest` + `ubuntu-latest` 矩阵，`build_for_darwin` + 基础测试 |
| 环境文档 | Go 1.26.6+、Node 18+、Docker Desktop 安装与配置指南 |
| 发布流程 | darwin arm64/amd64 二进制打包脚本 |

---

## 9. 纯 Web 决策下的排除项

以下功能在路线 A（纯 Web）中 **明确不做**:

| 排除项 | 原扫描来源 | 原因 |
|--------|-----------|------|
| Electron / Tauri 原生壳 | UI 项 #11 | 用户决策纯 Web |
| Swift 原生 App | UI 项 #11 | 同上 |
| 菜单栏托盘图标 | UI 项 #12 | 无原生壳 |
| Dock 图标定制 | UI 项 #12 | 浏览器访问即可 |
| macOS 系统通知（原生） | UI 项 #13 | 可用 Web Notification 替代，非必须 |
| Code Signing | UI 项 #14 | 无 .app bundle |
| Notarization / Hardened Runtime | UI 项 #14 | 同上 |
| 自动安装 LaunchAgent（Phase 1） | 服务管理 | 手动脚本优先 |
| Docker Desktop 自动安装 | Docker 项 | 用户决策 |
| sudo 任何场景 | 权限项 | 用户决策 |
| 防火墙 / Fail2ban 移植到 pf | 防火墙项 | 成本过高，直接禁用 |

---

## 10. 迁移路线对比（A vs B）

用户已选择 **路线 A**。

| 维度 | 路线 A：纯 Web 本地面板（已选） | 路线 B：原生壳 + 系统集成 |
|------|--------------------------------|---------------------------|
| **产品形态** | 浏览器访问 `localhost:9999` | Electron/Tauri/Swift 包装 WebView |
| **开发成本** | 低：仅后端平台适配 + 前端 Capability | 高：原生壳、打包、签名、更新通道 |
| **Docker** | 用户使用已有 Docker Desktop | 可集成 Docker Desktop 检测与引导 |
| **服务管理** | 手动脚本 / 可选 LaunchAgent | 原生菜单栏启停、自动注册 LaunchAgent |
| **系统通知** | 浏览器 API（可选） | 原生 NSUserNotification |
| **签名/公证** | 不需要 | 必须（分发 .app 时） |
| **托盘/Dock** | 无 | 有 |
| **Linux 功能禁用** | Capability API 统一处理 | 同样需 Capability，另加原生 UI |
| **跨平台一致性** | 与 Linux 版 Web UI 一致 | 原生控件可能不一致 |
| **维护负担** | 低 | 高（三端：Go + Vue + 原生壳） |
| **适用场景** | 开发者本地 Docker 管理 | 面向非技术用户的「Mac 应用」体验 |

**选择 A 的理由（用户决策）**:

- 最小化迁移范围，聚焦 Core/Agent 平台解耦
- 避免签名、公证、原生更新等 Mac 分发复杂度
- 与 1Panel 上游 Web 架构保持一致，便于合并

---

## 11. Makefile 构建目标设计

### 11.1 当前状态

```makefile
# 现有 target（节选）
build_frontend          # npm build → assets
build_core_on_linux     # GOOS=$(GOOS) 本地 Linux 编译
build_agent_on_linux
build_core_on_darwin    # GOOS=linux GOARCH=amd64 交叉编译（非 Mac 原生）
build_agent_on_darwin
build_all               # frontend + linux core + linux agent
build_on_local          # clean_assets + frontend + darwin-cross-to-linux
```

### 11.2 目标设计

```makefile
# 新增/改造 target

# Mac 本地原生编译（Phase 1 核心）
build_for_darwin: build_frontend
	cd $(CORE_PATH) && CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) $(GOBUILD) \
		-trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(CORE_NAME) $(CORE_MAIN)
	cd $(AGENT_PATH) && CGO_ENABLED=0 GOOS=darwin GOARCH=$(GOARCH) $(GOBUILD) \
		-trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(AGENT_NAME) $(AGENT_MAIN)

# Linux 服务器编译（保留，支持架构参数）
build_for_linux: build_frontend
	cd $(CORE_PATH) && CGO_ENABLED=0 GOOS=linux GOARCH=$(or $(GOARCH),amd64) $(GOBUILD) \
		-trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(CORE_NAME) $(CORE_MAIN)
	cd $(AGENT_PATH) && CGO_ENABLED=0 GOOS=linux GOARCH=$(or $(GOARCH),amd64) $(GOBUILD) \
		-trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(AGENT_NAME) $(AGENT_MAIN)

# 语义澄清
# build_on_local  → 改为调用 build_for_darwin（在 Mac 上本地开发）
# build_all       → 等价于 build_for_linux（CI/发布 Linux 版）
```

### 11.3 使用示例

```bash
# 在 Mac 上开发
make build_for_darwin
./scripts/mac/start.sh

# 交叉编译 Linux amd64 发布包
make build_for_linux GOARCH=amd64

# 交叉编译 Linux arm64（树莓派/ARM 服务器）
make build_for_linux GOARCH=arm64
```

### 11.4 Embed 构建顺序（强制）

```
1. make clean_assets          # 可选，清除旧 assets
2. make build_frontend        # Vue → core/cmd/server/web/assets/
3. make build_for_darwin      # go embed 打包
4. 运行二进制
```

**注意**: 跳过 `build_frontend` 直接 `go build` 将导致嵌入过期静态资源。

### 11.5 统一二进制 `macpanel`（已实现）

MacPanel 在 macOS 上使用 **单一可执行文件** 同时启动 Core 与 Agent，避免手动管理双进程。

```
┌─────────────────────────────────────────┐
│           macpanel (单进程)              │
│  ┌─────────────┐    ┌─────────────────┐ │
│  │ goroutine   │    │ main goroutine  │ │
│  │ agent.Start │───▶│ core.Start      │ │
│  └─────────────┘    └─────────────────┘ │
│         │                    │          │
│         └──── agent.sock ────┘          │
└─────────────────────────────────────────┘
```

| 组件 | 路径 |
|------|------|
| 入口 | `cmd/macpanel/main.go` |
| 模块链接 | 根目录 `go.work`（core / agent / cmd/macpanel / pkg/platform） |
| 构建产物 | `build/macpanel` |
| 启动脚本 | `scripts/mac/start.sh` → `exec build/macpanel` |

**子命令（向后兼容）**:

- `macpanel` — 默认：Agent goroutine + 等待 socket + Core 阻塞
- `macpanel core` — 仅 Core
- `macpanel agent` — 仅 Agent

**构建命令**:

```bash
make build_macpanel          # frontend + darwin 统一二进制
./scripts/mac/start.sh       # 启动
```

### 11.6 GitHub Actions macOS 构建（已实现）

工作流：`.github/workflows/build-macos.yml`

| 项 | 说明 |
|----|------|
| 触发 | push 到 `main` / `dev-v2`、PR、`workflow_dispatch` |
| 矩阵 | `macos-latest` (arm64)、`macos-13` (amd64) |
| 步骤 | checkout → Go 1.26.6 → Node 20 → `npm ci` + `build:pro` → `go build cmd/macpanel` |
| 产物 | Artifact `macpanel-darwin-{arch}` |
| 发布 | tag `v*` 时自动 attach 双架构二进制到 Release |

---

## 12. 各阶段验收标准

### Phase 1 验收标准

- [ ] `make build_for_darwin` 在 macOS（arm64 或 amd64）上零错误完成
- [ ] 无需 sudo 即可创建 `{Base}` 下全部目录
- [ ] `./scripts/mac/start.sh` 启动 Core + Agent，浏览器访问 `http://localhost:9999` 成功
- [ ] 登录、Dashboard 可正常加载
- [ ] Docker Desktop 运行时可列出容器、查看日志
- [ ] 调用防火墙/磁盘/fstab/sshd API 返回 `NOT_SUPPORTED_ON_DARWIN`
- [ ] 代码库 darwin 路径无 `/opt`、`/etc/1panel`、`/usr/local/bin/1pctl` 硬编码（paths 包除外 Linux 分支）
- [ ] 代码库无 `sudo` 调用（darwin build tag 路径）
- [ ] `GET /api/v2/platform/capabilities` 返回正确 feature map

### Phase 2 验收标准

- [ ] LaunchAgent plist 模板存在，文档说明安装步骤
- [ ] `launchctl load` 后重启登录可自动启动 MacPanel
- [ ] 网站创建默认 UID/GID 在 macOS 上为 501（或当前用户 UID）
- [ ] 在线升级: darwin 包可用 **或** UI 明确提示「macOS 请手动升级」
- [ ] gopsutil 进程列表在 macOS 上可用（若实现）

### Phase 3 验收标准

- [ ] 防火墙、Fail2ban、磁盘、swap 等菜单在 macOS 上 disabled 且不可路由进入
- [ ] Safari 可完成文件上传、基本操作
- [ ] 终端 WebSocket 在 Safari/Chrome 可用，键盘说明文档存在
- [ ] 文件管理尝试访问 Keychains 等路径时有警告

### Phase 4 验收标准

- [ ] GitHub Actions macOS job 绿色
- [ ] README/docs 含 Mac 环境要求与快速开始
- [ ] darwin arm64 二进制可下载或 CI artifact 可用

---

## 13. 优先级总览

```
P0 (Phase 1) — 阻塞性，必须首先完成
├── paths 路径抽象层
├── IPC Socket 用户目录迁移
├── build_for_darwin Makefile
├── Docker socket 探测 + 禁用安装
├── Capability 注册表 + API
├── 禁用模块统一 ErrNotSupportedOnMac
├── scripts/mac/start.sh
└── 移除 darwin sudo 与硬编码路径

P1 (Phase 2) — 运维与权限
├── LaunchAgent plist（可选安装）
├── UID 501 适配
├── darwin 升级包或禁用在线升级
└── gopsutil 替代 /proc（可选）

P2 (Phase 3) — 体验完善
├── 前端 Capability 菜单 disable
├── Safari 兼容测试与文档
├── 终端键盘说明
└── macOS 敏感路径前端校验

P2 (Phase 4, 可选) — 工程化
├── GitHub Actions macos-latest
└── 环境文档与发布脚本
```

---

## 14. 关键受影响文件索引

### 14.1 路径与配置

| 文件 | 改动类型 |
|------|----------|
| `core/utils/ctl_conf/ctl_conf.go` | 改用 paths 包 |
| `agent/utils/ctl_conf/ctl_conf.go` | 同上 |
| `core/init/viper/viper.go` | `/opt/1panel` → paths |
| `agent/init/viper/viper.go` | 同上 |
| **新建** `core/internal/platform/paths/` | 路径解析 |
| **新建** `agent/internal/platform/paths/` | 或共享 `pkg/paths` |

### 14.2 IPC / 代理

| 文件 | 改动类型 |
|------|----------|
| `agent/server/server.go` | Socket 路径 + 目录权限 |
| `core/utils/req_helper/proxy_local/req_to_local.go` | Socket 路径 |
| `core/init/proxy/proxy.go` | `SockPath` 常量 |
| `agent/middleware/certificate.go` | `.nodeProxyID` 路径 |

### 14.3 服务管理 / CLI

| 文件 | 改动类型 |
|------|----------|
| `core/utils/controller/controller.go` | darwin 分支 / LaunchAgent |
| `agent/utils/controller/controller.go` | 同上 |
| `core/utils/controller/manager/systemd.go` | darwin 禁用 |
| `agent/utils/controller/manager/systemd.go` | 同上 |
| `core/cmd/server/cmd/root.go` | 平台感知 |
| `core/cmd/server/cmd/restore.go` | 禁用或改写 |
| `core/cmd/server/cmd/update.go` | 禁用或改写 |
| `core/cmd/server/cmd/reset.go` | 禁用或改写 |
| `core/app/service/upgrade.go` | darwin 升级逻辑 |

### 14.4 Docker

| 文件 | 改动类型 |
|------|----------|
| `agent/utils/docker/docker.go` | Socket 探测顺序 |
| `agent/app/service/docker.go` | 禁用安装/systemctl |
| `frontend/src/views/container/setting/index.vue` | 置灰 + 文案 |
| `frontend/src/views/container/dashboard/index.vue` | Docker 未运行提示 |

### 14.5 防火墙 / 网络

| 文件 | 改动类型 |
|------|----------|
| `agent/init/firewall/firewall.go` | darwin 跳过初始化 |
| `agent/utils/firewall/`（整个目录） | capability 禁用 |
| `agent/utils/toolbox/fail2ban.go` | 禁用 |
| `frontend/src/views/host/firewall/`（整个目录） | Phase 3 disable |
| `frontend/src/views/toolbox/fail2ban/` | Phase 3 disable |

### 14.6 主机 / 监控

| 文件 | 改动类型 |
|------|----------|
| `agent/utils/websocket/process_data.go` | 禁用 /proc 或 gopsutil |
| `agent/utils/psutil/cpu.go` | 禁用 /proc/stat 或 gopsutil |
| `agent/utils/psutil/process.go` | 注释中的 /proc 问题 |
| `agent/app/service/ssh.go` | darwin 禁用 |
| `frontend/src/views/toolbox/device/swap/index.vue` | 禁用 |

### 14.7 构建 / 脚本

| 文件 | 改动类型 |
|------|----------|
| `Makefile` | 新增 build_for_darwin / build_for_linux |
| **新建** `scripts/mac/start.sh` | Mac 启动脚本 |
| **新建** `scripts/mac/com.macpanel.*.plist` | Phase 2 LaunchAgent |
| `scripts/openresty-modules/diagnose-install.sh` | darwin API 禁用 |

### 14.8 前端 i18n / 校验

| 文件 | 改动类型 |
|------|----------|
| `frontend/src/lang/modules/zh.ts`（及所有语言） | macOS 文案、Docker 提示 |
| `frontend/src/views/host/file-management/` | `.1panel_clash`、敏感路径 |
| `frontend/src/views/website/website/config/basic/site-folder/index.vue` | UID 501 |
| **新建** `frontend/src/api/platform/capabilities.ts` | Capability API 客户端 |
| **新建** `frontend/src/composables/usePlatformCapabilities.ts` | 菜单 disable 逻辑 |

### 14.9 Capability / API（新建）

| 文件 | 改动类型 |
|------|----------|
| **新建** `core/app/api/v2/platform.go` | Capability 端点 |
| **新建** `core/constant/capability.go` | Feature 常量 |
| **新建** `agent/constant/capability.go` | 或共享 pkg |
| **新建** `core/buserr/platform.go` | `ErrNotSupportedOnMac` |

### 14.10 CI / 文档

| 文件 | 改动类型 |
|------|----------|
| `.github/workflows/build-macos.yml` | macOS CI 矩阵（已实现） |
| `cmd/macpanel/main.go` | 统一二进制入口（已实现） |
| `go.work` | 多模块工作区（已实现） |
| `docs/MAC_MIGRATION_PLAN.md` | 本文档 |
| **新建** `docs/MAC_QUICKSTART.md` | Phase 4 快速开始（可选） |

---

## 附录 A：`scripts/mac/start.sh`（已实现）

当前实现为轻量包装，直接 exec 统一二进制（目录/bootstrap 由 `paths.Bootstrap` 在进程内完成）:

```bash
#!/bin/bash
exec "$(dirname "$0")/../../build/macpanel"
```

使用前请先 `make build_macpanel`。首次启动会在 `~/Library/Application Support/MacPanel/` 创建目录与默认 `config/1pctl`。

---

## 附录 B：术语表

| 术语 | 含义 |
|------|------|
| Core | `1panel-core`，主控进程，嵌入 Vue 前端 |
| Agent | `1panel-agent`，执行进程，操作 Docker 与 OS |
| 1pctl | 控制配置文件，存储 BASE_DIR、端口、凭据等 |
| Capability | 平台能力声明，告知前端哪些功能可用 |
| Base | `~/Library/Application Support/MacPanel/` |
| FHS | Filesystem Hierarchy Standard，Linux 目录规范 |
| LaunchAgent | macOS 用户级守护进程机制 |

---

*本文档将随迁移进展更新。代码实现时请以此为准，如有冲突以用户最新决策为准。*
