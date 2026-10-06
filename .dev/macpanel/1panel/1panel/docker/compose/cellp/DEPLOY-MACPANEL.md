# cellp 栈 — MacPanel 部署

单容器 **cellpd** 镜像内含控制面 API、Gateway 与 Dashboard。运行时按 version 动态拉起 celld 子进程（不在 compose 里单独定义）。

## 文件说明

| 文件 | 用途 |
|------|------|
| `docker-compose.yml` | 主 compose，仅 `cellpd` 服务 |
| `.env` | S3、端口、OTEL、Agent 证书路径等（勿提交 git） |
| `data/certs/elastic/` | AD-15 嵌入式 Node Agent 的 mTLS 证书（`ca.pem`、`agent-server*.pem`、`agent-client*.pem`） |
| `docker-compose.observability.yml` | **已废弃**，勿 `-f` 叠加 |
| `otel-collector-config.yaml` | **已废弃**（collector 在 `../otel/`） |
| `.env.observability.example` | OTEL 相关 env 片段，可复制进 `.env` |

## 前置：RustFS

cellpd 通过 `1panel-network` 访问 RustFS：

```env
S3_ENDPOINT=http://1Panel-rustfs-R3pn:9000
OFFSHOOT_S3_ENDPOINT=http://1Panel-rustfs-R3pn:9000
RUSTFS_ACCESS_KEY=rustfsadmin
RUSTFS_SECRET_KEY=<你的密钥>
```

Bucket（须已存在）：

```env
CELLD_BUCKET=s3://cellp-celld/demo-app
CELLP_ARTIFACTS_BUCKET=cellp-artifacts
OFFSHOOT_STORE=s3://cellp-offshoot
```

容器名 `1Panel-rustfs-R3pn` 随 1Panel 应用实例 ID 变化；以 MacPanel 容器列表中的 RustFS 容器名为准。

## 首次部署

```bash
cd /Users/mino/code/MacPanel/.dev/macpanel/1panel/1panel/docker/compose/cellp

# 若无 .env，从示例合并（S3/OTEL 段落见 .env.observability.example）
# 编辑 .env：RustFS 密钥、CELLP_ADMIN_TOKEN 等

docker compose up -d
```

## 端口与访问

| 端口 | 组件 |
|------|------|
| `8787` | Gateway（Host 选 version，见 cellp AD-12） |
| `8790` | Platform API |
| `5190` | Dashboard SPA |

```bash
curl -sf http://127.0.0.1:8790/v1/health
# Dashboard: http://127.0.0.1:5190 ，Bearer = CELLP_ADMIN_TOKEN
```

## 网络

cellpd 加入：

- compose 默认 network
- **`1panel-network`**（external）— 访问 RustFS 与 `otel-collector`

## 嵌入式 Node Agent（AD-15）

证书挂载：`./data/certs/elastic` → 容器内 `/data/certs/elastic`。

`.env` 中 `CELLP_ELASTIC_CERT_DIR=/data/certs/elastic` 及 `CELLP_AGENT_*_CERT_FILE` 路径均指向该目录。缺证书时 Agent 相关功能无法启动；部署前确认 PEM 文件齐全。

## 接入可观测（OTEL）

详见 [OBSERVABILITY.md](./OBSERVABILITY.md) 与 [../otel/README.md](../otel/README.md)。

要点：

1. 先启动 `../otel` 栈
2. 在 `.env` 设置 `CELLP_OTEL_COLLECTOR=http://otel-collector:4318`
3. `docker compose up -d --force-recreate cellpd`
4. **Trace 在 OpenObserve** `http://127.0.0.1:5080` 查看，不是 Dashboard Telemetry 页

## 升级镜像

编辑 `docker-compose.yml` 中 `image: ghcr.io/konghayao/cellp:<tag>`，然后：

```bash
docker compose pull
docker compose up -d
```

## MacPanel UI 注意

- **重启** cellpd 容器不会加载 `.env` 变更 → 使用 `docker compose up -d --force-recreate cellpd`
- 编排路径应仅为本目录 `docker-compose.yml`（单文件）

## 故障排查

| 问题 | 检查 |
|------|------|
| cellpd unhealthy | `docker compose logs cellpd`；8790 `/v1/health` |
| S3 连接失败 | RustFS 是否在 `1panel-network`；endpoint / 密钥 / bucket |
| 无 trace | OTEL 栈是否 up；`CELLP_OTEL_COLLECTOR`；是否 force-recreate |
| Agent TLS 错误 | `data/certs/elastic/` 文件是否完整、挂载是否 ro |
