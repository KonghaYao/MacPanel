# cellp 可观测（MacPanel）

OpenObserve + OTel Collector 在**独立目录**运行，不是本目录 compose 的 overlay：

**`../otel/`** — 详见 [../otel/README.md](../otel/README.md) · 总览 [../README.md](../README.md)

## 快速步骤

```bash
# 1. 部署 OTEL 栈
cd ../otel && cp .env.example .env && docker compose up -d

# 2. 在 cellp/.env 加入（或合并 .env.observability.example）
CELLP_OTEL_COLLECTOR=http://otel-collector:4318
CELLP_OTEL_BACKEND=jaeger

# 3. 重建 cellpd（仅 UI 重启不够）
cd ../cellp && docker compose up -d --force-recreate cellpd
```

## 连接说明

| 变量 | MacPanel 典型值 |
|------|-----------------|
| `CELLP_OTEL_COLLECTOR` | `http://otel-collector:4318` |
| `CELLP_OTEL_BACKEND` | `jaeger` 或 `memory`（查询门面；**不**代表 trace 存在 Jaeger） |

collector 容器名固定为 **`otel-collector`**，且挂载 **`1panel-network`**，cellpd 才能解析。

Collector 将 OTLP 转发至 `http://openobserve:5080/api/default`（Basic Auth 见 `otel/.env` 的 `OO_BASIC_AUTH`）。

## 在哪里看数据

| 数据 | 位置 |
|------|------|
| **Traces / Logs / Metrics（OTLP  ingest）** | OpenObserve **`http://127.0.0.1:5080`** |
| Prometheus 风格平台 metrics | cellpd **`http://127.0.0.1:8790/metrics`** |
| Dashboard Telemetry 调查页 | 依赖 Tempo/Jaeger 等 query API — **MacPanel OpenObserve 路径下请用 OO UI** |

## 已废弃（勿用）

| 文件 | 说明 |
|------|------|
| `docker-compose.observability.yml` | 已迁到 `../otel/` |
| `otel-collector-config.yaml`（本目录） | 使用 `../otel/otel-collector-config.yaml` |
| `docker-compose.cellp-bridge.yml` | 已删除；OTEL 仅单文件 compose |

## 故障排查

见 [../README.md#故障排查](../README.md#故障排查) 与 [../otel/README.md](../otel/README.md)。
