# MacPanel + cellp 部署指南

在 MacPanel（macOS 版 1Panel）上运行 **cellp** 控制面与 **OpenObserve + OTel Collector** 可观测栈的实操文档。基于已验证的本地部署经验整理。

## 前提条件

| 项 | 要求 |
|----|------|
| **MacPanel** | 已安装并运行，`macpanel` 可访问 Docker Desktop |
| **Docker Desktop** | 已启动；MacPanel 不会替你安装 Docker |
| **RustFS** | 已在 `1panel-network` 上运行（1Panel 应用商店安装或手动 compose） |
| **cellp 源码** | 可选，用于构建/调试；运行时镜像来自 `ghcr.io/konghayao/cellp` |

RustFS 默认 endpoint（容器内）：`http://1Panel-rustfs-R3pn:9000`

需预先创建 bucket（名称与 cellp `.env` 一致）：

- `cellp-celld`
- `cellp-artifacts`
- `cellp-offshoot`

## 目录布局

```
compose/
├── README.md                 ← 本文（总览）
├── cellp/                    ← cellp 单栈（cellpd + Gateway + Dashboard）
│   ├── docker-compose.yml
│   ├── .env                  ← 运行时配置（gitignore）
│   ├── DEPLOY-MACPANEL.md    ← cellp 栈详细说明
│   ├── OBSERVABILITY.md      ← 如何接入 ../otel
│   └── data/certs/elastic/   ← AD-15 嵌入式 Node Agent mTLS 证书
└── otel/                     ← 独立可观测栈（OpenObserve + Collector）
    ├── docker-compose.yml    ← 唯一 compose 文件（勿再用 bridge overlay）
    ├── .env.example
    ├── README.md
    └── otel-collector-config.yaml
```

**重要：** OTEL 栈与 cellp 栈是**两个独立 compose 项目**，通过 Docker 外部网络 `1panel-network` 互通。不要把 OTEL 服务 merge 进 `cellp/docker-compose.yml`。

## 部署顺序

1. 确保 RustFS 与 `1panel-network` 存在
2. 部署 **cellp** 栈 → [cellp/DEPLOY-MACPANEL.md](./cellp/DEPLOY-MACPANEL.md)
3. 部署 **otel** 栈 → [otel/README.md](./otel/README.md)
4. 在 cellp `.env` 中配置 `CELLP_OTEL_COLLECTOR`，并 **force-recreate cellpd**
5. 按下方「验证」检查

## cellp 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Gateway | `8787` | Host 路由入口 |
| Platform API | `8790` | 控制面 REST API |
| Dashboard | `5190` | Web UI，登录 token = `CELLP_ADMIN_TOKEN` |

## 连接 cellp → OTEL

在 `compose/cellp/.env` 中：

```env
CELLP_OTEL_BACKEND=jaeger
CELLP_OTEL_COLLECTOR=http://otel-collector:4318
```

说明：

- `CELLP_OTEL_BACKEND=jaeger`（或 `memory`）控制 cellp **查询门面**的后端类型；MacPanel 生产路径下 trace **存储与检索在 OpenObserve**，不在 cellp Dashboard 的 Telemetry 调查页。
- `otel-collector` 容器同时加入 `otel-internal` 与 **`1panel-network`**，cellpd 通过容器名访问。

修改 `.env` 后必须重建 cellpd：

```bash
cd compose/cellp
docker compose up -d --force-recreate cellpd
```

## 查看 Trace

| 方式 | 地址 | 适用 |
|------|------|------|
| **OpenObserve UI** | `http://127.0.0.1:5080` | MacPanel OTEL 栈（推荐） |
| cellp Dashboard Telemetry | `http://127.0.0.1:5190` | 仅当后端为 Tempo/Jaeger 且 API 可达；**OpenObserve  setup 下不可用** |

登录 OpenObserve：`OO_ROOT_EMAIL` / `OO_ROOT_PASSWORD`（见 `otel/.env`）。

## 1Panel / MacPanel UI 操作注意

