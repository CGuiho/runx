#!/usr/bin/env sh
set -eu

REPO="${RUNX_REPO:-CGuiho/runx}"
VERSION="${RUNX_VERSION:-latest}"
INSTALL_DIR="${RUNX_INSTALL_DIR:-$HOME/.local/bin}"
DOWNLOAD_BASE_URL="${RUNX_DOWNLOAD_BASE_URL:-}"
ARCH_OVERRIDE=""
VARIANT_OVERRIDE=""
OS=""
ARCH=""
TARGET_VERSION=""
TMP=""

cleanup() { [ -z "$TMP" ] || rm -rf -- "$TMP"; }
trap cleanup EXIT HUP INT TERM
fail() { printf 'error: %s\n' "$*" >&2; exit 1; }
require_command() { command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"; }
require_value() { [ -n "${2:-}" ] || fail "$1 requires a value"; }

usage() {
  cat <<EOF
Install GUIHO RunX as a verified native CLI binary from GitHub Releases.

Usage: install.sh [--version VERSION] [--arch x64|arm64] [--variant baseline|default|modern] [--install-dir DIR]

VERSION may be an exact stable or prerelease version. The default is latest stable.
EOF
}

parse_args() {
  while [ "$#" -gt 0 ]; do
    case "$1" in
      -v|--version) require_value "$1" "${2:-}"; VERSION="$2"; shift 2 ;;
      --version=*) VERSION=${1#*=}; shift ;;
      --arch) require_value "$1" "${2:-}"; ARCH_OVERRIDE="$2"; shift 2 ;;
      --arch=*) ARCH_OVERRIDE=${1#*=}; shift ;;
      --variant) require_value "$1" "${2:-}"; VARIANT_OVERRIDE="$2"; shift 2 ;;
      --variant=*) VARIANT_OVERRIDE=${1#*=}; shift ;;
      --install-dir) require_value "$1" "${2:-}"; INSTALL_DIR="$2"; shift 2 ;;
      --install-dir=*) INSTALL_DIR=${1#*=}; shift ;;
      -h|--help) usage; exit 0 ;;
      *) fail "unknown flag: $1" ;;
    esac
  done
}

detect_os() { case "$(uname -s)" in Linux) printf 'linux\n' ;; Darwin) printf 'darwin\n' ;; *) fail "unsupported OS: $(uname -s)" ;; esac; }
detect_arch() {
  detected=${ARCH_OVERRIDE:-$(uname -m)}
  case "$detected" in x64|x86_64|amd64) printf 'x64\n' ;; arm64|aarch64) printf 'arm64\n' ;; *) fail "unsupported architecture: $detected" ;; esac
}

normalize_version() {
  value=$1
  value=${value#@guiho/runx@}
  value=${value#v}
  printf '%s\n' "$value" | grep -Eq '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$' || fail "invalid RunX version: $1"
  case "$value" in
    *-*)
      prerelease=${value#*-}
      prerelease=${prerelease%%+*}
      old_ifs=$IFS
      IFS=.
      set -- $prerelease
      IFS=$old_ifs
      for identifier do
        if printf '%s\n' "$identifier" | grep -Eq '^[0-9]+$'; then
          case "$identifier" in 0|[1-9]*) ;; 0*) fail "invalid RunX version: $value" ;; esac
        fi
      done
      ;;
  esac
  printf '%s\n' "$value"
}

resolve_target_version() {
  if [ "$VERSION" != latest ]; then normalize_version "$VERSION"; return; fi
  printf 'Resolving latest stable RunX release...\n' >&2
  effective=$(curl --fail --location --silent --show-error --proto '=https' --tlsv1.2 --output /dev/null --write-out '%{url_effective}' "https://github.com/${REPO}/releases/latest")
  tag=${effective##*/}
  tag=$(printf '%s' "$tag" | sed 's/%40/@/g; s/%2[Ff]/\//g')
  normalize_version "$tag"
}

build_candidates() {
  variant=${VARIANT_OVERRIDE:-baseline}
  if [ "$ARCH" = arm64 ]; then
    [ -z "$VARIANT_OVERRIDE" ] || fail '--variant is only valid for x64 installs'
    printf 'runx-%s-arm64\n' "$OS"
    return
  fi
  case "$variant" in
    baseline) printf 'runx-%s-x64-baseline\nrunx-%s-x64\nrunx-%s-x64-modern\n' "$OS" "$OS" "$OS" ;;
    default) printf 'runx-%s-x64\nrunx-%s-x64-baseline\nrunx-%s-x64-modern\n' "$OS" "$OS" "$OS" ;;
    modern) printf 'runx-%s-x64-modern\nrunx-%s-x64\nrunx-%s-x64-baseline\n' "$OS" "$OS" "$OS" ;;
    *) fail "invalid variant: $variant" ;;
  esac
}

