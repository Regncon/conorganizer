#!/usr/bin/env sh
set -eu

projectRoot="$(cd "$(dirname "$0")/.." && pwd)"

if [ ! -d "$projectRoot/.agents/skills" ]; then
    echo "Missing $projectRoot/.agents/skills" >&2
    exit 1
fi

mkdir -p "$projectRoot/.claude"

if [ -e "$projectRoot/.claude/skills" ]; then
    echo ".claude/skills already exists. Nothing to do."
    exit 0
fi

ln -s ../.agents/skills "$projectRoot/.claude/skills"
echo "Linked .claude/skills to .agents/skills"
