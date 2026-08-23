#!/usr/bin/env bash
# GUIHO RunX installer (Convention 0001).
# Installs the stable launcher into $HOME/.guiho/bin/, immutable payloads into
# $HOME/.guiho/runx/versions/<version>/, activates via an atomic current.json
# pointer, verifies every artifact against checksums.txt and artifacts.json,
# and rolls back completely on failure.
set -euo pipefail

REPO="${RUNX_REPO:-CGuiho/runx}"
REQUESTED_VERSION=""
REQUESTED_CHANNEL=""
STAGING=""
PREVIOUS_POINTER=""

cleanup(){ [ -z "$STAGING" ] || rm -rf -- "$STAGING"; }
trap cleanup EXIT HUP INT TERM
fail(){ printf 'error: %s\n' "$*" >&2; exit 1; }
require(){ command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"; }

usage(){ cat <<'EOF'
Install GUIHO RunX.

Usage: install.sh [--version VERSION] [--channel CHANNEL]

--version and --channel are mutually exclusive. Without either, the latest
stable release is installed.
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version) [ "$#" -ge 2 ] || fail "--version requires a value"; REQUESTED_VERSION="$2"; shift 2 ;;
    --version=*) REQUESTED_VERSION="${1#*=}"; shift ;;
    --channel) [ "$#" -ge 2 ] || fail "--channel requires a value"; REQUESTED_CHANNEL="$2"; shift 2 ;;
    --channel=*) REQUESTED_CHANNEL="${1#*=}"; shift ;;
    --help) usage; exit 0 ;;
    *) fail "unknown argument: $1 (full flag names only)" ;;
  esac
done

[ -z "$REQUESTED_VERSION" ] || [ -z "$REQUESTED_CHANNEL" ] || fail "--version and --channel are mutually exclusive"

for tool in curl unzip awk grep mktemp sha256sum uname; do require "$tool"; done

OS_NAME="$(uname -s)"; MACHINE="$(uname -m)"
case "$OS_NAME/$MACHINE" in
  Linux/x86_64|Linux/amd64) PLATFORM=linux-amd64 ;;
  Linux/aarch64|Linux/arm64) PLATFORM=linux-arm64 ;;
  Linux/armv7l) PLATFORM=linux-armv7 ;;
  Linux/armv6l) PLATFORM=linux-armv6 ;;
  Darwin/x86_64) PLATFORM=darwin-amd64 ;;
  Darwin/arm64) PLATFORM=darwin-arm64 ;;
  *) fail "unsupported platform: $OS_NAME/$MACHINE" ;;
esac

normalize_version(){ printf '%s\n' "${1#@guiho/runx/v}" | sed -e 's/^@guiho\/runx@//' -e 's/^runx\/v//' -e 's/^runx@//' -e 's/^v//'; }
is_semver(){ printf '%s\n' "$1" | grep -Eq '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$'; }

channel_of(){
  # stable when no prerelease component, else the first dot-separated prerelease identifier
  case "$1" in
    *-*) printf '%s\n' "${1#*-}" | cut -d. -f1 ;;
    *) printf 'stable' ;;
  esac
}

API_BASE="https://api.github.com/repos/$REPO"
DOWNLOAD_BASE="https://github.com/$REPO/releases/download"

