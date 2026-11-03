#!/usr/bin/env bash
# Installs a commit-msg hook in this clone that runs check-commit-msg.sh on
# every commit, so a message that breaks the rules is refused before it lands.
# The hook calls the script inside the checkout, so updates to the script
# apply without reinstalling. Hooks are per clone and shared by its worktrees;
# nothing is committed.
#
# Usage, from anywhere inside the checkout:
#   bash ".agents/skills/conorganizer-commits/scripts/install-commit-msg-hook.sh" [--force]
# Refuses to replace an existing commit-msg hook unless --force is given.
# Remove it again with: rm "$(git rev-parse --git-path hooks)/commit-msg"
set -eu

if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
	sed -n '2,11p' "$0" | sed 's/^# \{0,1\}//'
	exit 0
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
checker="$script_dir/check-commit-msg.sh"
if [ ! -f "$checker" ]; then
	echo "check-commit-msg.sh not found next to this script" >&2
	exit 1
fi

hooks_dir="$(git rev-parse --git-path hooks)"
hook="$hooks_dir/commit-msg"
if [ -e "$hook" ] && [ "${1:-}" != "--force" ]; then
	echo "$hook already exists; run again with --force to replace it" >&2
	exit 1
fi

mkdir -p "$hooks_dir"
cat > "$hook" <<'EOF'
#!/bin/sh
exec bash "$(git rev-parse --show-toplevel)/.agents/skills/conorganizer-commits/scripts/check-commit-msg.sh" "$1"
EOF
chmod +x "$hook"
echo "Installed $hook"
echo "Every commit now runs check-commit-msg.sh; a message that breaks the rules is refused with the reason."
