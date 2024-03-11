#!/usr/bin/env sh

rm -rf "../web_assets/assets"

yarn && yarn build


original_file="../web_assets/index.html"
output_file="tmp"

sed 's|src="/assets|src="/web_assets/assets|g; s|href="/assets|href="/web_assets/assets|g' "$original_file" > "$output_file"

mv $output_file $original_file