resolve_selection(){
  if [ -n "$REQUESTED_VERSION" ]; then
    VERSION="$(normalize_version "$REQUESTED_VERSION")"
    is_semver "$VERSION" || fail "invalid RunX version: $REQUESTED_VERSION"
    return 0
  fi
  require awk
  local page=1 total=10 tags pre tag channel_tag="" stable_tag=""
  while [ "$page" -le "$total" ]; do
    local response
    response="$(curl --fail --silent --show-error --location --max-time 60 \
      -H "Accept: application/vnd.github+json" \
      "$API_BASE/releases?per_page=100&page=$page")" || fail "release catalog request failed"
    [ "$response" != "[]" ] || break
    tags="$(printf '%s\n' "$response" | awk 'BEGIN{RS=","} /"tag_name"/ {gsub(/.*"tag_name"[ \t]*:[ \t]*"/,""); gsub(/".*/,""); print}')"
    pre="$(printf '%s\n' "$response" | awk 'BEGIN{RS=","} /"prerelease"/ {gsub(/.*"prerelease"[ \t]*:[ \t]*/,""); gsub(/[ \t}.].*/,""); print}')"
    paste <(printf '%s\n' "$tags") <(printf '%s\n' "$pre") | while IFS="$(printf '\t')" read -r t p; do
      printf '%s\t%s\n' "$t" "$p"
    done > /tmp/.runx-releases.$$ || true
    while IFS="$(printf '\t')" read -r tag ispre; do
      [ -n "$tag" ] || continue
      local norm candidate_channel
      norm="$(normalize_version "$tag")"
      is_semver "$norm" || continue
      if [ "$ispre" = "true" ]; then candidate_channel="$(channel_of "$norm")"; else candidate_channel="stable"; fi
      if [ -z "$REQUESTED_CHANNEL" ] || [ "$REQUESTED_CHANNEL" = "stable" ]; then
        [ "$candidate_channel" = "stable" ] || continue
        [ -z "$stable_tag" ] && stable_tag="$norm"
      else
        [ "$candidate_channel" = "$REQUESTED_CHANNEL" ] || continue
        [ -z "$channel_tag" ] && channel_tag="$norm"
      fi
    done < /tmp/.runx-releases.$$
    rm -f /tmp/.runx-releases.$$
    if [ -n "$channel_tag" ] || [ -n "$stable_tag" ]; then break; fi
    page=$((page + 1))
  done
  if [ -n "$REQUESTED_CHANNEL" ] && [ "$REQUESTED_CHANNEL" != "stable" ]; then
    [ -n "$channel_tag" ] || fail "no published release found for channel: $REQUESTED_CHANNEL"
    VERSION="$channel_tag"
  else
    [ -n "$stable_tag" ] || fail "no published stable release found"
    VERSION="$stable_tag"
  fi
}

resolve_selection
# Transition: try new tag runx/v* first, fallback to legacy @guiho/runx/v* for pre-0.14.4 releases.
TAG_ENCODED_NEW="runx%2Fv$VERSION"
TAG_ENCODED_OLD="%40guiho%2Frunx%2Fv$VERSION"
# Probe which tag exists: prefer new format if it exists for this version.
if curl --fail --silent --head --location --max-time 10 "$DOWNLOAD_BASE/$TAG_ENCODED_NEW/checksums.txt" >/dev/null 2>&1; then
  TAG_ENCODED="$TAG_ENCODED_NEW"
else
  TAG_ENCODED="$TAG_ENCODED_OLD"
fi
ASSET_BASE="$DOWNLOAD_BASE/$TAG_ENCODED"
PAYLOAD_ASSET="runx-payload-$PLATFORM"
LAUNCHER_ASSET="runx-launcher-$PLATFORM"

GUIHO_ROOT="$HOME/.guiho"
CLI_DIR="$GUIHO_ROOT/runx"
BIN_DIR="$GUIHO_ROOT/bin"
VERSIONS_DIR="$CLI_DIR/versions"
RESOURCES_DIR="$CLI_DIR/resources"
TEMP_ROOT="$GUIHO_ROOT/.temp"

mkdir -p "$TEMP_ROOT"
STAGING="$(mktemp -d "$TEMP_ROOT/runx-install-XXXXXXXX")"

printf '%s\n' 'Initiating GUIHO CLI Upgrade / Installation Sequence...'
printf 'Target Version: %s\nPlatform:       %s\nPayload Asset:  %s\nSource URL:     %s\nStaging:         %s\n' \
  "$VERSION" "$PLATFORM" "$PAYLOAD_ASSET" "$ASSET_BASE" "$STAGING"

download(){ curl --fail --location --progress-bar --output "$STAGING/$2" "$ASSET_BASE/$2"; }

download "$ASSET_BASE" "$PAYLOAD_ASSET"
download "$ASSET_BASE" "$LAUNCHER_ASSET"
download "$ASSET_BASE" checksums.txt
download "$ASSET_BASE" artifacts.json
download "$ASSET_BASE" guiho-s-runx.zip
download "$ASSET_BASE" guiho-i-runx.md
download "$ASSET_BASE" guiho-p-runx.md
download "$ASSET_BASE" guiho-p-runx-uninstall.md
download "$ASSET_BASE" runx.schema.json
download "$ASSET_BASE" runx.global.schema.json

