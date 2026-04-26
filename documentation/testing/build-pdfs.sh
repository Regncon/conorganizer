#!/usr/bin/env bash
set -euo pipefail

shopt -s nullglob

for md_file in *.md; do
  pdf_file="${md_file%.md}.pdf"

  echo "Generating ${pdf_file} from ${md_file}"

  rm -f "$pdf_file"

  pandoc "$md_file" \
    -o "$pdf_file" \
    --pdf-engine=xelatex \
    -V geometry:margin=25mm \
    -V fontsize=11pt \
    --highlight-style=tango
done
