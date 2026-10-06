# MacPanel OTEL 栈（OpenObserve + Collector）

独立、面向生产的可观测 ingest 栈，供 cellp（AD-14）使用。**不是** `compose/cellp/docker-compose.yml` 的 overlay。

| 组件 | 角色 |
|------|------|
| **OpenObserve** `v1.0.4` | Trace / Log / Metric 存储 + UI |
| **OTel Collector (contrib)** `0.115.1` | OTLP 接入、batch、导出至 OpenObserve |

镜像版本在 `docker-compose.yml` 中**固定 pin**（无 `:latest`）。

## 架构要点

- **单一 compose 文件**：仅 `docker-compose.yml`。历史上的 `docker-compose.cellp-bridge.yml` 已合并并**删除**，勿再文档化或 `-f` 双文件启动。
- **网络**：`openobserve` 仅在 `otel-internal`；`otel-collector` 在 `otel-internal` + **`1panel-network`**（external），供 cellpd 通过 `http://otel-collector:4318` 访问。
- **Collector 导出**：`http://openobserve:5080/api/${OO_ORG}`（默认 `default`），**无尾部斜杠**；`Authorization: Basic ${OO_BASIC_AUTH}`。

## 安全

- UI（`5080`）与 OTLP（`4317`/`4318`）在 compose 中绑定 **127.0.0.1**。勿在未防护时发布到 `0.0.0.0`。
- 凭据在 **`.env`**（gitignore）。从 `.env.example` 复制。
- **`OO_ROOT_PASSWORD`**：至少 16 字符，含大小写、数字、符号；不得与邮箱 local-part 相同。
- **`OO_BASIC_AUTH`**：`OO_ROOT_EMAIL:OO_ROOT_PASSWORD` 的 base64（与 OpenObserve API Basic auth 同一对）。

```bash
echo -n 'admin@example.com:YourSecurePassword' | base64
```

## 部署

依赖外部 Docker 网络 **`1panel-network`**（与 MacPanel cellp / RustFS 相同）。若不存在，先启动任一使用该网络的栈或手动创建。

```bash
cd /Users/mino/code/MacPanel/.dev/macpanel/1panel/1panel/docker/compose/otel
cp .env.example .env
# 编辑 .env：OO_ROOT_*、OO_BASIC_AUTH

docker compose up -d
```

仅校验 compose 语法：

```bash
docker compose config
```

## 端口（主机 localhost）

| 端口 | 服务 |
|------|------|
| `5080` | OpenObserve UI |
| `4317` | OTLP gRPC |
| `4318` | OTLP HTTP（cellp 默认） |

主机健康探测 OpenObserve：`curl -sf http://127.0.0.1:5080/healthz`

## cellp 如何连接

在 `../cellp/.env`：

```env
CELLP_OTEL_COLLECTOR=http://otel-collector:4318
CELLP_OTEL_BACKEND=jaeger
```

| 场景 | 网络 | `CELLP_OTEL_COLLECTOR` |
|------|------|-------------------------|
| Docker cellpd（MacPanel） | `1panel-network` + 容器名 `otel-collector` | `http://otel-collector:4318` |
| 宿主机 cellp | 本机 published 端口 | `http://127.0.0.1:4318` |

修改 cellp env 后**必须**重建 cellpd：

```bash
cd ../cellp
docker compose up -d --force-recreate cellpd
```

## 在哪里看 Trace

**OpenObserve UI：** `http://127.0.0.1:5080`（`OO_ROOT_EMAIL` / `OO_ROOT_PASSWORD`）。

本 setup **没有** Tempo / Jaeger query API。cellp Dashboard → Storage → Telemetry **不能**作为 MacPanel OpenObserve 路径的 trace 浏览器；OTLP 数据经 collector 写入 OpenObserve 后，在 OO UI 的 Traces 视图查询。

## OpenObserve v1.0.4 与 healthcheck

`public.ecr.aws/zinclabs/openobserve:v1.0.4` 为 **distroless**，镜像内无 `wget`/`curl`。

若 compose 中对 OpenObserve 配置 `healthcheck` 且 collector 使用 `depends_on: service_healthy`，collector 会一直处于 **Created**、永不启动。

当前 `docker-compose.yml` 已：

- OpenObserve / collector：`healthcheck: disable: true`
- collector：`depends_on: openobserve: condition: service_started`

升级 OpenObserve 时勿恢复基于 wget 的 healthcheck，除非换用带 shell 的镜像或改用外部探测。

## 1Panel / MacPanel UI：重启与重建

### dual-file 标签（已修复的历史问题）

旧部署曾用：

```text
docker compose -f docker-compose.yml -f docker-compose.cellp-bridge.yml up -d
```

bridge 文件已删除。若容器仍带标签：

```text
com.docker.compose.project.config_files=.../docker-compose.yml,.../docker-compose.cellp-bridge.yml
```

在 MacPanel 编排页点「重启」会失败：

```text
open .../docker-compose.cellp-bridge.yml: no such file or directory
```

**修复：**

```bash
cd compose/otel
docker compose up -d --force-recreate
```

若 UI 仍显示错误 compose 路径，检查 Agent DB（`agent.db`）中 compose 记录的 path，或删除旧编排条目后重新 `up -d`。

### 修改 .env 后

UI「重启」不 reload env。改 `OO_*` 后：

```bash
docker compose up -d --force-recreate
```

## 验证

```bash
docker compose ps
# openobserve、otel-collector 均 Up

curl -sf http://127.0.0.1:5080/healthz

# 从 cellpd 网络侧（可选，在 cellpd 容器内）
# curl -sf http://otel-collector:4318
```

产生请求后，在 OpenObserve UI 选择 Traces、适当时间范围查询。

Collector 进程内 health extension 监听 `13133`（仅容器内，未 publish 到主机）。

## 升级

1. 修改 `docker-compose.yml` 中 pin 的镜像 tag（对照 OpenObserve release notes + collector contrib release）。
2. `docker compose pull`
3. `docker compose up -d`
4. `docker compose ps`；必要时查看 `docker compose logs otel-collector`

## 备份

持久卷：**`macpanel-openobserve-data`**（OpenObserve `/data`）。

```bash
docker run --rm \
  -v macpanel-openobserve-data:/data:ro \
  -v "$(pwd)/backups:/backup" \
  alpine tar czf /backup/openobserve-$(date +%Y%m%d).tar.gz -C /data .
```

恢复：停栈 → 还原卷数据 → 再 `up -d`。

## 停止 / 删除

```bash
docker compose down          # 保留 macpanel-openobserve-data
docker compose down -v       # 删除数据卷 — 破坏性
```

## 防火墙

若必须将 OTLP 或 UI 绑定到 LAN，用本地 override 改 port mapping，并用主机防火墙限制来源 IP（仅 cellp 节点）。

## 总览文档

[../README.md](../README.md) · cellp 接入：[../cellp/OBSERVABILITY.md](../cellp/OBSERVABILITY.md)