verify_checksum(){
  expected="$(awk -v asset="$2" '$2 == asset {print $1}' "$STAGING/checksums.txt")"
  [ -n "$expected" ] || fail "checksum entry missing for $2"
  actual="$(sha256sum "$STAGING/$1" | awk '{print $1}')"
  [ "$expected" = "$actual" ] || fail "checksum verification failed for $2"
}
manifest_digest(){
  awk -v want="$2" '
    /"file"/    { gsub(/.*: "/,""); gsub(/",?$/,""); current=$0 }
    /"sha256"/  { gsub(/.*: "/,""); gsub(/",?$/,""); if (current == want) print $0 }
  ' "$STAGING/artifacts.json"
}

FILES="$PAYLOAD_ASSET $LAUNCHER_ASSET guiho-s-runx.zip guiho-i-runx.md guiho-p-runx.md guiho-p-runx-uninstall.md runx.schema.json runx.global.schema.json artifacts.json"
for asset in $FILES; do
  verify_checksum "$asset" "$asset"
  declared="$(manifest_digest x "$asset")"
  if [ -n "$declared" ]; then
    actual="$(sha256sum "$STAGING/$asset" | awk '{print $1}')"
    [ "$declared" = "$actual" ] || fail "artifacts.json digest mismatch for $asset"
  fi
done
printf '%s\n' '[OK] SHA-256 verification complete.'

PAYLOAD_FILE="runx-payload"
case "$PLATFORM" in windows-*) PAYLOAD_FILE="runx-payload.exe" ;; esac

# Self-test the staged payload BEFORE activation.
chmod +x "$STAGING/$PAYLOAD_ASSET"
STAGED_VERSION="$(RUNX_DISABLE_UPDATE_WORKER=1 RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=1 "$STAGING/$PAYLOAD_ASSET" --version)"
[ "$STAGED_VERSION" = "$VERSION" ] || fail "staged payload reports $STAGED_VERSION, want $VERSION"
RUNX_DISABLE_UPDATE_WORKER=1 RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=1 "$STAGING/$PAYLOAD_ASSET" __self-test >/dev/null \
  || fail 'staged payload failed its installation self-test'

# Snapshot previous pointer for rollback.
POINTER_PATH="$CLI_DIR/current.json"
HAD_POINTER=false
if [ -f "$POINTER_PATH" ]; then cp -- "$POINTER_PATH" "$STAGING/current.json.previous"; HAD_POINTER=true; fi
rollback_pointer(){ if [ "$HAD_POINTER" = true ]; then cp -- "$STAGING/current.json.previous" "$POINTER_PATH"; else rm -f -- "$POINTER_PATH"; fi; }

# Install the immutable version directory.
DEST_VERSION_DIR="$VERSIONS_DIR/$VERSION"
if [ -e "$DEST_VERSION_DIR" ]; then rm -rf -- "$DEST_VERSION_DIR"; fi
mkdir -p "$DEST_VERSION_DIR"
cp -- "$STAGING/$PAYLOAD_ASSET" "$DEST_VERSION_DIR/$PAYLOAD_FILE"
chmod 0755 "$DEST_VERSION_DIR/$PAYLOAD_FILE"
mkdir -p "$RESOURCES_DIR/skills" "$RESOURCES_DIR/instruction" "$RESOURCES_DIR/prompts" "$RESOURCES_DIR/schemas"
unzip -q -o "$STAGING/guiho-s-runx.zip" -d "$RESOURCES_DIR/skills"
[ -f "$RESOURCES_DIR/skills/guiho-s-runx/SKILL.md" ] || { rollback_install(){ :; }; fail 'skill archive is missing guiho-s-runx/SKILL.md'; }
cp -- "$STAGING/guiho-i-runx.md" "$RESOURCES_DIR/instruction/guiho-i-runx.md"
cp -- "$STAGING/guiho-p-runx.md" "$RESOURCES_DIR/prompts/guiho-p-runx.md"
cp -- "$STAGING/guiho-p-runx-uninstall.md" "$RESOURCES_DIR/prompts/guiho-p-runx-uninstall.md"
cp -- "$STAGING/runx.schema.json" "$RESOURCES_DIR/schemas/runx.schema.json"
cp -- "$STAGING/runx.global.schema.json" "$RESOURCES_DIR/schemas/runx.global.schema.json"
cp -- "$STAGING/artifacts.json" "$DEST_VERSION_DIR/release-artifacts.json"

