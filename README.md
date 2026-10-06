# MacPanel

**MacPanel** is a macOS-focused fork and adaptation of [1Panel](https://github.com/1Panel-dev/1Panel), the open-source Linux server management panel. It ships as a single local `macpanel` binary (Apple Silicon / Intel universal) rather than a Linux server panel installer.

## Prerequisites

1. Install [mise](https://mise.jdx.dev/) and enable it in your shell (see [Getting Started](https://mise.jdx.dev/getting-started.html)).
2. **macOS** only (Apple Silicon or Intel).

## Install (mise + GitHub Release)

When releases are available, install via mise from GitHub Releases:

```bash
mise use -g "github:KonghaYao/MacPanel[bin=macpanel]@0.2.0"
```

Downloads: [GitHub Releases](https://github.com/KonghaYao/MacPanel/releases)

## Run

**Foreground** — terminal must stay open; Ctrl+C stops the service:

```bash
macpanel
```

**Background** (recommended) — survives closing the terminal:

```bash
macpanel -d
# or: macpanel --daemon
```

Daemon mode writes:

- PID: `~/Library/Application Support/MacPanel/run/macpanel.pid`
- Log: `~/Library/Application Support/MacPanel/run/macpanel.log`

Stop a background instance:

```bash
kill $(cat ~/Library/Application\ Support/MacPanel/run/macpanel.pid)
```

Alternative without built-in daemon mode:

```bash
nohup macpanel &
```

Open [http://127.0.0.1:9999](http://127.0.0.1:9999) in your browser after startup.

Config file:

```
~/Library/Application Support/MacPanel/config/1pctl
```

## Docker

Container features require **Docker Desktop**. MacPanel does not install Docker for you — install and start Docker Desktop yourself.

## Development build

From the repository root:

```bash
mise install          # installs Go 1.26.6, Node 20 (see .mise.toml)
make build_macpanel   # builds unified binary at build/macpanel
./build/macpanel
```

## Upstream

MacPanel is derived from [1Panel](https://github.com/1Panel-dev/1Panel). See the upstream project for the original Linux server panel documentation and community.

## License

Licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.html).
