#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")"

shopt -s nullglob

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

linebreak_filter="$tmp_dir/html-linebreaks.lua"
cat > "$linebreak_filter" <<'LUA'
function RawInline(el)
  if el.format:match("html") and el.text:match("^%s*<br%s*/?>%s*$") then
    return pandoc.LineBreak()
  end
end
LUA

for md_file in *.md; do
  pdf_file="${md_file%.md}.pdf"

  echo "Generating ${pdf_file} from ${md_file}"

  rm -f "$pdf_file"

  pandoc "$md_file" \
    -o "$pdf_file" \
    --pdf-engine=xelatex \
    --lua-filter="$linebreak_filter" \
    -V geometry:margin=25mm \
    -V fontsize=11pt \
    --highlight-style=tango
done
