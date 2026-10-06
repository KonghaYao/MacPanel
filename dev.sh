#!/usr/bin/env bash
# Local MacPanel dev launcher — isolated data dir and non-production port.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
DEV_HOME="${MACPANEL_HOME:-$ROOT/.dev/macpanel}"
DEV_PORT="${MACPANEL_PORT:-19999}"
BINARY="${MACPANEL_BIN:-$ROOT/build/macpanel}"

usage() {
	cat <<EOF
Usage: $(basename "$0") [options]

Options:
  -p, --port PORT   HTTP port (default: 19999)
  -b, --build       Full rebuild: clean embed, frontend (npm run build:pro), macpanel
  -q, --quick       Quick go build only (default when binary is missing)
  -h, --help        Show this help

Environment:
  MACPANEL_HOME     Dev data root (default: $ROOT/.dev/macpanel)
  MACPANEL_PORT     Same as --port
  MACPANEL_BIN      macpanel binary path (default: build/macpanel)

Production MacPanel (port 9999) is untouched when using the default dev home.
EOF
}

BUILD=0
QUICK=0
WEB_INDEX="$ROOT/core/cmd/server/web/index.html"

while [[ $# -gt 0 ]]; do
	case "$1" in
	-p | --port)
		DEV_PORT="$2"
		shift 2
		;;
	-b | --build)
		BUILD=1
		shift
		;;
	-q | --quick)
		QUICK=1
		shift
		;;
	-h | --help)
		usage
		exit 0
		;;
	*)
		echo "Unknown option: $1" >&2
		usage
		exit 1
		;;
	esac
done

export MACPANEL_HOME="$DEV_HOME"
export MACPANEL_PORT="$DEV_PORT"

quick_build() {
	echo "==> quick build macpanel (go build, reuse embedded frontend)..."
	mkdir -p "$(dirname "$BINARY")"
	(
		cd "$ROOT/cmd/macpanel"
		CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags '-s -w' -o "$BINARY" .
	)
}

needs_frontend_rebuild=0
if [[ -f "$WEB_INDEX" ]] && grep -q 'unpkg.com' "$WEB_INDEX"; then
	echo "==> detected stale CDN frontend embed; rebuilding frontend..."
	needs_frontend_rebuild=1
fi

full_build_with_clean() {
	echo "==> cleaning embedded web assets..."
	make -C "$ROOT" clean_assets
	echo "==> building frontend (npm run build:pro)..."
	echo "==> building macpanel binary..."
	make -C "$ROOT" build_macpanel
}

if [[ "$BUILD" -eq 1 ]]; then
	full_build_with_clean
elif [[ "$needs_frontend_rebuild" -eq 1 ]] || [[ ! -f "$WEB_INDEX" ]]; then
	echo "==> full build (frontend + macpanel)..."
	make -C "$ROOT" build_macpanel
elif [[ ! -x "$BINARY" ]] || [[ "$QUICK" -eq 1 ]]; then
	quick_build
fi

if [[ ! -x "$BINARY" ]]; then
	echo "Binary not found: $BINARY" >&2
	exit 1
fi

CONFIG="$DEV_HOME/config/1pctl"
mkdir -p "$(dirname "$CONFIG")"
if [[ -f "$CONFIG" ]] && grep -q "^ORIGINAL_PORT=" "$CONFIG"; then
	current_port="$(grep "^ORIGINAL_PORT=" "$CONFIG" | cut -d= -f2-)"
	if [[ "$current_port" != "$DEV_PORT" ]]; then
		sed -i '' "s/^ORIGINAL_PORT=.*/ORIGINAL_PORT=$DEV_PORT/" "$CONFIG"
		echo "==> updated ORIGINAL_PORT: $current_port -> $DEV_PORT"
	fi
fi

echo "==> MacPanel dev"
echo "    home:   $DEV_HOME"
echo "    port:   $DEV_PORT"
echo "    config: $CONFIG"
echo "    url:    http://127.0.0.1:$DEV_PORT"
if [[ -f "$CONFIG" ]]; then
	echo "    user:   $(grep '^ORIGINAL_USERNAME=' "$CONFIG" | cut -d= -f2-)"
	echo "    pass:   $(grep '^ORIGINAL_PASSWORD=' "$CONFIG" | cut -d= -f2-)"
fi
echo ""
echo "Press Ctrl+C to stop (foreground dev mode)."
echo ""

exec "$BINARY"
