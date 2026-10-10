#!/bin/bash

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

if ! command -v zip >/dev/null 2>&1; then
    echo "Error: zip is required to bundle skills." >&2
    exit 1
fi

output_dir="$(pwd)/.local/skills"
mkdir -p "$output_dir"

while IFS= read -r -d '' skill_file; do
    skill_dir="$(dirname "$skill_file")"
    skill_name="$(basename "$skill_dir")"

    (
        cd "$(dirname "$skill_dir")"
        zip -q -r -FS "$output_dir/$skill_name.zip" "$skill_name"
    )

    echo "Bundled $skill_name -> .local/skills/$skill_name.zip"
done < <(find agents/skills -path agents/skills/etc -prune -o -type f -name SKILL.md -print0)