verify_native_binary() {
  magic=$(od -An -tx1 -N4 "$1" 2>/dev/null | tr -d ' \n')
  case "$OS:$magic" in
    linux:7f454c46|darwin:cffaedfe|darwin:cefaedfe|darwin:cafebabe|darwin:feedfacf|darwin:feedface) return 0 ;;
    *) return 1 ;;
  esac
}

verify_markdown_asset() {
  path=$1
  expected_name=$2
  [ -s "$path" ] || fail "downloaded Markdown asset is empty: $path"
  magic=$(od -An -tx1 -N2 "$path" 2>/dev/null | tr -d ' \n')
  [ "$magic" != 4d5a ] || fail "downloaded Markdown asset has a Windows executable header: $path"
  if od -An -tx1 "$path" 2>/dev/null | grep -Eq '(^|[[:space:]])00([[:space:]]|$)'; then
    fail "downloaded Markdown asset contains binary NUL bytes: $path"
  fi
  iconv -f UTF-8 -t UTF-8 "$path" >/dev/null 2>&1 || fail "downloaded Markdown asset is not valid UTF-8 text: $path"
  first_line=$(sed -n '1{s/\r$//;p;}' "$path")
  [ "$first_line" = '---' ] || fail "downloaded Markdown asset does not begin with YAML frontmatter: $path"
  awk -v expected="$expected_name" '
    { sub(/\r$/, "") }
    $0 == "name: " expected { found=1 }
    END { exit(found ? 0 : 1) }
  ' "$path" || fail "downloaded Markdown asset identity does not match $expected_name: $path"
}

verify_installed_version() {
  verify_output=$(mktemp "${TMPDIR:-/tmp}/runx-version.XXXXXX") || return 1
  "$1" --version >"$verify_output" 2>&1 &
  verify_pid=$!
  verify_waited=0
  while kill -0 "$verify_pid" 2>/dev/null; do
    if [ "$verify_waited" -ge 10 ]; then
      kill "$verify_pid" 2>/dev/null || true
      wait "$verify_pid" 2>/dev/null || true
      rm -f -- "$verify_output"
      printf 'error: installed RunX version check timed out after 10 seconds\n' >&2
      return 1
    fi
    sleep 1
    verify_waited=$((verify_waited + 1))
  done
  if ! wait "$verify_pid"; then
    cat "$verify_output" >&2
    rm -f -- "$verify_output"
    return 1
  fi
  actual=$(cat "$verify_output")
  rm -f -- "$verify_output"
  [ "$actual" = "$TARGET_VERSION" ] || { printf 'error: installed RunX reported %s; expected %s\n' "${actual:-<empty>}" "$TARGET_VERSION" >&2; return 1; }
}

ensure_path() {
  PATH="$INSTALL_DIR:$PATH"; export PATH
  profile="$HOME/.profile"
  case "${SHELL##*/}" in zsh) profile="$HOME/.zshrc" ;; bash) profile="$HOME/.bashrc" ;; fish) profile="$HOME/.config/fish/config.fish" ;; esac
  mkdir -p -- "$(dirname "$profile")"
  if ! grep -Fq "$INSTALL_DIR" "$profile" 2>/dev/null; then
    if [ "${SHELL##*/}" = fish ]; then printf '\n# Added by runx installer\nfish_add_path %s\n' "$INSTALL_DIR" >>"$profile"
    else printf '\n# Added by runx installer\nexport PATH=%s:$PATH\n' "$INSTALL_DIR" >>"$profile"; fi
  fi
}

download_asset() {
  url=$1
  output=$2
  set +e
  if [ -n "$DOWNLOAD_BASE_URL" ]; then
    http_code=$(curl --location --progress-bar --show-error --output "$output" --write-out '%{http_code}' "$url")
  else
    http_code=$(curl --location --progress-bar --show-error --proto '=https' --tlsv1.2 --output "$output" --write-out '%{http_code}' "$url")
  fi
  curl_exit=$?
  set -e
  [ "$curl_exit" -eq 0 ] || fail "download failed for $url (curl exit $curl_exit)"
  case "$http_code" in
    404) return 44 ;;
    2??) return 0 ;;
    *) fail "download failed for $url (HTTP $http_code)" ;;
  esac
}

install_transactional() {
  source_path=$1
  destination=$2
  backup="${destination}.old.$$.$(date +%s)"
  original_moved=false
  if [ -e "$destination" ]; then mv -- "$destination" "$backup" || fail "could not back up existing RunX at $destination"; original_moved=true; fi
  if ! install -m 0755 "$source_path" "$destination"; then
    [ "$original_moved" = false ] || mv -- "$backup" "$destination" || fail "installation and automatic rollback both failed; backup remains at $backup"
    fail 'installation failed; the previous executable was restored'
  fi
  printf 'Verifying...\n'
  if ! verify_installed_version "$destination"; then
    rm -f -- "$destination"
    [ "$original_moved" = false ] || mv -- "$backup" "$destination" || fail "verification and automatic rollback both failed; backup remains at $backup"
    fail 'installation verification failed; the previous executable was restored'
  fi
  [ "$original_moved" = false ] || rm -f -- "$backup"
}

