#!/usr/bin/env bash
# Update server/assets/top-1m.csv with the latest Tranco top-1M list.
#
#   scripts/update-tranco.sh            # latest daily list
#   scripts/update-tranco.sh Y83YG      # a specific list, by Tranco list ID
#
# The list is validated before it replaces the current file, and the list ID
# is written to server/assets/top-1m.source. Tranco asks users to cite the ID:
# https://tranco-list.eu/
#
# The backend loads the file at startup, so restart it afterwards. Ranks are
# also cached (domain_rank, 24h), so old ranks can linger until they expire.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEST="$ROOT/server/assets/top-1m.csv"
SOURCE_FILE="$ROOT/server/assets/top-1m.source"
API="https://tranco-list.eu/api/lists"
EXPECTED_LINES=1000000

die() { echo "error: $*" >&2; exit 1; }
command -v curl >/dev/null || die "curl is required"

# --- Look up the list --------------------------------------------------------
if [[ $# -ge 1 ]]; then
  meta=$(curl -fsS --max-time 30 "$API/id/$1") || die "couldn't fetch list $1"
else
  meta=$(curl -fsS --max-time 30 "$API/date/latest") || die "couldn't reach the Tranco API"
fi

# Small, fixed-shape JSON: pull fields with sed rather than depend on jq.
field() { sed -n "s/.*\"$1\": *\"\([^\"]*\)\".*/\1/p" <<<"$meta"; }
list_id=$(field list_id)
download=$(field download)
created=$(field created_on)
[[ -n "$list_id" && -n "$download" ]] || die "unexpected API response: $meta"
grep -q '"available": *true' <<<"$meta" || die "list $list_id isn't available yet"

if [[ -f "$SOURCE_FILE" ]] && grep -qx "list_id=$list_id" "$SOURCE_FILE"; then
  echo "Already on Tranco list $list_id. Nothing to do."
  exit 0
fi

# --- Download and validate ---------------------------------------------------
tmp=$(mktemp "${DEST}.XXXXXX")
trap 'rm -f "$tmp"' EXIT

echo "Downloading Tranco list $list_id (created ${created%%T*})…"
curl -fsS --max-time 300 -o "$tmp" "$download" || die "download failed"

lines=$(wc -l <"$tmp")
[[ "$lines" -eq "$EXPECTED_LINES" ]] || die "expected $EXPECTED_LINES lines, got $lines"
head -1 "$tmp" | grep -q '^1,' || die "first line isn't rank 1: $(head -1 "$tmp")"
bad=$(grep -cvE '^[0-9]+,[^,[:space:]]+'$'\r''?$' "$tmp" || true)
[[ "$bad" -eq 0 ]] || die "$bad lines aren't in rank,domain format"

# --- Swap in -----------------------------------------------------------------
chmod 644 "$tmp"
mv "$tmp" "$DEST"
trap - EXIT
printf 'list_id=%s\ncreated_on=%s\nurl=https://tranco-list.eu/list/%s/1000000\n' \
  "$list_id" "$created" "$list_id" >"$SOURCE_FILE"

echo "Updated $(realpath --relative-to="$ROOT" "$DEST") to Tranco list $list_id."
echo "Restart the backend to load it."
