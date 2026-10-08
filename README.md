# MacPanel

**MacPanel** is a macOS-focused fork and adaptation of [1Panel](https://github.com/1Panel-dev/1Panel), the open-source Linux server management panel. It ships as a single local `macpanel` binary (Apple Silicon / Intel universal) rather than a Linux server panel installer.

## Prerequisites

1. Install [mise](https://mise.jdx.dev/) and enable it in your shell:

   ```bash
   eval "$(mise activate zsh)"   # zsh
   # eval "$(mise activate bash)"  # bash
   ```

   Add the `eval` line to `~/.zshrc` (or `~/.bashrc`) so `macpanel` is on PATH in new terminals. See [Getting Started](https://mise.jdx.dev/getting-started.html).

2. **macOS** only (Apple Silicon or Intel).

## Install (mise + GitHub Release)

When releases are available, install the **latest** release via mise from GitHub Releases:

```bash
mise use -g "github:KonghaYao/MacPanel[bin=macpanel]@latest"
```

Releases are tagged `v*` (e.g. `v0.4.5`). `@latest` resolves to the newest stable [GitHub release](https://github.com/KonghaYao/MacPanel/releases/latest) (non-draft, non-prerelease).

Downloads: [GitHub Releases](https://github.com/KonghaYao/MacPanel/releases)

## Update

To upgrade to the latest release:

```bash
mise use -g "github:KonghaYao/MacPanel[bin=macpanel]@latest"
mise reshim
```

Or, if `macpanel` is already installed via mise:

```bash
mise upgrade macpanel
mise reshim
```

To pin a specific version instead of latest, replace `@latest` with the version (without the `v` prefix), e.g. `@0.4.5`.

## Run

Foreground (debug): `macpanel` — keeps the process in the terminal; Ctrl+C stops it.

Background service:

```bash
macpanel start
macpanel stop
macpanel restart
macpanel status
```

| | |
|---|---|
| PID file | `~/Library/Application Support/MacPanel/run/macpanel.pid` |
| Log file | `~/Library/Application Support/MacPanel/run/macpanel.log` |

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