install_agent_assets() {
  tag_base=${DOWNLOAD_BASE_URL:-"https://github.com/${REPO}/releases/download"}
  skill_url="${tag_base%/}/${encoded_tag}/guiho-s-runx.md"
  prompt_url="${tag_base%/}/${encoded_tag}/guiho-i-runx.md"
  printf 'Downloading skill asset: %s\n' "$skill_url"
  download_asset "$skill_url" "$TMP/guiho-s-runx.md" || fail 'could not download guiho-s-runx.md'
  verify_markdown_asset "$TMP/guiho-s-runx.md" 'guiho-s-runx'
  printf 'Downloading instruction asset: %s\n' "$prompt_url"
  download_asset "$prompt_url" "$TMP/guiho-i-runx.md" || fail 'could not download guiho-i-runx.md'
  verify_markdown_asset "$TMP/guiho-i-runx.md" 'guiho-i-runx'
  for skill_root in "$HOME/.agents/skills/guiho-s-runx" "$HOME/.claude/skills/guiho-s-runx"; do
    mkdir -p -- "$skill_root"
    install -m 0644 "$TMP/guiho-s-runx.md" "$skill_root/SKILL.md"
    printf 'Installed skill: %s\n' "$skill_root/SKILL.md"
  done
  targets=""
  [ ! -f AGENTS.md ] || targets="AGENTS.md"
  [ ! -f CLAUDE.md ] || targets="${targets}${targets:+ }CLAUDE.md"
  [ -n "$targets" ] || targets="AGENTS.md"
  for instruction_file in $targets; do
    printf 'Reconciling instruction file: %s\n' "$instruction_file"
    clean_file="$TMP/instruction-clean"
    if [ -f "$instruction_file" ]; then
      awk '
        $0 == "<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->" { managed=1; next }
        $0 == "<!-- END RUNX -->" { managed=0; next }
        !managed { print }
      ' "$instruction_file" >"$clean_file"
    else
      : >"$clean_file"
    fi
    {
      cat "$clean_file"
      [ ! -s "$clean_file" ] || printf '\n'
      printf '<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->\n'
      cat "$TMP/guiho-i-runx.md"
      printf '\n<!-- END RUNX -->\n'
    } >"$instruction_file"
  done
}

main() {
  parse_args "$@"
  for command in awk cat curl date dirname grep head iconv install mktemp mv od rm sed tr uname; do require_command "$command"; done
  OS=$(detect_os); ARCH=$(detect_arch); TARGET_VERSION=$(resolve_target_version)
  TMP=$(mktemp -d "${TMPDIR:-/tmp}/runx-install.XXXXXX")
  mkdir -p -- "$INSTALL_DIR"
  destination="$INSTALL_DIR/runx"
  encoded_tag="%40guiho%2Frunx%40${TARGET_VERSION}"
  downloaded=""
  first_asset=$(build_candidates | head -n 1)
  source_url="${DOWNLOAD_BASE_URL:-https://github.com/${REPO}/releases/download}/${encoded_tag}/${first_asset}"
  printf 'Initiating GUIHO CLI Upgrade / Installation Sequence...\n'
  printf 'Target Version: v%s\nArchitecture:   %s\nVariant:        %s\nSource URL:     %s\n' "$TARGET_VERSION" "$ARCH" "${VARIANT_OVERRIDE:-baseline}" "$source_url"
  build_candidates | while IFS= read -r asset; do printf '%s\n' "$asset"; done >"$TMP/candidates"
  while IFS= read -r asset; do
    if [ -n "$DOWNLOAD_BASE_URL" ]; then url="${DOWNLOAD_BASE_URL%/}/${encoded_tag}/${asset}"; else url="https://github.com/${REPO}/releases/download/${encoded_tag}/${asset}"; fi
    candidate="$TMP/$asset"
    printf 'Downloading %s\n' "$url"
    if download_asset "$url" "$candidate"; then :; else
      status=$?
      if [ "$status" -eq 44 ]; then printf '  %s is not available; trying the next compatible candidate.\n' "$asset"; continue; fi
      exit "$status"
    fi
    verify_native_binary "$candidate" || fail "downloaded asset $asset is not a native $OS executable"
    downloaded=$candidate
    break
  done <"$TMP/candidates"
  [ -n "$downloaded" ] || fail "no compatible RunX $TARGET_VERSION binary found at https://github.com/${REPO}/releases"
  printf 'Replacing...\n'
  install_transactional "$downloaded" "$destination"
  install_agent_assets
  [ "${RUNX_SKIP_PATH_UPDATE:-0}" = 1 ] || ensure_path
  printf 'Final verification: %s --version\n' "$destination"
  printf 'Installed and verified RunX %s at %s\n' "$TARGET_VERSION" "$destination"
}

if [ "${RUNX_INSTALLER_SOURCE_ONLY:-0}" = 1 ]; then
  return 0 2>/dev/null || exit 0
fi

main "$@"
