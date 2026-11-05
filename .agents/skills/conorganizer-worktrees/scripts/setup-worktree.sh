#!/usr/bin/env bash
set -eu

git_dir="$(cd "$(git rev-parse --git-dir)" && pwd -P)"
common_dir="$(cd "$(git rev-parse --git-common-dir)" && pwd -P)"
if [ "$git_dir" = "$common_dir" ]; then
  echo "This is the main checkout, not a worktree. Its dev server generates the templ code itself; nothing to do." >&2
  exit 1
fi

cd "$(git rev-parse --show-toplevel)"

for dir in components pages layouts service; do
  if [ -d "$dir" ]; then
    go tool templ generate -path "$dir" -log-level error
  fi
done
echo "templ: generated for components, pages, layouts and service"

if [ -e .claude/skills ]; then
  echo "skills: .claude/skills is already linked"
else
  go tool task setup:skills
fi

go build ./...
echo "build: go build ./... passed"

if [ ! -f .env ]; then
  echo "missing: .env (copy it from the main checkout before running the app; never print it)"
fi
if [ ! -f database/events.db ]; then
  echo "missing: database/events.db (tests do not need it; running the app does, copy it with sqlite .backup)"
fi
echo "next: go test ./<pkg>/... for the packages you touch; never go tool task test here"