MacPanel 容器编排页通过 Docker 标签 `com.docker.compose.project.config_files` 记录创建容器时使用的 compose 文件列表。

### 旧 dual-file 导致重启失败

早期 OTEL 栈曾使用 `docker-compose.yml` + `docker-compose.cellp-bridge.yml` 双文件。**bridge 文件已删除并合并进单一 `docker-compose.yml`。**

若容器仍带旧标签，在 UI 点「重启」可能报错：

```text
open .../docker-compose.cellp-bridge.yml: no such file or directory
```

**修复（终端，在 otel 目录）：**

```bash
cd compose/otel
docker compose up -d --force-recreate
```

若 MacPanel 编排列表仍显示错误或路径不对，检查 Agent 数据库 `agent.db` 中 compose 记录的 `path` 字段是否仍指向已删除的 bridge 文件；在 UI 中删除旧编排记录后重新 `docker compose up -d`，或手动修正 DB 中的 path 为当前 `docker-compose.yml` 绝对路径。

### 修改 cellp 环境变量

UI「重启」**不会**重新加载 `.env` 中新变量。改 OTEL 相关 env 后务必：

```bash
cd compose/cellp && docker compose up -d --force-recreate cellpd
```

## 验证清单

```bash
# 1. cellp API
curl -sf http://127.0.0.1:8790/v1/health

# 2. Gateway
curl -sf http://127.0.0.1:8787/health 2>/dev/null || curl -sf http://127.0.0.1:8787/

# 3. OpenObserve
curl -sf http://127.0.0.1:5080/healthz

# 4. OTLP HTTP 端口（collector 监听）
curl -sf -o /dev/null -w '%{http_code}' http://127.0.0.1:4318/v1/traces
# 405 或无路由均表示端口可达

# 5. 容器状态
cd compose/otel && docker compose ps
cd compose/cellp && docker compose ps
```

产生 trace 后，在 OpenObserve UI → Traces 中按时间范围查询（非 Dashboard）。

## 故障排查

| 现象 | 原因 | 处理 |
|------|------|------|
| 1Panel 重启 OTEL 报 `no such file`（bridge yml） | 容器标签仍引用已删的 `docker-compose.cellp-bridge.yml` | `cd otel && docker compose up -d --force-recreate` |
| `otel-collector` 一直 **Created** 不启动 | 旧 compose 中 OpenObserve `healthcheck` 失败（distroless 无 wget），`depends_on: service_healthy` 阻塞 | 已修复：healthcheck 禁用，`depends_on: service_started`；拉最新 `docker-compose.yml` 并 force-recreate |
| cellp 发 OTLP 无数据 | `.env` 未设 collector 或未 recreate cellpd | 检查 `CELLP_OTEL_COLLECTOR`，force-recreate cellpd |
| Collector 401/403 写 OpenObserve | `OO_BASIC_AUTH` 与 root 凭据不一致 | 重新 `echo -n 'email:password' \| base64` 写入 `otel/.env` 并 recreate collector |
| cellpd 连不上 S3 | RustFS 未在 `1panel-network` 或 bucket 未建 | 检查 RustFS 容器网络与 bucket 名 |
| Dashboard 看不到 trace | 正常：OO 栈无 Tempo/Jaeger API | 改用 `http://127.0.0.1:5080` |

Collector 导出 endpoint（勿加尾部斜杠）：`http://openobserve:5080/api/default`，Basic Auth 来自 `OO_BASIC_AUTH`。

## 相关文档

- cellp 栈：[cellp/DEPLOY-MACPANEL.md](./cellp/DEPLOY-MACPANEL.md)
- 可观测接入：[cellp/OBSERVABILITY.md](./cellp/OBSERVABILITY.md)
- OTEL 栈详情：[otel/README.md](./otel/README.md)
- cellp 产品设计：[cellp 仓库 DESIGN.md / AD-14 OTEL-OBSERVABILITY.md](https://github.com/konghayao/cellp)
