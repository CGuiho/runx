#!/usr/bin/env bash
# GUIHO RunX uninstaller (Convention 0001).
# Shares the uninstallation contract with `runx uninstall` and
# devops/uninstall.ps1. By default removes everything RunX owns; supports
# --preserve-config, --preserve-data, --dry-run, and --yes.
set -euo pipefail

PRESERVE_CONFIG=false
PRESERVE_DATA=false
DRY_RUN=false
ASSUME_YES=false

fail(){ printf 'error: %s\n' "$*" >&2; exit 1; }

usage(){ cat <<'EOF'
Uninstall GUIHO RunX.

Usage: uninstall.sh [--preserve-config] [--preserve-data] [--dry-run] [--yes]

By default every RunX-owned artifact is removed, including global
configuration, persistent data, and the current project's runx.yaml.
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --preserve-config) PRESERVE_CONFIG=true; shift ;;
    --preserve-data) PRESERVE_DATA=true; shift ;;
    --dry-run) DRY_RUN=true; shift ;;
    --yes) ASSUME_YES=true; shift ;;
    --help) usage; exit 0 ;;
    *) fail "unknown argument: $1 (full flag names only)" ;;
  esac
done

GUIHO_ROOT="$HOME/.guiho"
CLI_DIR="$GUIHO_ROOT/runx"
BIN_DIR="$GUIHO_ROOT/bin"
LAUNCHER="$BIN_DIR/runx"
PROJECT_DIR="$(pwd)"

REMOVE=()
PRESERVE=()
add_remove(){ REMOVE+=("$1"); }
add_preserve(){ PRESERVE+=("$1"); }

# Shared infrastructure is never removed.
add_preserve "$GUIHO_ROOT (shared GUIHO home)"
[ -d "$BIN_DIR" ] && add_preserve "$BIN_DIR (shared launcher directory)"
add_preserve "$GUIHO_ROOT/.temp (shared staging root)"
add_preserve "user PATH entry for $BIN_DIR"

if [ -f "$LAUNCHER" ]; then add_remove "$LAUNCHER"; fi
if [ -d "$CLI_DIR" ]; then
  if [ "$PRESERVE_CONFIG" = true ] || [ "$PRESERVE_DATA" = true ]; then
    add_preserve "$CLI_DIR/runx.global.yaml"
    add_preserve "$CLI_DIR/data"
    add_preserve "$CLI_DIR directory itself (preserved children remain)"
    add_remove "$CLI_DIR/versions"
    add_remove "$CLI_DIR/resources"
    add_remove "$CLI_DIR/current.json"
    add_remove "$CLI_DIR/installed-artifacts.json"
    add_remove "other non-preserved children of $CLI_DIR"
  else
    add_remove "$CLI_DIR"
  fi
fi
for ROOT in "$HOME/.agents/skills" "$HOME/.claude/skills"; do
  [ -e "$ROOT/guiho-s-runx" ] && add_remove "$ROOT/guiho-s-runx"
done
for MARKER in AGENTS.md CLAUDE.md; do
  if [ -f "$PROJECT_DIR/$MARKER" ]; then add_remove "RunX managed block in $MARKER"; fi
done
if [ -f "$PROJECT_DIR/runx.yaml" ]; then
  if [ "$PRESERVE_CONFIG" = true ]; then add_preserve "$PROJECT_DIR/runx.yaml"; else add_remove "$PROJECT_DIR/runx.yaml"; fi
fi

printf 'Uninstallation plan\n\nREMOVE:\n'
if [ "${#REMOVE[@]}" -eq 0 ]; then printf '  (nothing)\n'; else printf '  %s\n' "${REMOVE[@]}"; fi
printf '\nPRESERVE:\n'
printf '  %s\n' "${PRESERVE[@]}"
printf '\n'

interactive(){ [ -t 0 ] && [ -t 1 ]; }

if [ "$DRY_RUN" = true ]; then
  printf 'Dry run: nothing was changed.\n'
  exit 0
fi

if [ "$ASSUME_YES" != true ]; then
  if ! interactive; then fail 'noninteractive uninstall requires --yes'; fi
  printf 'Proceed with uninstallation? [y/N] '
  read -r answer
  case "$answer" in y|Y|yes|YES) ;; *) printf 'Aborted.\n'; exit 0 ;; esac
fi

managed_block_remove(){
  file="$1"; [ -f "$file" ] || return 0
  awk 'BEGIN{skip=0}
    /^<!-- BEGIN RUNX/ { skip=1; next }
    /^<!-- END RUNX/   { skip=0; next }
    skip==0 { print }' "$file" > "$file.tmp" && mv -- "$file.tmp" "$file"
}

remove_path(){
  target="$1"
  # Fail closed on anything that escapes CLI-owned roots.
  case "$target" in
    "$CLI_DIR"|"$CLI_DIR"/*|"$LAUNCHER"|"$HOME/.agents/skills/guiho-s-runx"|"$HOME/.claude/skills/guiho-s-runx") ;;
    *) fail "refusing to remove unowned path: $target" ;;
  esac
  rm -rf -- "$target"
}

[ -e "$LAUNCHER" ] && remove_path "$LAUNCHER"
if [ -d "$CLI_DIR" ]; then
  if [ "$PRESERVE_CONFIG" = false ] && [ "$PRESERVE_DATA" = false ]; then
    remove_path "$CLI_DIR"
  else
    remove_path "$CLI_DIR/versions" 2>/dev/null || true
    remove_path "$CLI_DIR/resources" 2>/dev/null || true
    rm -f -- "$CLI_DIR/current.json" "$CLI_DIR/installed-artifacts.json" 2>/dev/null || true
    find "$CLI_DIR" -mindepth 1 -maxdepth 1 \
      ! -name 'runx.global.yaml' ! -name 'data' -exec rm -rf -- {} +
    rmdir "$CLI_DIR" 2>/dev/null || true
  fi
fi
remove_path "$HOME/.agents/skills/guiho-s-runx" 2>/dev/null || true
remove_path "$HOME/.claude/skills/guiho-s-runx" 2>/dev/null || true
for MARKER in AGENTS.md CLAUDE.md; do managed_block_remove "$PROJECT_DIR/$MARKER"; done
if [ "$PRESERVE_CONFIG" = false ] && [ -f "$PROJECT_DIR/runx.yaml" ]; then
  rm -f -- "$PROJECT_DIR/runx.yaml"
fi

printf '[OK] RunX uninstalled.\n'
