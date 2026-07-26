#!/usr/bin/env bash
set -euo pipefail

REPO="${RUNX_REPO:-CGuiho/runx}"
VERSION="${RUNX_VERSION:-latest}"
INSTALL_DIR="${RUNX_INSTALL_DIR:-$HOME/.local/bin}"
DOWNLOAD_BASE_URL="${RUNX_DOWNLOAD_BASE_URL:-}"
TMP=""
cleanup(){ [ -z "$TMP" ] || rm -rf -- "$TMP"; }
trap cleanup EXIT HUP INT TERM
fail(){ printf 'error: %s\n' "$*" >&2; exit 1; }
require(){ command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"; }

usage(){ cat <<'EOF'
Install GUIHO RunX from the canonical verified Go release matrix.

Usage: install.sh [--version VERSION] [--install-dir DIRECTORY]
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version|-v) [ "$#" -ge 2 ] || fail "$1 requires a value"; VERSION="$2"; shift 2 ;;
    --version=*) VERSION="${1#*=}"; shift ;;
    --install-dir) [ "$#" -ge 2 ] || fail "$1 requires a value"; INSTALL_DIR="$2"; shift 2 ;;
    --install-dir=*) INSTALL_DIR="${1#*=}"; shift ;;
    --help|-h) usage; exit 0 ;;
    *) fail "unknown argument: $1" ;;
  esac
done

require curl; require unzip; require awk; require install; require mktemp; require uname
OS_NAME="$(uname -s)"; MACHINE="$(uname -m)"
case "$OS_NAME/$MACHINE" in
  Linux/x86_64|Linux/amd64) ASSET=runx-linux-amd64 ;;
  Linux/aarch64|Linux/arm64) ASSET=runx-linux-arm64 ;;
  Linux/armv7l) ASSET=runx-linux-armv7 ;;
  Linux/armv6l) ASSET=runx-linux-armv6 ;;
  Darwin/x86_64) ASSET=runx-darwin-amd64 ;;
  Darwin/arm64) ASSET=runx-darwin-arm64 ;;
  *) fail "unsupported platform: $OS_NAME/$MACHINE" ;;
esac

if [ "$VERSION" = latest ]; then
  BASE_URL="https://github.com/$REPO/releases/latest/download"
  TARGET_LABEL=latest
else
  NORMALIZED="${VERSION#@guiho/runx/v}"; NORMALIZED="${NORMALIZED#@guiho/runx@}"; NORMALIZED="${NORMALIZED#v}"
  printf '%s\n' "$NORMALIZED" | grep -Eq '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$' || fail "invalid RunX version: $VERSION"
  VERSION="$NORMALIZED"; TARGET_LABEL="$VERSION"
  BASE_URL="https://github.com/$REPO/releases/download/%40guiho%2Frunx%2Fv$VERSION"
fi
if [ -n "$DOWNLOAD_BASE_URL" ]; then BASE_URL="${DOWNLOAD_BASE_URL%/}"; fi

TMP="$(mktemp -d "${TMPDIR:-/tmp}/runx-install.XXXXXX")"
printf '%s\n' 'Initiating GUIHO CLI Upgrade / Installation Sequence...'
printf 'Target Version: %s\nArchitecture:   %s\nTarget Asset:   %s\nSource URL:     %s/%s\n' "$TARGET_LABEL" "$MACHINE" "$ASSET" "$BASE_URL" "$ASSET"
curl --fail --location --progress-bar --output "$TMP/$ASSET" "$BASE_URL/$ASSET"
curl --fail --location --progress-bar --output "$TMP/checksums.txt" "$BASE_URL/checksums.txt"
curl --fail --location --progress-bar --output "$TMP/guiho-s-runx.zip" "$BASE_URL/guiho-s-runx.zip"

verify_checksum(){ expected="$(awk -v asset="$2" '$2 == asset {print $1}' "$TMP/checksums.txt")"; [ -n "$expected" ] || fail "checksum entry missing for $2"; if command -v sha256sum >/dev/null 2>&1; then actual="$(sha256sum "$1" | awk '{print $1}')"; else require shasum; actual="$(shasum -a 256 "$1" | awk '{print $1}')"; fi; [ "$expected" = "$actual" ] || fail "checksum verification failed for $2"; }
verify_checksum "$TMP/$ASSET" "$ASSET"; verify_checksum "$TMP/guiho-s-runx.zip" guiho-s-runx.zip
printf '%s\n' '[OK] SHA-256 verification complete.'

mkdir -p -- "$INSTALL_DIR"; DESTINATION="$INSTALL_DIR/runx"; BACKUP="$DESTINATION.old.$$"; HAD_BACKUP=false
if [ -e "$DESTINATION" ]; then mv -- "$DESTINATION" "$BACKUP"; HAD_BACKUP=true; fi
if ! install -m 0755 "$TMP/$ASSET" "$DESTINATION"; then [ "$HAD_BACKUP" = false ] || mv -- "$BACKUP" "$DESTINATION"; fail 'binary installation failed and was rolled back'; fi
ACTUAL_VERSION="$(RUNX_DISABLE_UPDATE_WORKER=1 RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=1 "$DESTINATION" --version)" || { rm -f -- "$DESTINATION"; [ "$HAD_BACKUP" = false ] || mv -- "$BACKUP" "$DESTINATION"; fail 'installed binary failed version verification'; }
if [ "$VERSION" != latest ] && [ "$ACTUAL_VERSION" != "$VERSION" ]; then rm -f -- "$DESTINATION"; [ "$HAD_BACKUP" = false ] || mv -- "$BACKUP" "$DESTINATION"; fail "installed version $ACTUAL_VERSION does not match $VERSION"; fi
[ "$HAD_BACKUP" = false ] || rm -f -- "$BACKUP"

mkdir -p "$TMP/skill"; unzip -q "$TMP/guiho-s-runx.zip" -d "$TMP/skill"; [ -f "$TMP/skill/guiho-s-runx/SKILL.md" ] || fail 'skill archive is missing guiho-s-runx/SKILL.md'
for ROOT in "$HOME/.agents/skills" "$HOME/.claude/skills"; do mkdir -p "$ROOT"; rm -rf "$ROOT/guiho-s-runx"; cp -R "$TMP/skill/guiho-s-runx" "$ROOT/guiho-s-runx"; printf '[OK] Installed agent skill: %s\n' "$ROOT/guiho-s-runx"; done
if [ -f AGENTS.md ] || [ -f CLAUDE.md ]; then RUNX_DISABLE_UPDATE_WORKER=1 RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=1 "$DESTINATION" agent instruction update; fi
if [ "${RUNX_SKIP_PATH_UPDATE:-0}" != 1 ]; then case ":$PATH:" in *":$INSTALL_DIR:"*) ;; *) printf '\nexport PATH="%s:$PATH"\n' "$INSTALL_DIR" >> "$HOME/.profile" ;; esac; fi
printf '[OK] Installed and verified RunX %s at %s\n' "$ACTUAL_VERSION" "$DESTINATION"
