#!/usr/bin/env sh
set -eu

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  darwin) platform='darwin' ;;
  linux) platform='linux' ;;
  *) echo "Unsupported operating system: $os" >&2; exit 1 ;;
esac

machine="$(uname -m)"
case "$machine" in
  x86_64|amd64) architecture='x64' ;;
  arm64|aarch64) architecture='arm64' ;;
  *) echo "Unsupported architecture: $machine" >&2; exit 1 ;;
esac

asset_name="runx-$platform-$architecture"
release_json="$(curl --fail --silent --show-error -H 'Accept: application/vnd.github+json' https://api.github.com/repos/CGuiho/runx/releases/latest)"
url="$(printf '%s' "$release_json" | grep -o '"browser_download_url":"[^"]*' | cut -d '"' -f4 | grep "/$asset_name$" | head -n 1)"
if [ -z "$url" ]; then
  echo "Release asset $asset_name was not found." >&2
  exit 1
fi

directory="${HOME}/.local/bin"
destination="$directory/runx"
mkdir -p "$directory"
curl --fail --location --silent --show-error "$url" -o "$destination"
chmod +x "$destination"

echo "Installed RunX to $destination"
case ":$PATH:" in
  *":$directory:"*) ;;
  *) echo "Add $directory to PATH, then run: runx --help" ;;
esac
