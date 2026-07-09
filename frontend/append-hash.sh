#!/usr/bin/env bash
set -euo pipefail

cd ./dist

# Clear old hashed files (academy-<hash>.js / academy-<hash>.css)
rm -f academy-*.js academy-*.css

for name in academy.css academy.js; do
    ext="${name##*.}"     # css / js
    stem="${name%.*}"     # academy

    hash=$(sha256sum "$name" | cut -c1-8)
    new_name="${stem}-${hash}.${ext}"

    mv "$name" "$new_name"

    # Rewrite references in index.html
    sed -i "s|/${name}|/${new_name}|g" index.html
done

echo "Updated to include hashes"