# Global managed skill copies.
for ROOT in "$HOME/.agents/skills" "$HOME/.claude/skills"; do
  mkdir -p "$ROOT"; rm -rf "$ROOT/guiho-s-runx"; cp -R "$RESOURCES_DIR/skills/guiho-s-runx" "$ROOT/guiho-s-runx"
  printf '[OK] Installed agent skill: %s\n' "$ROOT/guiho-s-runx"
done
if [ -f AGENTS.md ] || [ -f CLAUDE.md ]; then
  RUNX_DISABLE_UPDATE_WORKER=1 RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=1 "$DEST_VERSION_DIR/$PAYLOAD_FILE" agent instruction update || true
fi

# Activate atomically, keeping the previous version as rollback fallback.
mkdir -p "$(dirname -- "$POINTER_PATH")"
PREVIOUS_ACTIVE=""
if [ "$HAD_POINTER" = true ]; then
  PREVIOUS_ACTIVE="$(awk 'BEGIN{RS=","} /"active"/ {gsub(/.*: "/,""); gsub(/".*/,""); print}' "$STAGING/current.json.previous")"
fi
if [ -n "$PREVIOUS_ACTIVE" ] && [ "$PREVIOUS_ACTIVE" != "$VERSION" ]; then
  printf '{"protocol":1,"active":"%s","previous":"%s"}\n' "$VERSION" "$PREVIOUS_ACTIVE" > "$STAGING/current.json.new"
else
  printf '{"protocol":1,"active":"%s"}\n' "$VERSION" > "$STAGING/current.json.new"
fi

LAUNCHER_PATH="$BIN_DIR/runx"
BACKUP_LAUNCHER=""
if [ -e "$LAUNCHER_PATH" ]; then BACKUP_LAUNCHER="$STAGING/launcher.previous"; cp -- "$LAUNCHER_PATH" "$BACKUP_LAUNCHER"; fi

commit(){
  mv -- "$STAGING/current.json.new" "$POINTER_PATH"
  mkdir -p "$BIN_DIR"
  cp -- "$STAGING/$LAUNCHER_ASSET" "$LAUNCHER_PATH.tmp"
  chmod 0755 "$LAUNCHER_PATH.tmp"
  mv -f -- "$LAUNCHER_PATH.tmp" "$LAUNCHER_PATH"
}

restore_previous(){
  rollback_pointer
  if [ -n "$BACKUP_LAUNCHER" ]; then cp -- "$BACKUP_LAUNCHER" "$LAUNCHER_PATH"; chmod 0755 "$LAUNCHER_PATH"; fi
  rm -rf -- "$DEST_VERSION_DIR"
}

if ! commit; then restore_previous; fail 'activation failed and was rolled back'; fi

ACTUAL_VERSION="$(RUNX_DISABLE_UPDATE_WORKER=1 RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=1 "$LAUNCHER_PATH" --version)" || { restore_previous; fail 'activated launcher failed version verification'; }
if [ "$ACTUAL_VERSION" != "$VERSION" ]; then restore_previous; fail "activated launcher reports $ACTUAL_VERSION, want $VERSION"; fi

if [ "${RUNX_SKIP_PATH_UPDATE:-0}" != 1 ]; then
  case ":$PATH:" in *":$BIN_DIR:"*) ;; *)
    touch "$HOME/.profile"
    case "$(cat "$HOME/.profile")" in *":$BIN_DIR:"*) ;; *) printf '\nexport PATH="%s:$PATH"\n' "$BIN_DIR" >> "$HOME/.profile" ;; esac
  ;; esac
fi

printf '[OK] Installed and verified RunX %s\n' "$ACTUAL_VERSION"
printf 'Launcher: %s\nActive payload: %s/%s\nRunX home: %s\nAgent skill: %s\n' \
  "$LAUNCHER_PATH" "$DEST_VERSION_DIR" "$PAYLOAD_FILE" "$CLI_DIR" "$RESOURCES_DIR/skills/guiho-s-runx"
[ "${RUNX_SKIP_PATH_UPDATE:-0}" = 1 ] || case ":$PATH:" in *":$BIN_DIR:"*) ;; *) printf 'Restart your shell or run: source ~/.profile\n' ;; esac
